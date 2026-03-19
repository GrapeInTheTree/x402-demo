package config

import (
	"log/slog"
	"os"
	"testing"
)

func TestLoadFacilitator_MissingKey(t *testing.T) {
	os.Unsetenv("FACILITATOR_PRIVATE_KEY")
	_, err := LoadFacilitator()
	if err == nil {
		t.Fatal("expected error for missing FACILITATOR_PRIVATE_KEY")
	}
}

func TestLoadFacilitator_Defaults(t *testing.T) {
	t.Setenv("FACILITATOR_PRIVATE_KEY", "0xdeadbeef")
	os.Unsetenv("RPC_URL")
	os.Unsetenv("NETWORK")
	os.Unsetenv("FACILITATOR_PORT")

	cfg, err := LoadFacilitator()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.RPCURL != "https://sepolia.base.org" {
		t.Errorf("expected default RPC_URL, got %s", cfg.RPCURL)
	}
	if cfg.Network != "eip155:84532" {
		t.Errorf("expected default network, got %s", cfg.Network)
	}
	if cfg.Port != "4022" {
		t.Errorf("expected default port 4022, got %s", cfg.Port)
	}
}

func TestLoadResource_MissingRequired(t *testing.T) {
	os.Unsetenv("FACILITATOR_URL")
	os.Unsetenv("PAY_TO_ADDRESS")
	_, err := LoadResource()
	if err == nil {
		t.Fatal("expected error for missing required fields")
	}
}

func TestLoadResource_Valid(t *testing.T) {
	t.Setenv("FACILITATOR_URL", "http://localhost:4022")
	t.Setenv("PAY_TO_ADDRESS", "0x1234")

	cfg, err := LoadResource()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.FacilitatorURL != "http://localhost:4022" {
		t.Errorf("unexpected FACILITATOR_URL: %s", cfg.FacilitatorURL)
	}
}

func TestLoadClient_MissingRequired(t *testing.T) {
	os.Unsetenv("CLIENT_PRIVATE_KEY")
	os.Unsetenv("RESOURCE_URL")
	_, err := LoadClient()
	if err == nil {
		t.Fatal("expected error for missing required fields")
	}
}

func TestLoadClient_Defaults(t *testing.T) {
	t.Setenv("CLIENT_PRIVATE_KEY", "0xdeadbeef")
	t.Setenv("RESOURCE_URL", "http://localhost:4021")
	os.Unsetenv("ENDPOINT_PATH")

	cfg, err := LoadClient()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.EndpointPath != "/weather" {
		t.Errorf("expected default endpoint /weather, got %s", cfg.EndpointPath)
	}
}

func TestParseLogLevel(t *testing.T) {
	tests := []struct {
		input string
		want  slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"DEBUG", slog.LevelDebug},
		{"warn", slog.LevelWarn},
		{"warning", slog.LevelWarn},
		{"error", slog.LevelError},
		{"info", slog.LevelInfo},
		{"", slog.LevelInfo},
		{"unknown", slog.LevelInfo},
	}
	for _, tt := range tests {
		if got := parseLogLevel(tt.input); got != tt.want {
			t.Errorf("parseLogLevel(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
