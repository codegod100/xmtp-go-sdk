package ffi

import (
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