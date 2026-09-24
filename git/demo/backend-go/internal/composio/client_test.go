package composio

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return NewClient(srv.URL, "test-key")
}

func TestNotConfiguredWithoutKey(t *testing.T) {
	c := NewClient("", "")
	if c.Configured() {
		t.Fatal("client without API key must not report Configured()")
	}
	if _, err := c.CreateConnectLink(context.Background(), "ac_1", "u1", ""); err == nil {
		t.Fatal("expected error when API key is missing")
	}
}

func TestCreateConnectLink(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/connected_accounts/link" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		if r.Header.Get("x-api-key") != "test-key" {
			t.Errorf("missing/incorrect x-api-key header")
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["auth_config_id"] != "ac_1" || body["user_id"] != "u1" || body["callback_url"] != "https://example.test/tiktok" {
			t.Errorf("unexpected body: %v", body)
		}
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"redirect_url":"https://connect.composio.dev/link/x","connected_account_id":"ca_1","link_token":"t"}`))
	})

	link, err := c.CreateConnectLink(context.Background(), "ac_1", "u1", "https://example.test/tiktok")
	if err != nil {
		t.Fatal(err)
	}
	if link.RedirectURL != "https://connect.composio.dev/link/x" || link.ConnectedAccountID != "ca_1" {
		t.Fatalf("unexpected link: %+v", link)
	}
}

func TestCreateConnectLinkAPIError(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"message":"auth config not found"}`))
	})
	_, err := c.CreateConnectLink(context.Background(), "nope", "u1", "")
	var apiErr *APIError
	if !errors.As(err, &apiErr) || apiErr.Status != http.StatusNotFound {
		t.Fatalf("expected *APIError 404, got %v", err)
	}
}

func TestFindActiveAccount(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("user_ids") != "u1" || q.Get("toolkit_slugs") != "tiktok" || q.Get("statuses") != "ACTIVE" {
			t.Errorf("unexpected query: %v", q)
		}
		// Includes entries a misbehaving server-side filter might leak:
		// another user, another toolkit, a non-active status.
		_, _ = w.Write([]byte(`{"items":[
			{"id":"old","status":"ACTIVE","user_id":"u1","created_at":"2026-01-01T00:00:00Z","toolkit":{"slug":"tiktok"}},
			{"id":"new","status":"ACTIVE","user_id":"u1","created_at":"2026-02-01T00:00:00Z","toolkit":{"slug":"tiktok"}},
			{"id":"other-user","status":"ACTIVE","user_id":"u2","created_at":"2027-01-01T00:00:00Z","toolkit":{"slug":"tiktok"}},
			{"id":"other-toolkit","status":"ACTIVE","user_id":"u1","created_at":"2027-01-01T00:00:00Z","toolkit":{"slug":"youtube"}},
			{"id":"inactive","status":"INITIATED","user_id":"u1","created_at":"2027-01-01T00:00:00Z","toolkit":{"slug":"tiktok"}}
		]}`))
	})

	acc, err := c.FindActiveAccount(context.Background(), "u1", "tiktok")
	if err != nil {
		t.Fatal(err)
	}
	if acc == nil || acc.ID != "new" {
		t.Fatalf("expected newest matching account 'new', got %+v", acc)
	}
}

func TestFindActiveAccountNone(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"items":[]}`))
	})
	acc, err := c.FindActiveAccount(context.Background(), "u1", "tiktok")
	if err != nil || acc != nil {
		t.Fatalf("expected (nil, nil), got (%+v, %v)", acc, err)
	}
}

func TestExecuteTool(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/tools/execute/TIKTOK_GET_USER_STATS" {
			t.Errorf("unexpected path %s", r.URL.Path)
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["connected_account_id"] != "ca_1" || body["user_id"] != "u1" {
			t.Errorf("unexpected body: %v", body)
		}
		_, _ = w.Write([]byte(`{"data":{"user":{"display_name":"Marcel"}},"error":null,"successful":true}`))
	})

	data, err := c.ExecuteTool(context.Background(), "TIKTOK_GET_USER_STATS", "u1", "ca_1", map[string]any{"fields": []string{"display_name"}})
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"user":{"display_name":"Marcel"}}` {
		t.Fatalf("unexpected data: %s", data)
	}
}

func TestExecuteToolUnsuccessful(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"data":{},"error":"missing scope user.info.basic","successful":false}`))
	})
	_, err := c.ExecuteTool(context.Background(), "TIKTOK_GET_USER_STATS", "u1", "ca_1", nil)
	if err == nil {
		t.Fatal("expected error for successful=false")
	}
}

func TestDeleteAccount(t *testing.T) {
	called := false
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		called = true
		if r.Method != http.MethodDelete || r.URL.Path != "/connected_accounts/ca_1" {
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		if r.URL.Query().Get("revoke_on_delete") != "true" {
			t.Errorf("expected revoke_on_delete=true")
		}
		_, _ = w.Write([]byte(`{"success":true}`))
	})
	if err := c.DeleteAccount(context.Background(), "ca_1", true); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("server never called")
	}
}
