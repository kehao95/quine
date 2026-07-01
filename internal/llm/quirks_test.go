package llm

import "testing"

func TestQuirksForEndpointRules(t *testing.T) {
	t.Run("copilot assistant prefill", func(t *testing.T) {
		q := quirksFor("https://api.githubcopilot.com", "")
		if !q.noAssistantPrefill {
			t.Fatal("expected GitHub Copilot endpoint to disable assistant prefill")
		}
	})

	t.Run("copilot oauth sentinel", func(t *testing.T) {
		q := quirksFor("https://example.invalid", "copilot-oauth")
		if !q.noAssistantPrefill {
			t.Fatal("expected copilot oauth sentinel to disable assistant prefill")
		}
	})

	t.Run("anthropic endpoint assistant prefill", func(t *testing.T) {
		q := quirksFor("https://api.anthropic.com", "")
		if !q.noAssistantPrefill {
			t.Fatal("expected anthropic endpoint to disable assistant prefill")
		}
	})

	t.Run("claude oauth sentinel assistant prefill", func(t *testing.T) {
		q := quirksFor("https://example.invalid", "claude-oauth")
		if !q.noAssistantPrefill {
			t.Fatal("expected claude oauth sentinel to disable assistant prefill")
		}
	})

	t.Run("openai compatible base path trim", func(t *testing.T) {
		for _, base := range []string{
			"https://api.moonshot.ai/v1",
			"https://generativelanguage.googleapis.com/v1beta/openai",
			"https://api.z.ai/api/coding/paas/v4",
			"https://api.z.ai/api/paas/v4",
		} {
			if q := quirksFor(base, ""); !q.trimV1Prefix {
				t.Fatalf("expected %q to trim protocol /v1 prefix", base)
			}
		}
	})

	t.Run("z.ai strips assistant reasoning echo", func(t *testing.T) {
		q := quirksFor("https://api.z.ai/api/coding/paas/v4", "")
		if !q.stripAssistantReasoning {
			t.Fatal("expected z.ai endpoint to strip echoed assistant reasoning_content")
		}
	})

	t.Run("non-z.ai endpoints keep reasoning echo", func(t *testing.T) {
		for _, base := range []string{
			"https://api.deepseek.com/v1",
			"https://api.moonshot.ai/v1",
			"https://api.openai.com/v1",
		} {
			if q := quirksFor(base, ""); q.stripAssistantReasoning {
				t.Fatalf("did not expect %q to strip assistant reasoning_content", base)
			}
		}
	})

	t.Run("business copilot gpt54 chat fallback", func(t *testing.T) {
		q := quirksFor("https://api.business.githubcopilot.com/v1", "")
		if !q.requiresThinkingBudgetFallback("gpt-5.4", true) {
			t.Fatal("expected gpt-5.4 tool call to require thinking budget fallback")
		}
		if q.requiresThinkingBudgetFallback("gpt-5.4", false) {
			t.Fatal("did not expect fallback without tools")
		}
		if q.requiresThinkingBudgetFallback("claude-opus-4.5", true) {
			t.Fatal("did not expect fallback for non-gpt-5.4 model")
		}
	})
}
