package update

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetLatestHash_Success(t *testing.T) {
	jsonResp := `[ { "name": "devctl-darwin-amd64-abc123" }, { "name": "devctl-darwin-amd64-def456" }, { "name": "devctl-linux-amd64-zzz999" } ]`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, jsonResp)
	}))
	defer ts.Close()

	hash, err := getLatestHash(ts.URL, "darwin", "amd64")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hash != "def456" {
		t.Errorf("expected def456, got %s", hash)
	}
}

func TestGetLatestHash_NoMatch(t *testing.T) {
	jsonResp := `[ { "name": "devctl-linux-amd64-zzz999" } ]`
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, jsonResp)
	}))
	defer ts.Close()

	_, err := getLatestHash(ts.URL, "darwin", "amd64")
	if err == nil || !strings.Contains(err.Error(), "no binary found") {
		t.Errorf("expected no binary found error, got %v", err)
	}
}

func TestGetLatestHash_HTTPError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
	}))
	defer ts.Close()

	_, err := getLatestHash(ts.URL, "darwin", "amd64")
	if err == nil || !strings.Contains(err.Error(), "GitHub API returned status 500") {
		t.Errorf("expected status 500 error, got %v", err)
	}
}

func TestGetLatestHash_BadJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "not json")
	}))
	defer ts.Close()

	_, err := getLatestHash(ts.URL, "darwin", "amd64")
	if err == nil || !strings.Contains(err.Error(), "failed to decode") {
		t.Errorf("expected decode error, got %v", err)
	}
}
