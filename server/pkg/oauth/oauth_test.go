package oauth

import (
	"testing"
)

func TestSanitizeReturnURL(t *testing.T) {
	cfg := &Config{ReturnURL: "https://scrimpton.com/oauth/%s/return"}

	tests := []struct {
		name    string
		rawURL  string
		want    string
	}{
		{name: "same origin is kept", rawURL: "https://scrimpton.com/search", want: "https://scrimpton.com/search"},
		{name: "same origin with path is kept", rawURL: "https://scrimpton.com/episode/xfm-s01e01", want: "https://scrimpton.com/episode/xfm-s01e01"},
		{name: "attacker origin is rejected", rawURL: "https://evil.com/search", want: ""},
		{name: "attacker subdomain is rejected", rawURL: "https://scrimpton.com.evil.com/search", want: ""},
		{name: "scheme downgrade is rejected", rawURL: "http://scrimpton.com/search", want: ""},
		{name: "different port is rejected", rawURL: "https://scrimpton.com:8443/search", want: ""},
		{name: "empty url is rejected", rawURL: "", want: ""},
		{name: "garbage url is rejected", rawURL: "not a url", want: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cfg.SanitizeReturnURL(tt.rawURL); got != tt.want {
				t.Fatalf("SanitizeReturnURL(%q) = %q, want %q", tt.rawURL, got, tt.want)
			}
		})
	}
}

func TestSanitizeReturnURLDevDefault(t *testing.T) {
	// the default ReturnURL points at the dev GUI origin
	cfg := &Config{ReturnURL: "http://localhost:4200/oauth/%s/return"}
	if got := cfg.SanitizeReturnURL("http://localhost:4200/search"); got != "http://localhost:4200/search" {
		t.Fatalf("dev origin rejected: %q", got)
	}
	if got := cfg.SanitizeReturnURL("https://scrimpton.com/search"); got != "" {
		t.Fatalf("prod origin must not be accepted against dev config: %q", got)
	}
}
