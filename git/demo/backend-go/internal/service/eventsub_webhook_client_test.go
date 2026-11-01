package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

const testCallback = "https://example.test/api/subathon/eventsub/callback"

type fakeAppToken struct {
	token string
	err   error
}

func (f fakeAppToken) AppToken(ctx context.Context) (string, error) { return f.token, f.err }

type recordedRequest struct {
	Method string
	Path   string
	Query  string
	Auth   string
	Client string
	Body   map[string]interface{}
}

// fakeTwitch records every request and answers GET (list), POST (create) and
// DELETE with whatever the test configured.
type fakeTwitch struct {
	mu          sync.Mutex
	requests    []recordedRequest
	existing    []map[string]interface{}
	listStatus  int
	listBody    string
	createCode  int
	createBody  string
	deleteCode  int
	srv         *httptest.Server
	created     []string
	deletedByID []string
}

func newFakeTwitch(t *testing.T) *fakeTwitch {
	f := &fakeTwitch{listStatus: 200, createCode: 202, deleteCode: 204}
	f.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		defer f.mu.Unlock()
		rec := recordedRequest{Method: r.Method, Path: r.URL.Path, Query: r.URL.RawQuery,
			Auth: r.Header.Get("Authorization"), Client: r.Header.Get("Client-Id")}
		if r.Method == http.MethodPost {
			raw, _ := io.ReadAll(r.Body)
			_ = json.Unmarshal(raw, &rec.Body)
			if typ, ok := rec.Body["type"].(string); ok {
				f.created = append(f.created, typ)
			}
		}
		f.requests = append(f.requests, rec)

		switch r.Method {
		case http.MethodGet:
			if f.listStatus >= 300 {
				w.WriteHeader(f.listStatus)
				_, _ = w.Write([]byte(f.listBody))
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"data": f.existing})
		case http.MethodPost:
			w.WriteHeader(f.createCode)
			_, _ = w.Write([]byte(f.createBody))
		case http.MethodDelete:
			f.deletedByID = append(f.deletedByID, r.URL.Query().Get("id"))
			w.WriteHeader(f.deleteCode)
		}
	}))
	t.Cleanup(f.srv.Close)
	return f
}

func (f *fakeTwitch) client(token AppTokenProvider) *eventSubWebhookClient {
	c := newEventSubWebhookClient(token, "client-123", "s3cret-s3cret", testCallback)
	c.apiURL = f.srv.URL + "/helix/eventsub/subscriptions"
	return c
}

func sub(id, typ, status, callback string) map[string]interface{} {
	return map[string]interface{}{
		"id": id, "type": typ, "status": status,
		"transport": map[string]interface{}{"method": "webhook", "callback": callback},
	}
}

func TestEnsureUser_UsesAppTokenNotUserToken(t *testing.T) {
	f := newFakeTwitch(t)
	err := f.client(fakeAppToken{token: "APP-TOKEN"}).ensureUser(context.Background(), "42", subathonSubscriptionSpecs, "Subathon")
	if err != nil {
		t.Fatalf("ensureUser: %v", err)
	}
	if len(f.requests) != 4 { // 1 list + 3 create
		t.Fatalf("expected 4 requests (list + 3 create), got %d", len(f.requests))
	}
	for _, r := range f.requests {
		if r.Auth != "Bearer APP-TOKEN" {
			t.Errorf("%s %s used Authorization %q, want the app token", r.Method, r.Path, r.Auth)
		}
		if r.Client != "client-123" {
			t.Errorf("%s: Client-Id %q", r.Method, r.Client)
		}
	}
	if f.requests[0].Method != http.MethodGet || !strings.Contains(f.requests[0].Query, "user_id=42") {
		t.Errorf("first request should list by user_id, got %+v", f.requests[0])
	}
}

func TestEnsureUser_CreatesWebhookWithSecretCallbackAndCondition(t *testing.T) {
	f := newFakeTwitch(t)
	_ = f.client(fakeAppToken{token: "t"}).ensureUser(context.Background(), "42", subathonSubscriptionSpecs, "Subathon")

	want := map[string]bool{"channel.subscribe": false, "channel.subscription.gift": false, "channel.cheer": false}
	for _, r := range f.requests {
		if r.Method != http.MethodPost {
			continue
		}
		typ := r.Body["type"].(string)
		want[typ] = true
		if r.Body["version"] != "1" {
			t.Errorf("%s: version %v", typ, r.Body["version"])
		}
		cond := r.Body["condition"].(map[string]interface{})
		if cond["broadcaster_user_id"] != "42" {
			t.Errorf("%s: condition %v", typ, cond)
		}
		tr := r.Body["transport"].(map[string]interface{})
		if tr["method"] != "webhook" || tr["callback"] != testCallback || tr["secret"] != "s3cret-s3cret" {
			t.Errorf("%s: transport %v", typ, tr)
		}
	}
	for typ, seen := range want {
		if !seen {
			t.Errorf("subscription %s was never created", typ)
		}
	}
}

func TestEnsureUser_SkipsEnabledAndPendingCreatesOnlyMissing(t *testing.T) {
	f := newFakeTwitch(t)
	f.existing = []map[string]interface{}{
		sub("a", "channel.subscribe", "enabled", testCallback),
		sub("b", "channel.subscription.gift", "webhook_callback_verification_pending", testCallback),
	}
	_ = f.client(fakeAppToken{token: "t"}).ensureUser(context.Background(), "42", subathonSubscriptionSpecs, "Subathon")

	if len(f.created) != 1 || f.created[0] != "channel.cheer" {
		t.Fatalf("expected only channel.cheer to be created, got %v", f.created)
	}
}

