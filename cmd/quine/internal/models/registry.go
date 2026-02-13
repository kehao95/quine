// Package models provides access to the models.dev registry for model metadata.
package models

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const registryURL = "https://models.dev/api.json"

// Provider represents a model provider from models.dev
type Provider struct {
	ID     string            `json:"id"`
	Name   string            `json:"name"`
	Env    []string          `json:"env"`
	API    string            `json:"api"` // API base URL
	Doc    string            `json:"doc"`
	Models map[string]*Model `json:"models"`
}

// Model represents a single model's metadata
type Model struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Family      string `json:"family"`
	ToolCall    bool   `json:"tool_call"`
	Reasoning   bool   `json:"reasoning"`
	Attachment  bool   `json:"attachment"`
	Temperature bool   `json:"temperature"`
	OpenWeights bool   `json:"open_weights"`
	Limit       struct {
		Context int `json:"context"`
		Output  int `json:"output"`
	} `json:"limit"`
	Cost struct {
		Input      float64 `json:"input"`
		Output     float64 `json:"output"`
		CacheRead  float64 `json:"cache_read"`
		CacheWrite float64 `json:"cache_write"`
	} `json:"cost"`
}

// Registry holds the parsed models.dev data
type Registry struct {
	providers map[string]*Provider
	mu        sync.RWMutex
}

var (
	globalRegistry *Registry
	once           sync.Once
	initErr        error
)

// Load fetches and parses the models.dev registry.
// Results are cached globally for the process lifetime.
func Load() (*Registry, error) {
	once.Do(func() {
		globalRegistry, initErr = fetch()
	})
	return globalRegistry, initErr
}

// fetch retrieves the registry from models.dev
func fetch() (*Registry, error) {
	client := &http.Client{Timeout: 30 * time.Second}

	resp, err := client.Get(registryURL)
	if err != nil {
		return nil, fmt.Errorf("fetching models registry: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("models registry returned HTTP %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading registry response: %w", err)
	}

	var providers map[string]*Provider
	if err := json.Unmarshal(body, &providers); err != nil {
		return nil, fmt.Errorf("parsing registry JSON: %w", err)
	}

	return &Registry{providers: providers}, nil
}

// LookupResult contains the result of looking up a model
type LookupResult struct {
	ProviderID    string   // e.g. "anthropic", "openai", "google", "amazon-bedrock", "openrouter"
	ProviderName  string   // e.g. "Anthropic", "OpenAI", "OpenRouter"
	Model         *Model   // Full model metadata
	ActualModelID string   // Model ID to use in API calls (may differ from input when provider prefix is used)
	EnvVars       []string // Environment variables for API key
	APIKey        string   // Resolved API key value
	APIBase       string   // API base URL from registry (e.g. "https://openrouter.ai/api/v1")
	ContextWindow int      // Context window size
}

// Lookup finds a model by ID and returns its provider info and API key.
// It supports two formats:
//   - "provider/model" - explicit provider selection (e.g., "openrouter/anthropic/claude-3.5-sonnet")
//   - "model" - auto-detect provider, preferring primary providers
//
// When a model exists in multiple providers, it prefers the primary providers
// (anthropic, openai, google, amazon-bedrock, openrouter) over third-party services.
func (r *Registry) Lookup(modelID string) (*LookupResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Check for explicit provider prefix: "provider/model..."
	if idx := strings.Index(modelID, "/"); idx > 0 {
		providerID := modelID[:idx]
		actualModelID := modelID[idx+1:]

		// Check if providerID is a known provider
		if provider, ok := r.providers[providerID]; ok {
			if model, ok := provider.Models[actualModelID]; ok {
				return r.buildResult(providerID, provider, model), nil
			}
			return nil, fmt.Errorf("model %q not found in provider %q", actualModelID, providerID)
		}
		// Not a known provider prefix, treat entire string as model ID
	}

	// Primary providers to prefer (in order)
	primaryProviders := []string{"anthropic", "openai", "google", "amazon-bedrock", "openrouter"}

	// First, try primary providers
	for _, providerID := range primaryProviders {
		if provider, ok := r.providers[providerID]; ok {
			if model, ok := provider.Models[modelID]; ok {
				return r.buildResult(providerID, provider, model), nil
			}
		}
	}

	// Fall back to any provider that has this model
	for providerID, provider := range r.providers {
		if model, ok := provider.Models[modelID]; ok {
			return r.buildResult(providerID, provider, model), nil
		}
	}

	return nil, fmt.Errorf("model %q not found in registry", modelID)
}

// buildResult constructs a LookupResult from provider and model data
func (r *Registry) buildResult(providerID string, provider *Provider, model *Model) *LookupResult {
	result := &LookupResult{
		ProviderID:    providerID,
		ProviderName:  provider.Name,
		Model:         model,
		ActualModelID: model.ID,
		EnvVars:       provider.Env,
		APIBase:       provider.API,
		ContextWindow: model.Limit.Context,
	}

	// Try to resolve API key from environment
	for _, envVar := range provider.Env {
		if val := os.Getenv(envVar); val != "" {
			result.APIKey = val
			break
		}
	}

	return result
}

// ProviderMapping maps models.dev provider IDs to internal provider names
var ProviderMapping = map[string]string{
	"anthropic":      "anthropic",
	"openai":         "openai",
	"google":         "google",
	"amazon-bedrock": "bedrock",
	"openrouter":     "openai", // OpenRouter uses OpenAI-compatible API
}

// InternalProvider returns the internal provider name for routing
func (r *LookupResult) InternalProvider() string {
	if mapped, ok := ProviderMapping[r.ProviderID]; ok {
		return mapped
	}
	return r.ProviderID
}
