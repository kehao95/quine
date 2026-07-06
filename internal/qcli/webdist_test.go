package qcli

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func newWebKernel(t *testing.T, webDist string) *Kernel {
	t.Helper()
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "pid"), 0o755); err != nil {
		t.Fatal(err)
	}
	kernel, err := NewKernel(context.Background(), KernelOptions{RuntimeRoot: root, WebDist: webDist})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(kernel.Close)
	return kernel
}

func get(t *testing.T, srv *httptest.Server, path string) (int, string, string) {
	t.Helper()
	res, err := http.Get(srv.URL + path)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	return res.StatusCode, string(body), res.Header.Get("Content-Type")
}

// WEB.md §6: with a dist dir, / serves index.html, unknown paths SPA-fall back
// to index.html, real assets serve as themselves, and API routes stay
// unshadowed. Without dist, / serves the plain-text endpoint index.
func TestServeWebDist(t *testing.T) {
	dist := t.TempDir()
	if err := os.WriteFile(filepath.Join(dist, "index.html"), []byte("<!doctype html><title>qcli</title>"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dist, "assets"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dist, "assets", "app.js"), []byte("console.log(1)"), 0o644); err != nil {
		t.Fatal(err)
	}
	kernel := newWebKernel(t, dist)
	srv := httptest.NewServer(kernel.Handler())
	defer srv.Close()

	if code, body, _ := get(t, srv, "/"); code != 200 || !strings.Contains(body, "<title>qcli</title>") {
		t.Fatalf("GET / = %d %q; want index.html", code, body)
	}
	if code, body, _ := get(t, srv, "/peers/some/spa/route"); code != 200 || !strings.Contains(body, "<title>qcli</title>") {
		t.Fatalf("GET spa route = %d %q; want index.html fallback", code, body)
	}
	if code, body, _ := get(t, srv, "/assets/app.js"); code != 200 || body != "console.log(1)" {
		t.Fatalf("GET asset = %d %q; want file content", code, body)
	}
	if code, body, ct := get(t, srv, "/healthz"); code != 200 || !strings.Contains(body, `"ok":true`) || !strings.Contains(ct, "application/json") {
		t.Fatalf("GET /healthz = %d %q ct=%q; want unshadowed JSON", code, body, ct)
	}
	// /events with no peer must still hit the API handler (409), not the SPA.
	if code, body, _ := get(t, srv, "/events"); code != http.StatusConflict || !strings.Contains(body, "no_peer") {
		t.Fatalf("GET /events = %d %q; want 409 no_peer", code, body)
	}
}

func TestServeWebDistTraversalContained(t *testing.T) {
	dist := t.TempDir()
	if err := os.WriteFile(filepath.Join(dist, "index.html"), []byte("shell"), 0o644); err != nil {
		t.Fatal(err)
	}
	secret := filepath.Join(filepath.Dir(dist), "secret.txt")
	if err := os.WriteFile(secret, []byte("nope"), 0o644); err != nil {
		t.Fatal(err)
	}
	kernel := newWebKernel(t, dist)
	srv := httptest.NewServer(kernel.Handler())
	defer srv.Close()
	if _, body, _ := get(t, srv, "/../secret.txt"); strings.Contains(body, "nope") {
		t.Fatalf("path traversal escaped dist: %q", body)
	}
	if _, body, _ := get(t, srv, "/%2e%2e/secret.txt"); strings.Contains(body, "nope") {
		t.Fatalf("encoded traversal escaped dist: %q", body)
	}
}

func TestServeNoDistFallsBackToEndpointIndex(t *testing.T) {
	// Force discovery to fail: point QCLI_WEB_DIST at a dir with no index.html.
	t.Setenv("QCLI_WEB_DIST", t.TempDir())
	kernel := newWebKernel(t, "")
	srv := httptest.NewServer(kernel.Handler())
	defer srv.Close()
	code, body, ct := get(t, srv, "/")
	if code != 200 || !strings.Contains(body, "/events") || !strings.Contains(ct, "text/plain") {
		t.Fatalf("GET / without dist = %d ct=%q %q; want plain-text endpoint index", code, ct, body)
	}
	if code, _, _ := get(t, srv, "/nonexistent"); code != http.StatusNotFound {
		t.Fatalf("GET unknown path without dist = %d; want 404", code)
	}
}