func TestEnsureUser_IgnoresSubscriptionsOnForeignCallbacks(t *testing.T) {
	f := newFakeTwitch(t)
	f.existing = []map[string]interface{}{
		sub("old", "channel.subscribe", "enabled", "https://timer.goalkeeper91.de/some/old/callback"),
	}
	_ = f.client(fakeAppToken{token: "t"}).ensureUser(context.Background(), "42", subathonSubscriptionSpecs, "Subathon")

	if len(f.created) != 3 {
		t.Fatalf("a subscription for another callback must not count - expected 3 creates, got %v", f.created)
	}
	if len(f.deletedByID) != 0 {
		t.Fatalf("foreign subscriptions must never be deleted, got %v", f.deletedByID)
	}
}

func TestEnsureUser_DeletesDeadSubscriptionAndRecreatesIt(t *testing.T) {
	f := newFakeTwitch(t)
	f.existing = []map[string]interface{}{
		sub("dead1", "channel.subscribe", "webhook_callback_verification_failed", testCallback),
		sub("dead2", "channel.cheer", "notification_failures_exceeded", testCallback),
		sub("ok", "channel.subscription.gift", "enabled", testCallback),
	}
	_ = f.client(fakeAppToken{token: "t"}).ensureUser(context.Background(), "42", subathonSubscriptionSpecs, "Subathon")

	if strings.Join(f.deletedByID, ",") != "dead1,dead2" {
		t.Fatalf("expected the two dead subscriptions to be deleted, got %v", f.deletedByID)
	}
	got := map[string]bool{}
	for _, typ := range f.created {
		got[typ] = true
	}
	if !got["channel.subscribe"] || !got["channel.cheer"] || got["channel.subscription.gift"] {
		t.Fatalf("expected recreate of the 2 dead types only, got %v", f.created)
	}
}

func TestEnsureUser_TwitchErrorBodyEndsUpInTheError(t *testing.T) {
	f := newFakeTwitch(t)
	f.listStatus = 401
	f.listBody = `{"error":"Unauthorized","status":401,"message":"Invalid OAuth token"}`
	err := f.client(fakeAppToken{token: "t"}).ensureUser(context.Background(), "42", subathonSubscriptionSpecs, "Subathon")
	if err == nil || !strings.Contains(err.Error(), "401") || !strings.Contains(err.Error(), "Invalid OAuth token") {
		t.Fatalf("list error should carry status and Twitch's message, got %v", err)
	}
}

func TestCreateError_IncludesTwitchMessage(t *testing.T) {
	f := newFakeTwitch(t)
	f.createCode = 400
	f.createBody = `{"error":"Bad Request","status":400,"message":"The secret must be between 10 and 100 characters"}`
	c := f.client(fakeAppToken{token: "t"})
	err := c.create(context.Background(), "t", "42", subathonSubscriptionSpecs[0])
	if err == nil || !strings.Contains(err.Error(), "secret must be between") {
		t.Fatalf("create error should include Twitch's explanation, got %v", err)
	}
}

func TestEnsureUser_CreateFailureDoesNotAbortTheRest(t *testing.T) {
	f := newFakeTwitch(t)
	f.createCode = 403
	f.createBody = `{"message":"subscription missing proper authorization"}`
	err := f.client(fakeAppToken{token: "t"}).ensureUser(context.Background(), "42", subathonSubscriptionSpecs, "Subathon")
	if err != nil {
		t.Fatalf("a failing create is logged, not returned: %v", err)
	}
	if len(f.created) != 3 {
		t.Fatalf("all 3 types should still be attempted, got %v", f.created)
	}
}

func TestEnsureUser_AppTokenFailureStopsBeforeCallingTwitch(t *testing.T) {
	f := newFakeTwitch(t)
	err := f.client(fakeAppToken{err: errors.New("token endpoint down")}).ensureUser(context.Background(), "42", subathonSubscriptionSpecs, "Subathon")
	if err == nil || !strings.Contains(err.Error(), "app access token") {
		t.Fatalf("expected an app-token error, got %v", err)
	}
	if len(f.requests) != 0 {
		t.Fatalf("no Twitch request may be made without a token, got %d", len(f.requests))
	}
}

func TestActivitySpecs_ConditionsAndVersions(t *testing.T) {
	byType := map[string]eventSubSpec{}
	for _, s := range activitySubscriptionSpecs {
		byType[s.Type] = s
	}
	if len(byType) != 5 {
		t.Fatalf("expected 5 activity subscription types, got %d", len(byType))
	}
	if s := byType["channel.follow"]; s.Version != "2" ||
		s.Condition("7")["broadcaster_user_id"] != "7" || s.Condition("7")["moderator_user_id"] != "7" {
		t.Errorf("channel.follow must be v2 with broadcaster + moderator, got v%s %v", s.Version, s.Condition("7"))
	}
	if c := byType["channel.raid"].Condition("7"); c["to_broadcaster_user_id"] != "7" || len(c) != 1 {
		t.Errorf("channel.raid must be keyed on the recipient, got %v", c)
	}
	for _, typ := range []string{"channel.subscribe", "channel.subscription.gift", "channel.cheer"} {
		if c := byType[typ].Condition("7"); c["broadcaster_user_id"] != "7" || byType[typ].Version != "1" {
			t.Errorf("%s: unexpected spec %v v%s", typ, c, byType[typ].Version)
		}
	}
}
