package server

import (
	"testing"

	x402 "github.com/coinbase/x402/go"
)

func TestBuildRoutes_ReturnsAllEndpoints(t *testing.T) {
	routes := BuildRoutes("0x1234", x402.Network("eip155:84532"))

	expectedRoutes := []string{"GET /weather", "GET /joke", "GET /premium-data"}
	if len(routes) != len(expectedRoutes) {
		t.Fatalf("expected %d routes, got %d", len(expectedRoutes), len(routes))
	}

	for _, name := range expectedRoutes {
		route, ok := routes[name]
		if !ok {
			t.Errorf("missing route %q", name)
			continue
		}
		if len(route.Accepts) == 0 {
			t.Errorf("route %q has no accepts", name)
		}
		if route.Accepts[0].Scheme != "exact" {
			t.Errorf("route %q: expected scheme 'exact', got %q", name, route.Accepts[0].Scheme)
		}
		if route.Accepts[0].PayTo != "0x1234" {
			t.Errorf("route %q: expected payTo '0x1234', got %q", name, route.Accepts[0].PayTo)
		}
		if string(route.Accepts[0].Network) != "eip155:84532" {
			t.Errorf("route %q: unexpected network %q", name, route.Accepts[0].Network)
		}
		if route.MimeType != "application/json" {
			t.Errorf("route %q: expected mimeType 'application/json', got %q", name, route.MimeType)
		}
	}
}

func TestBuildRoutes_Description(t *testing.T) {
	routes := BuildRoutes("0xABCD", x402.Network("eip155:1"))

	for name, route := range routes {
		if route.Description == "" {
			t.Errorf("route %q has empty description", name)
		}
	}
}
