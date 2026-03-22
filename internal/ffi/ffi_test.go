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

	// This will fail without a real backend, but we can test the call path
	_, err := ConnectToBackend("invalid-host:1234", "", "test-sdk/0.1.0")
	if err != nil {
		t.Logf("Expected error connecting to invalid backend: %v", err)
		// This is expected - we're just testing the call path works
	} else {
		t.Log("Connected to backend (unexpected)")
	}
}