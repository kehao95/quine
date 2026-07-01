package transport

import (
	"net/http"
	"testing"

	"github.com/kehao95/quine/internal/llm/copilotoauth"
)

func TestCopilotOAuthTransportSignAddsCopilotHeaders(t *testing.T) {
	t.Setenv("GITHUB_COPILOT_INTEGRATION_ID", "copilot-developer-cli")

	tr := &CopilotOAuthTransport{
		token: copilotoauth.TokenState{
			AccessToken: "test-token",
		},
		userAgent:     "quine-test",
		integrationID: copilotIntegrationID(),
		interactionID: "interaction-test",
	}

	req, err := http.NewRequest(http.MethodPost, "https://api.business.githubcopilot.com/v1/chat/completions", nil)
	if err != nil {
		t.Fatalf("http.NewRequest() error = %v", err)
	}

	if err := tr.Sign(req, nil); err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	if got := req.URL.Path; got != "/chat/completions" {
		t.Fatalf("path = %q, want %q", got, "/chat/completions")
	}
	if got := req.Header.Get("Authorization"); got != "Bearer test-token" {
		t.Fatalf("Authorization = %q, want %q", got, "Bearer test-token")
	}
	if got := req.Header.Get("Openai-Intent"); got != "conversation-agent" {
		t.Fatalf("Openai-Intent = %q, want %q", got, "conversation-agent")
	}
	if got := req.Header.Get("X-Initiator"); got != "user" {
		t.Fatalf("X-Initiator = %q, want %q", got, "user")
	}
	if got := req.Header.Get("X-GitHub-Api-Version"); got != "2026-01-09" {
		t.Fatalf("X-GitHub-Api-Version = %q, want %q", got, "2026-01-09")
	}
	if got := req.Header.Get("Copilot-Integration-Id"); got != "copilot-developer-cli" {
		t.Fatalf("Copilot-Integration-Id = %q, want %q", got, "copilot-developer-cli")
	}
	if got := req.Header.Get("User-Agent"); got != "quine-test" {
		t.Fatalf("User-Agent = %q, want %q", got, "quine-test")
	}
	if got := req.Header.Get("X-Interaction-Id"); got != "interaction-test" {
		t.Fatalf("X-Interaction-Id = %q, want %q", got, "interaction-test")
	}
}
