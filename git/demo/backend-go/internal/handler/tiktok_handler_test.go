package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"demo/backend-go/internal/composio"
)

type fakeComposio struct {
	configured bool

	linkURL string
	linkErr error
	gotUser string
	gotCB   string

	account    *composio.Account
	accountErr error

	toolData json.RawMessage
	toolErr  error

	deleted   string
	deleteErr error
}

func (f *fakeComposio) Configured() bool { return f.configured }
func (f *fakeComposio) CreateConnectLink(_ context.Context, _ string, userID, cb string) (*composio.ConnectLink, error) {
	f.gotUser, f.gotCB = userID, cb
	if f.linkErr != nil {
		return nil, f.linkErr
	}
	return &composio.ConnectLink{RedirectURL: f.linkURL, ConnectedAccountID: "ca_1"}, nil
}
func (f *fakeComposio) FindActiveAccount(_ context.Context, userID, _ string) (*composio.Account, error) {
	f.gotUser = userID
	return f.account, f.accountErr
}
func (f *fakeComposio) ExecuteTool(_ context.Context, _, _, _ string, _ map[string]any) (json.RawMessage, error) {
	return f.toolData, f.toolErr
}
func (f *fakeComposio) DeleteAccount(_ context.Context, id string, _ bool) error {
	f.deleted = id
	return f.deleteErr
}

func newHandler(f *fakeComposio) *TikTokHandler {
	return NewTikTokHandler(f, "ac_1", "https://goalkeeper91.de/", true)
}

func do(h http.HandlerFunc, method string, cookie *http.Cookie) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, "/x", nil)
	req.Header.Set("X-Real-IP", "203.0.113.7")
	if cookie != nil {
		req.AddCookie(cookie)
	}
	rec := httptest.NewRecorder()
	h(rec, req)
	return rec
}

func decode(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var out map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("invalid json %q: %v", rec.Body.String(), err)
	}
	return out
}

const validUID = "0123456789abcdef0123456789abcdef"

func TestNotConfigured(t *testing.T) {
	h := NewTikTokHandler(&fakeComposio{configured: false}, "ac_1", "https://goalkeeper91.de", true)
	if rec := do(h.Connect, "POST", nil); rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("connect: want 503, got %d", rec.Code)
	}
	h2 := NewTikTokHandler(&fakeComposio{configured: true}, "", "https://goalkeeper91.de", true)
	if rec := do(h2.Me, "GET", nil); rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("me without auth config id: want 503, got %d", rec.Code)
	}
}

func TestConnectSetsHardenedCookieAndCallback(t *testing.T) {
	f := &fakeComposio{configured: true, linkURL: "https://connect.composio.dev/link/x"}
	rec := do(newHandler(f).Connect, "POST", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("want 200, got %d: %s", rec.Code, rec.Body.String())
	}
	if got := decode(t, rec)["redirectUrl"]; got != "https://connect.composio.dev/link/x" {
		t.Fatalf("unexpected redirectUrl %v", got)
	}
	if f.gotCB != "https://goalkeeper91.de/tiktok" {
		t.Fatalf("callback url: %q", f.gotCB)
	}
	if !strings.HasPrefix(f.gotUser, tiktokUserPrefix) {
		t.Fatalf("composio user id not namespaced: %q", f.gotUser)
	}

	var c *http.Cookie
	for _, ck := range rec.Result().Cookies() {
		if ck.Name == tiktokCookieName {
			c = ck
		}
	}
	if c == nil {
		t.Fatal("visitor cookie not set")
	}
	if !c.HttpOnly || !c.Secure || c.SameSite != http.SameSiteLaxMode || c.MaxAge != tiktokCookieMaxAge {
		t.Fatalf("cookie not hardened: %+v", c)
	}
	if !tiktokUIDPattern.MatchString(c.Value) {
		t.Fatalf("cookie value has wrong shape: %q", c.Value)
	}
}

func TestConnectReusesExistingCookie(t *testing.T) {
	f := &fakeComposio{configured: true, linkURL: "https://connect.composio.dev/l"}
	rec := do(newHandler(f).Connect, "POST", &http.Cookie{Name: tiktokCookieName, Value: validUID})
	if len(rec.Result().Cookies()) != 0 {
		t.Fatal("must not mint a new cookie when a valid one is present")
	}
	if f.gotUser != tiktokUserPrefix+validUID {
		t.Fatalf("wrong user id: %q", f.gotUser)
	}
}

func TestConnectRejectsNonHTTPSRedirect(t *testing.T) {
	f := &fakeComposio{configured: true, linkURL: "javascript:alert(1)"}
	if rec := do(newHandler(f).Connect, "POST", nil); rec.Code != http.StatusBadGateway {
		t.Fatalf("want 502 for non-https redirect, got %d", rec.Code)
	}
}

