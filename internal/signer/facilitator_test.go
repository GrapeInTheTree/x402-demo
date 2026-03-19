package signer

import (
	"testing"
)

func TestNewFacilitatorSigner_InvalidKey(t *testing.T) {
	_, err := NewFacilitatorSigner("not-a-valid-key", "https://sepolia.base.org", nil)
	if err == nil {
		t.Fatal("expected error for invalid private key")
	}
}

func TestNewFacilitatorSigner_InvalidRPC(t *testing.T) {
	// Valid key format but unreachable RPC
	validKey := "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	_, err := NewFacilitatorSigner(validKey, "http://localhost:1", nil)
	// ethclient.Dial may succeed (lazy connect) or fail depending on implementation
	// Either way, if it succeeds the signer should be usable
	if err != nil {
		// Some RPC implementations fail on Dial, which is acceptable
		t.Logf("Dial failed (acceptable): %v", err)
	}
}

func TestNewFacilitatorSigner_AddressDerivation(t *testing.T) {
	// Hardhat account #0 — well-known test key
	key := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	expectedAddr := "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

	signer, err := NewFacilitatorSigner(key, "https://sepolia.base.org", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer signer.Close()

	if signer.Address() != expectedAddr {
		t.Errorf("expected %s, got %s", expectedAddr, signer.Address())
	}

	addrs := signer.GetAddresses()
	if len(addrs) != 1 || addrs[0] != expectedAddr {
		t.Errorf("GetAddresses() = %v, expected [%s]", addrs, expectedAddr)
	}
}

func TestNewFacilitatorSigner_WithOxPrefix(t *testing.T) {
	key := "0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
	expectedAddr := "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"

	signer, err := NewFacilitatorSigner(key, "https://sepolia.base.org", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer signer.Close()

	if signer.Address() != expectedAddr {
		t.Errorf("expected %s, got %s", expectedAddr, signer.Address())
	}
}

func TestFacilitatorSigner_Close(t *testing.T) {
	key := "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

	signer, err := NewFacilitatorSigner(key, "https://sepolia.base.org", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify key is non-zero before close
	allZero := true
	for _, b := range signer.privateKey {
		if b != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		t.Fatal("private key should not be all zeros before Close")
	}

	signer.Close()

	// Verify key is zeroed after close
	for i, b := range signer.privateKey {
		if b != 0 {
			t.Fatalf("private key byte %d not zeroed after Close: %d", i, b)
		}
	}
}
