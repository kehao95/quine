package llm

import "strings"

type endpointQuirks struct {
	noAssistantPrefill      bool
	stripAssistantReasoning bool
	thinkingBudgetFallback  func(model string, hasTools bool) bool
	trimV1Prefix            bool
}

func (q endpointQuirks) requiresThinkingBudgetFallback(model string, hasTools bool) bool {
	if q.thinkingBudgetFallback == nil {
		return false
	}
	return q.thinkingBudgetFallback(model, hasTools)
}

func quirksFor(apiBase, apiKey string) endpointQuirks {
	base := strings.ToLower(strings.TrimRight(strings.TrimSpace(apiBase), "/"))

	var quirks endpointQuirks

	// OpenAI-compatible custom bases may already include the API version or
	// Gemini's OpenAI compatibility prefix, so protocol paths must not add
	// another /v1 segment. Zhipu/z.ai expose their OpenAI-compatible surface at
	// .../paas/v4 (e.g. https://api.z.ai/api/coding/paas/v4) and expect
	// /chat/completions, not /v1/chat/completions.
	if strings.HasSuffix(base, "/v1") || strings.HasSuffix(base, "/openai") ||
		strings.HasSuffix(base, "/paas/v4") {
		quirks.trimV1Prefix = true
	}

	// Some providers reject trailing assistant-prefill messages and require the
	// conversation to end on a user turn.
	if strings.Contains(base, "api.githubcopilot.com") ||
		strings.Contains(base, "api.anthropic.com") ||
		apiKey == "copilot-oauth" ||
		apiKey == "claude-oauth" {
		quirks.noAssistantPrefill = true
	}

	// z.ai (Zhipu GLM) rejects requests whose assistant turns echo prior
	// reasoning_content back into the conversation (HTTP 400 code 1214). The
	// runtime's memory manager only trims visible content, so re-sent reasoning
	// accumulates untrimmed until the payload is rejected. Strip the outgoing
	// echo for this provider; response parsing still stores reasoning in the
	// tape. DeepSeek, by contrast, requires reasoning on assistant messages, so
	// this stays provider-scoped rather than a global drop.
	if strings.Contains(base, "z.ai") {
		quirks.stripAssistantReasoning = true
	}

	// GitHub Copilot Business chat completions currently rejects gpt-5.4 tool
	// requests when a thinking budget is sent. Responses API traffic keeps the
	// normal budget; NewProvider drops this quirk for non-chat protocols.
	if strings.Contains(base, "api.business.githubcopilot.com") {
		quirks.thinkingBudgetFallback = func(model string, hasTools bool) bool {
			if !hasTools {
				return false
			}
			return strings.HasPrefix(strings.ToLower(strings.TrimSpace(model)), "gpt-5.4")
		}
	}

	return quirks
}
