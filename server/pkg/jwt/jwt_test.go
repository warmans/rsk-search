package jwt

import (
	"strings"
	"testing"
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr string
	}{
		{name: "empty key is rejected", key: "", wantErr: "not set"},
		{name: "historic default key is rejected", key: "insecure", wantErr: "known default"},
		{name: "short key is rejected", key: "short", wantErr: "at least 32 bytes"},
		{name: "long random key is accepted", key: "a-very-long-random-key-that-is-definitely-not-guessable", wantErr: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{SigningKey: tt.key}
			err := cfg.Validate()
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("Validate() unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatal("Validate() expected an error but got none")
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Validate() error %q does not contain %q", err.Error(), tt.wantErr)
			}
		})
	}
}
