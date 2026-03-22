package ffi

import (
	"fmt"
	"testing"
)

func TestLibraryLoad(t *testing.T) {
	if !IsLoaded() {
		t.Fatalf("Failed to load library: %v", LoadError())
	}
	t.Log("Library loaded successfully")
}

func TestBufferRoundTrip(t *testing.T) {
	if !IsLoaded() {
		t.Skip("Library not loaded")
	}

	original := []byte("hello world")
	result, err := BufferRoundTrip(original)
	if err != nil {
		t.Fatalf("BufferRoundTrip failed: %v", err)
	}

	if string(result) != string(original) {
		t.Errorf("Expected %q, got %q", original, result)
	} else {
		t.Log("Buffer round-trip successful")
	}
}

func TestEmptyBuffer(t *testing.T) {
	if !IsLoaded() {
		t.Skip("Library not loaded")
	}

	original := []byte{}
	result, err := BufferRoundTrip(original)
	if err != nil {
		t.Fatalf("BufferRoundTrip failed: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("Expected empty buffer, got %v", result)
	} else {
		t.Log("Empty buffer round-trip successful")
	}
}

func TestConnectToBackend(t *testing.T) {
	if !IsLoaded() {
		t.Skip("Library not loaded")
	}

	// Skip this test - it hangs waiting for a real backend
	t.Skip("Skipping backend connection test - requires live XMTP network")
}

func TestEthereumUtilities(t *testing.T) {
	if !IsLoaded() {
		t.Skip("Library not loaded")
	}

	// Valid test private key (DO NOT USE IN PRODUCTION)
	// This is from the Rust test: 90b7388a7427358cb7fc7e9042805b1942eae47ee783e627a989719da35e76fb
	privateKey, err := hexDecode("90b7388a7427358cb7fc7e9042805b1942eae47ee783e627a989719da35e76fb")
	if err != nil {
		t.Fatalf("Failed to decode private key: %v", err)
	}

	// Test public key generation
	pubKey, err := EthereumGeneratePublicKey(privateKey)
	if err != nil {
		t.Fatalf("EthereumGeneratePublicKey failed: %v", err)
	}
	if len(pubKey) != 65 {
		t.Errorf("Public key should be 65 bytes, got %d", len(pubKey))
	}
	if pubKey[0] != 0x04 {
		t.Errorf("Public key should start with 0x04, got 0x%02x", pubKey[0])
	}
	t.Logf("Generated public key: %x (len=%d)", pubKey[:min(10, len(pubKey))], len(pubKey))

	// Test address derivation
	address, err := EthereumAddressFromPublicKey(pubKey)
	if err != nil {
		t.Fatalf("EthereumAddressFromPublicKey failed: %v", err)
	}
	if len(address) != 42 || address[:2] != "0x" {
		t.Errorf("Invalid Ethereum address: %s", address)
	}
	t.Logf("Derived address: %s", address)

	// Test message hashing
	message := "Hello XMTP!"
	hash, err := EthereumHashPersonal(message)
	if err != nil {
		t.Fatalf("EthereumHashPersonal failed: %v", err)
	}
	if len(hash) != 32 {
		t.Errorf("Hash should be 32 bytes, got %d", len(hash))
	}
	t.Logf("Message hash: %x", hash)

	// Test signing
	signature, err := EthereumSignRecoverable([]byte(message), privateKey, 1) // 1 = hash with Ethereum prefix
	if err != nil {
		t.Fatalf("EthereumSignRecoverable failed: %v", err)
	}
	if len(signature) != 65 {
		t.Errorf("Signature should be 65 bytes, got %d", len(signature))
	}
	if signature[64] != 27 && signature[64] != 28 {
		t.Errorf("Signature v value should be 27 or 28, got %d", signature[64])
	}
	t.Logf("Signature: %x (len=%d)", signature[:min(10, len(signature))], len(signature))
}

func hexDecode(s string) ([]byte, error) {
	result := make([]byte, len(s)/2)
	for i := 0; i < len(s); i += 2 {
		var b byte
		_, err := fmt.Sscanf(s[i:i+2], "%02x", &b)
		if err != nil {
			return nil, err
		}
		result[i/2] = b
	}
	return result, nil
}

func TestGenerateInboxID(t *testing.T) {
	if !IsLoaded() {
		t.Skip("Library not loaded")
	}

	address := "0x1234567890123456789012345678901234567890"
	nonce := uint64(0)

	inboxID, err := GenerateInboxID(address, nonce)
	if err != nil {
		t.Fatalf("GenerateInboxID failed: %v", err)
	}
	if len(inboxID) == 0 {
		t.Fatal("Inbox ID should not be empty")
	}
	t.Logf("Generated inbox ID: %s", inboxID)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestGetVersionInfo(t *testing.T) {
	if !IsLoaded() {
		t.Skip("Library not loaded")
	}

	version, err := GetVersionInfo()
	if err != nil {
		t.Fatalf("GetVersionInfo failed: %v", err)
	}
	if len(version) == 0 {
		t.Fatal("Version should not be empty")
	}
	t.Logf("Version: %s", version)
}