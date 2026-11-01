package main

import "testing"

func TestResolveActivityWebhookSecret(t *testing.T) {
	cases := []struct {
		name, activity, subathon string
		want                     string
		wantFallback             bool
	}{
		{"dedicated secret wins", "activity-secret-1", "subathon-secret-1", "activity-secret-1", false},
		{"falls back to the subathon secret", "", "subathon-secret-1", "subathon-secret-1", true},
		{"nothing configured", "", "", "", false},
		{"only a dedicated secret", "activity-secret-1", "", "activity-secret-1", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, fb := resolveActivityWebhookSecret(c.activity, c.subathon)
			if got != c.want || fb != c.wantFallback {
				t.Fatalf("got (%q, %v), want (%q, %v)", got, fb, c.want, c.wantFallback)
			}
		})
	}
}