func TestConnectUpstreamError(t *testing.T) {
	f := &fakeComposio{configured: true, linkErr: errors.New("boom")}
	rec := do(newHandler(f).Connect, "POST", nil)
	if rec.Code != http.StatusBadGateway {
		t.Fatalf("want 502, got %d", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "boom") {
		t.Fatal("upstream error detail leaked to the client")
	}
}

func TestMeWithoutOrWithGarbageCookie(t *testing.T) {
	f := &fakeComposio{configured: true}
	h := newHandler(f)
	for _, ck := range []*http.Cookie{nil, {Name: tiktokCookieName, Value: "../../etc"}, {Name: tiktokCookieName, Value: "ABC"}} {
		rec := do(h.Me, "GET", ck)
		if rec.Code != http.StatusOK || decode(t, rec)["connected"] != false {
			t.Fatalf("want connected=false for %+v, got %d %s", ck, rec.Code, rec.Body.String())
		}
	}
	if f.gotUser != "" {
		t.Fatal("must not query Composio for an invalid/missing cookie")
	}
}

func TestMeNotConnected(t *testing.T) {
	f := &fakeComposio{configured: true, account: nil}
	rec := do(newHandler(f).Me, "GET", &http.Cookie{Name: tiktokCookieName, Value: validUID})
	if rec.Code != http.StatusOK || decode(t, rec)["connected"] != false {
		t.Fatalf("got %d %s", rec.Code, rec.Body.String())
	}
}

func TestMeConnectedReturnsProfile(t *testing.T) {
	f := &fakeComposio{
		configured: true,
		account:    &composio.Account{ID: "ca_1"},
		toolData:   json.RawMessage(`{"data":{"user":{"display_name":"Marcel","avatar_url":"https://p16.tiktokcdn.com/a.jpg"}}}`),
	}
	rec := do(newHandler(f).Me, "GET", &http.Cookie{Name: tiktokCookieName, Value: validUID})
	out := decode(t, rec)
	if rec.Code != http.StatusOK || out["connected"] != true || out["displayName"] != "Marcel" || out["avatarUrl"] != "https://p16.tiktokcdn.com/a.jpg" {
		t.Fatalf("got %d %v", rec.Code, out)
	}
	if rec.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("profile responses must not be cacheable")
	}
}

func TestMeDropsNonHTTPSAvatar(t *testing.T) {
	f := &fakeComposio{
		configured: true,
		account:    &composio.Account{ID: "ca_1"},
		toolData:   json.RawMessage(`{"display_name":"M","avatar_url":"javascript:alert(1)"}`),
	}
	out := decode(t, do(newHandler(f).Me, "GET", &http.Cookie{Name: tiktokCookieName, Value: validUID}))
	if out["avatarUrl"] != "" {
		t.Fatalf("non-https avatar must be dropped, got %v", out["avatarUrl"])
	}
}

func TestMeProfileToolFailure(t *testing.T) {
	f := &fakeComposio{configured: true, account: &composio.Account{ID: "ca_1"}, toolErr: errors.New("missing scope")}
	rec := do(newHandler(f).Me, "GET", &http.Cookie{Name: tiktokCookieName, Value: validUID})
	if rec.Code != http.StatusBadGateway || decode(t, rec)["connected"] != true {
		t.Fatalf("got %d %s", rec.Code, rec.Body.String())
	}
	if strings.Contains(rec.Body.String(), "missing scope") {
		t.Fatal("upstream error detail leaked to the client")
	}
}

func TestDisconnect(t *testing.T) {
	f := &fakeComposio{configured: true, account: &composio.Account{ID: "ca_9"}}
	rec := do(newHandler(f).Disconnect, "POST", &http.Cookie{Name: tiktokCookieName, Value: validUID})
	if rec.Code != http.StatusOK || decode(t, rec)["disconnected"] != true || f.deleted != "ca_9" {
		t.Fatalf("got %d %s deleted=%q", rec.Code, rec.Body.String(), f.deleted)
	}

	f2 := &fakeComposio{configured: true, account: nil}
	rec2 := do(newHandler(f2).Disconnect, "POST", &http.Cookie{Name: tiktokCookieName, Value: validUID})
	if decode(t, rec2)["disconnected"] != false || f2.deleted != "" {
		t.Fatal("nothing to disconnect must be a no-op")
	}
}

func TestRateLimiter(t *testing.T) {
	l := newRateLimiter(3, time.Hour)
	for i := 0; i < 3; i++ {
		if !l.allow("a") {
			t.Fatalf("hit %d should be allowed", i+1)
		}
	}
	if l.allow("a") {
		t.Fatal("4th hit must be blocked")
	}
	if !l.allow("b") {
		t.Fatal("other keys are independent")
	}

	short := newRateLimiter(1, 20*time.Millisecond)
	short.allow("k")
	if short.allow("k") {
		t.Fatal("second hit inside the window must be blocked")
	}
	time.Sleep(30 * time.Millisecond)
	if !short.allow("k") {
		t.Fatal("window expiry must free the key again")
	}
}

func TestConnectIsRateLimited(t *testing.T) {
	f := &fakeComposio{configured: true, linkURL: "https://connect.composio.dev/l"}
	h := newHandler(f)
	for i := 0; i < 10; i++ {
		if rec := do(h.Connect, "POST", nil); rec.Code != http.StatusOK {
			t.Fatalf("request %d: %d", i+1, rec.Code)
		}
	}
	if rec := do(h.Connect, "POST", nil); rec.Code != http.StatusTooManyRequests {
		t.Fatalf("11th request must be 429, got %d", rec.Code)
	}
}

func TestExtractTikTokProfileShapes(t *testing.T) {
	cases := map[string]struct{ name, avatar string }{
		`{"display_name":"A","avatar_url":"https://x/a"}`:                          {"A", "https://x/a"},
		`{"data":{"user":{"display_name":"B","avatar_large_url":"https://x/b"}}}`:  {"B", "https://x/b"},
		`"{\"user\":{\"display_name\":\"C\",\"avatar_url_100\":\"https://x/c\"}}"`: {"C", "https://x/c"},
		`[{"display_name":"D"}]`: {"D", ""},
		`not json`:               {"", ""},
	}
	for in, want := range cases {
		n, a := extractTikTokProfile(json.RawMessage(in))
		if n != want.name || a != want.avatar {
			t.Errorf("%s -> (%q,%q), want (%q,%q)", in, n, a, want.name, want.avatar)
		}
	}
}
