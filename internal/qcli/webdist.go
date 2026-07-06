package qcli

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// DiscoverWebDist resolves the built web client directory (WEB.md §3):
// QCLI_WEB_DIST wins when it holds an index.html; otherwise walk up from the
// working directory and the executable directory looking for
// operator/qcli/web/dist — the same discovery pattern the TUI launcher uses.
// Empty string means no web app is available (the kernel then serves a
// plain-text endpoint index at /).
func DiscoverWebDist() string {
	if env := strings.TrimSpace(os.Getenv("QCLI_WEB_DIST")); env != "" {
		if hasIndexHTML(env) {
			return env
		}
		return ""
	}
	var candidates []string
	if wd, err := os.Getwd(); err == nil {
		for dir := wd; ; dir = filepath.Dir(dir) {
			candidates = append(candidates, filepath.Join(dir, "operator", "qcli", "web", "dist"))
			if parent := filepath.Dir(dir); parent == dir {
				break
			}
		}
	}
	if exe, err := os.Executable(); err == nil {
		dir := filepath.Dir(exe)
		candidates = append(candidates,
			filepath.Join(dir, "operator", "qcli", "web", "dist"),
			filepath.Join(dir, "..", "operator", "qcli", "web", "dist"),
		)
	}
	for _, dir := range candidates {
		if hasIndexHTML(dir) {
			return dir
		}
	}
	return ""
}

func hasIndexHTML(dir string) bool {
	info, err := os.Stat(filepath.Join(dir, "index.html"))
	return err == nil && !info.IsDir()
}

// handleRoot serves the built web client under / (SPA fallback to index.html)
// when a dist directory is available, and a plain-text endpoint index
// otherwise. API routes are registered as exact mux patterns, so they always
// take precedence over this catch-all (WEB.md §3).
func (k *Kernel) handleRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	dist := k.webDist
	if dist == "" {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprintf(w, "qcli/1 kernel\n\nendpoints:\n  GET  /events\n  GET  /context\n  GET  /status\n  GET  /roster\n  GET  /peer-contract\n  POST /command\n  GET  /healthz\n  GET  /history (reserved)\n\nweb client: build with `make build-qcli-web` (operator/qcli/web/dist), then reload.\n")
		return
	}
	rel := strings.TrimPrefix(filepath.Clean("/"+r.URL.Path), "/")
	full := filepath.Join(dist, filepath.FromSlash(rel))
	if within, err := filepath.Rel(dist, full); err != nil || strings.HasPrefix(within, "..") {
		http.NotFound(w, r)
		return
	}
	if info, err := os.Stat(full); err == nil && !info.IsDir() {
		http.ServeFile(w, r, full)
		return
	}
	// SPA fallback: any unknown non-API path renders the app shell.
	http.ServeFile(w, r, filepath.Join(dist, "index.html"))
}
