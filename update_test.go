package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVersionGreater(t *testing.T) {
	tests := []struct {
		remote, local string
		want          bool
	}{
		{"1.2.0", "1.1.0", true},
		{"v1.2.0", "1.1.0", true},
		{"1.1.0", "1.1.0", false},
		{"v1.1.0", "v1.1.0", false},
		{"1.1.0", "1.2.0", false},
		{"1.10.0", "1.9.0", true},
		{"2.0.0", "1.9.9", true},
		{"1.1.1", "1.1.0", true},
		{"1.1", "1.0.9", true},
		{"1.0.0", "1.0", false},
		{"v1.2.0-beta", "1.1.0", true},
	}
	for _, tt := range tests {
		if got := versionGreater(tt.remote, tt.local); got != tt.want {
			t.Errorf("versionGreater(%q, %q) = %v, want %v", tt.remote, tt.local, got, tt.want)
		}
	}
}

func TestTagFromReleaseURL(t *testing.T) {
	tests := []struct {
		url  string
		want string
		ok   bool
	}{
		{"https://github.com/RykivSale/faceit-voicechat/releases/tag/v1.1.0", "1.1.0", true},
		{"/RykivSale/faceit-voicechat/releases/tag/v1.2.0", "1.2.0", true},
		{"https://github.com/RykivSale/faceit-voicechat/releases/tag/1.3.0?foo=1", "1.3.0", true},
		{"https://github.com/RykivSale/faceit-voicechat/releases/latest", "", false},
		{"https://example.com/not-a-release", "", false},
	}
	for _, tt := range tests {
		got, err := tagFromReleaseURL(tt.url)
		if tt.ok {
			if err != nil {
				t.Errorf("tagFromReleaseURL(%q) unexpected error: %v", tt.url, err)
				continue
			}
			if got != tt.want {
				t.Errorf("tagFromReleaseURL(%q) = %q, want %q", tt.url, got, tt.want)
			}
		} else if err == nil {
			t.Errorf("tagFromReleaseURL(%q) = %q, want error", tt.url, got)
		}
	}
}

func TestFetchLatestFromRedirect(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://github.com/RykivSale/faceit-voicechat/releases/tag/v9.9.9", http.StatusFound)
	}))
	defer srv.Close()

	old := githubLatestURL
	githubLatestURL = srv.URL
	defer func() { githubLatestURL = old }()

	got, err := fetchLatestFromRedirect()
	if err != nil {
		t.Fatalf("fetchLatestFromRedirect: %v", err)
	}
	if got != "9.9.9" {
		t.Fatalf("got %q, want 9.9.9", got)
	}
}

func TestFetchLatestFromAPI(t *testing.T) {
	payload, _ := json.Marshal(map[string]string{"tag_name": "v3.2.1"})
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(payload)
	}))
	defer srv.Close()

	old := githubAPIURL
	githubAPIURL = srv.URL
	defer func() { githubAPIURL = old }()

	got, err := fetchLatestFromAPI()
	if err != nil {
		t.Fatalf("fetchLatestFromAPI: %v", err)
	}
	if got != "3.2.1" {
		t.Fatalf("got %q, want 3.2.1", got)
	}
}

func TestFetchLatestVersionFallsBackToAPI(t *testing.T) {
	redir := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "nope", http.StatusInternalServerError)
	}))
	defer redir.Close()

	payload, _ := json.Marshal(map[string]string{"tag_name": "v4.0.0"})
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(payload)
	}))
	defer api.Close()

	oldLatest, oldAPI := githubLatestURL, githubAPIURL
	githubLatestURL = redir.URL
	githubAPIURL = api.URL
	defer func() {
		githubLatestURL = oldLatest
		githubAPIURL = oldAPI
	}()

	got, err := fetchLatestVersion()
	if err != nil {
		t.Fatalf("fetchLatestVersion: %v", err)
	}
	if got != "4.0.0" {
		t.Fatalf("got %q, want 4.0.0", got)
	}
}

func TestParseConfigCheckUpdates(t *testing.T) {
	t.Run("missing field defaults to true", func(t *testing.T) {
		cfg := parseConfig([]byte(`{"gameFolder":"C:\\cs2","keys":["F5","F6","F7"]}`))
		if !cfg.CheckUpdates {
			t.Fatal("expected CheckUpdates to default to true")
		}
	})
	t.Run("explicit false is kept", func(t *testing.T) {
		cfg := parseConfig([]byte(`{"checkUpdates":false}`))
		if cfg.CheckUpdates {
			t.Fatal("expected CheckUpdates false to be kept")
		}
	})
	t.Run("invalid json falls back to defaults", func(t *testing.T) {
		cfg := parseConfig([]byte(`{not json`))
		if !cfg.CheckUpdates || len(cfg.Keys) != 3 {
			t.Fatalf("unexpected defaults: %+v", cfg)
		}
	})
}
