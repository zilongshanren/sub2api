package responseheaders

import (
	"net/http"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
)

func TestFilterHeadersDisabledUsesDefaultAllowlist(t *testing.T) {
	src := http.Header{}
	src.Add("Content-Type", "application/json")
	src.Add("X-Request-Id", "req-123")
	src.Add("X-Test", "ok")
	src.Add("Connection", "keep-alive")
	src.Add("Content-Length", "123")

	cfg := config.ResponseHeaderConfig{
		Enabled:     false,
		ForceRemove: []string{"x-request-id"},
	}

	filtered := FilterHeaders(src, CompileHeaderFilter(cfg))
	if filtered.Get("Content-Type") != "application/json" {
		t.Fatalf("expected Content-Type passthrough, got %q", filtered.Get("Content-Type"))
	}
	if filtered.Get("X-Request-Id") != "req-123" {
		t.Fatalf("expected X-Request-Id allowed, got %q", filtered.Get("X-Request-Id"))
	}
	if filtered.Get("X-Test") != "" {
		t.Fatalf("expected X-Test removed, got %q", filtered.Get("X-Test"))
	}
	if filtered.Get("Connection") != "" {
		t.Fatalf("expected Connection to be removed, got %q", filtered.Get("Connection"))
	}
	if filtered.Get("Content-Length") != "" {
		t.Fatalf("expected Content-Length to be removed, got %q", filtered.Get("Content-Length"))
	}
}

func TestFilterHeadersEnabledUsesAllowlist(t *testing.T) {
	src := http.Header{}
	src.Add("Content-Type", "application/json")
	src.Add("X-Extra", "ok")
	src.Add("X-Remove", "nope")
	src.Add("X-Blocked", "nope")

	cfg := config.ResponseHeaderConfig{
		Enabled:           true,
		AdditionalAllowed: []string{"x-extra"},
		ForceRemove:       []string{"x-remove"},
	}

	filtered := FilterHeaders(src, CompileHeaderFilter(cfg))
	if filtered.Get("Content-Type") != "application/json" {
		t.Fatalf("expected Content-Type allowed, got %q", filtered.Get("Content-Type"))
	}
	if filtered.Get("X-Extra") != "ok" {
		t.Fatalf("expected X-Extra allowed, got %q", filtered.Get("X-Extra"))
	}
	if filtered.Get("X-Remove") != "" {
		t.Fatalf("expected X-Remove removed, got %q", filtered.Get("X-Remove"))
	}
	if filtered.Get("X-Blocked") != "" {
		t.Fatalf("expected X-Blocked removed, got %q", filtered.Get("X-Blocked"))
	}
}

func TestFilterHeadersAnthropicUnifiedRateLimitPassthrough(t *testing.T) {
	src := http.Header{}
	src.Add("anthropic-ratelimit-unified-5h-status", "ok")
	src.Add("anthropic-ratelimit-unified-5h-utilization", "0.42")
	src.Add("anthropic-ratelimit-unified-5h-reset", "1234567890")
	src.Add("anthropic-ratelimit-unified-5h-surpassed-threshold", "false")
	src.Add("anthropic-ratelimit-unified-7d-status", "ok")
	src.Add("anthropic-ratelimit-unified-7d-utilization", "0.13")
	src.Add("anthropic-ratelimit-unified-7d-reset", "1234567890")
	src.Add("anthropic-ratelimit-unified-fallback-status", "ok")
	src.Add("anthropic-ratelimit-unified-fallback-reset", "1234567890")
	src.Add("anthropic-ratelimit-unified-reset", "1234567890")

	filtered := FilterHeaders(src, CompileHeaderFilter(config.ResponseHeaderConfig{}))
	for key, want := range map[string]string{
		"anthropic-ratelimit-unified-5h-status":              "ok",
		"anthropic-ratelimit-unified-5h-utilization":         "0.42",
		"anthropic-ratelimit-unified-5h-reset":               "1234567890",
		"anthropic-ratelimit-unified-5h-surpassed-threshold": "false",
		"anthropic-ratelimit-unified-7d-status":              "ok",
		"anthropic-ratelimit-unified-7d-utilization":         "0.13",
		"anthropic-ratelimit-unified-7d-reset":               "1234567890",
		"anthropic-ratelimit-unified-fallback-status":        "ok",
		"anthropic-ratelimit-unified-fallback-reset":         "1234567890",
		"anthropic-ratelimit-unified-reset":                  "1234567890",
	} {
		if got := filtered.Get(key); got != want {
			t.Fatalf("expected %s=%q passthrough, got %q", key, want, got)
		}
	}
}

func TestFilterHeadersCodexRateLimitPassthrough(t *testing.T) {
	src := http.Header{}
	src.Add("x-codex-primary-used-percent", "42")
	src.Add("x-codex-primary-reset-after-seconds", "384607")
	src.Add("x-codex-primary-window-minutes", "10080")
	src.Add("x-codex-secondary-used-percent", "3")
	src.Add("x-codex-secondary-reset-after-seconds", "17369")
	src.Add("x-codex-secondary-window-minutes", "300")

	filtered := FilterHeaders(src, CompileHeaderFilter(config.ResponseHeaderConfig{}))
	for key, want := range map[string]string{
		"x-codex-primary-used-percent":          "42",
		"x-codex-primary-reset-after-seconds":   "384607",
		"x-codex-primary-window-minutes":        "10080",
		"x-codex-secondary-used-percent":        "3",
		"x-codex-secondary-reset-after-seconds": "17369",
		"x-codex-secondary-window-minutes":      "300",
	} {
		if got := filtered.Get(key); got != want {
			t.Fatalf("expected %s=%q passthrough, got %q", key, want, got)
		}
	}
}
