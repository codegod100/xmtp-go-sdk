package ffi

import (
	"sync"
	"testing"
)

// TestGoString tests the GoString conversion function
func TestGoString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty", "", ""},
		{"simple", "hello", "hello"},
		{"with spaces", "hello world", "hello world"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a C string
			ptr := CString(tt.input)
			if ptr == 0 && tt.input != "" {
				t.Fatalf("CString returned 0 for non-empty input")
			}

			// Convert back to Go string
			result := GoString(ptr)
			if result != tt.expected {
				t.Errorf("GoString: got %q, want %q", result, tt.expected)
			}
		})
	}

	// Test null pointer
	t.Run("null pointer", func(t *testing.T) {
		result := GoString(0)
		if result != "" {
			t.Errorf("GoString(0): got %q, want empty string", result)
		}
	})
}

// TestGoBytes tests the GoBytes conversion function
func TestGoBytes(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		result := GoBytes(0, 0)
		if result != nil {
			t.Errorf("GoBytes(0, 0): got %v, want nil", result)
		}
	})

	t.Run("zero length", func(t *testing.T) {
		ptr := CString("test")
		result := GoBytes(ptr, 0)
		if result != nil {
			t.Errorf("GoBytes with len 0: got %v, want nil", result)
		}
	})
}

// TestCFree tests that CFree is a no-op
func TestCFree(t *testing.T) {
	// Just verify it doesn't panic
	CFree(0)
	CFree(12345)
}

// TestLibraryPaths tests that library paths are defined
func TestLibraryPaths(t *testing.T) {
	if len(libraryPaths) == 0 {
		t.Error("libraryPaths should not be empty")
	}

	// Check for expected paths
	expectedPaths := []string{
		"libxmtp_ffi.so",
		"libxmtp_ffi.dylib",
		"xmtp_ffi.dll",
	}

	for _, expected := range expectedPaths {
		found := false
		for _, path := range libraryPaths {
			if path == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected library path %q not found", expected)
		}
	}
}

// TestErrLibraryNotLoaded tests the error value
func TestErrLibraryNotLoaded(t *testing.T) {
	if ErrLibraryNotLoaded == nil {
		t.Error("ErrLibraryNotLoaded should not be nil")
	}
	if ErrLibraryNotLoaded.Error() == "" {
		t.Error("ErrLibraryNotLoaded should have a message")
	}
}

// TestGetLibraryWithoutLoad tests GetLibrary behavior when not loaded
func TestGetLibraryWithoutLoad(t *testing.T) {
	// Reset the library state for this test
	originalLibrary := library
	originalOnce := libraryOnce
	defer func() {
		library = originalLibrary
		libraryOnce = originalOnce
	}()

	library = nil
	libraryOnce = sync.Once{}

	lib, err := GetLibrary()
	if err != ErrLibraryNotLoaded {
		t.Errorf("GetLibrary error: got %v, want %v", err, ErrLibraryNotLoaded)
	}
	if lib != nil {
		t.Error("GetLibrary should return nil when not loaded")
	}
}

// TestTypes tests FFI type definitions
func TestTypes(t *testing.T) {
	t.Run("XmtpEnv", func(t *testing.T) {
		envs := []XmtpEnv{
			XmtpEnvLocal,
			XmtpEnvDev,
			XmtpEnvProduction,
			XmtpEnvTestnetStaging,
			XmtpEnvTestnetDev,
			XmtpEnvTestnet,
			XmtpEnvMainnet,
		}
		for i, env := range envs {
			if int(env) != i {
				t.Errorf("XmtpEnv %d: got %v", i, env)
			}
		}
	})

	t.Run("XmtpConversationType", func(t *testing.T) {
		types := []XmtpConversationType{
			XmtpConversationTypeDm,
			XmtpConversationTypeGroup,
		}
		for i, ct := range types {
			if int(ct) != i {
				t.Errorf("XmtpConversationType %d: got %v", i, ct)
			}
		}
	})

	t.Run("XmtpConsentState", func(t *testing.T) {
		states := []XmtpConsentState{
			XmtpConsentStateUnknown,
			XmtpConsentStateAllowed,
			XmtpConsentStateDenied,
		}
		for i, state := range states {
			if int(state) != i {
				t.Errorf("XmtpConsentState %d: got %v", i, state)
			}
		}
	})

	t.Run("XmtpContentType", func(t *testing.T) {
		types := []XmtpContentType{
			XmtpContentTypeText,
			XmtpContentTypeMarkdown,
			XmtpContentTypeReply,
			XmtpContentTypeReaction,
			XmtpContentTypeAttachment,
		}
		for i, ct := range types {
			if int(ct) != i {
				t.Errorf("XmtpContentType %d: got %v", i, ct)
			}
		}
	})

	t.Run("XmtpDeliveryStatus", func(t *testing.T) {
		statuses := []XmtpDeliveryStatus{
			XmtpDeliveryStatusUnpublished,
			XmtpDeliveryStatusPublished,
			XmtpDeliveryStatusFailed,
		}
		for i, status := range statuses {
			if int(status) != i {
				t.Errorf("XmtpDeliveryStatus %d: got %v", i, status)
			}
		}
	})
}

// TestClientOptionsStruct tests the XmtpClientOptions struct
func TestClientOptionsStruct(t *testing.T) {
	opts := XmtpClientOptions{
		Env:                 XmtpEnvProduction,
		DbPath:              CString("/tmp/test.db"),
		DbEncryptionKeyLen:  32,
		AppVersion:          CString("test/1.0.0"),
		DisableAutoRegister: true,
		StructuredLogging:   true,
		LogLevel:            4,
	}

	if opts.Env != XmtpEnvProduction {
		t.Errorf("Env: got %v", opts.Env)
	}
	if opts.DisableAutoRegister != true {
		t.Error("DisableAutoRegister should be true")
	}
	if opts.StructuredLogging != true {
		t.Error("StructuredLogging should be true")
	}
	if opts.LogLevel != 4 {
		t.Errorf("LogLevel: got %v", opts.LogLevel)
	}
}

// TestResultTypes tests FFI result types
func TestResultTypes(t *testing.T) {
	t.Run("XmtpStringResult", func(t *testing.T) {
		r := XmtpStringResult{
			Value: CString("test"),
			Error: 0,
		}
		if GoString(r.Value) != "test" {
			t.Errorf("Value: got %v", GoString(r.Value))
		}
	})

	t.Run("XmtpIntResult", func(t *testing.T) {
		r := XmtpIntResult{
			Value: 42,
			Error: 0,
		}
		if r.Value != 42 {
			t.Errorf("Value: got %v", r.Value)
		}
	})

	t.Run("XmtpBoolResult", func(t *testing.T) {
		r := XmtpBoolResult{
			Value: true,
			Error: 0,
		}
		if r.Value != true {
			t.Error("Value should be true")
		}
	})

	t.Run("XmtpBytesResult", func(t *testing.T) {
		r := XmtpBytesResult{
			Data:  0,
			Len:   3,
			Error: 0,
		}
		if r.Len != 3 {
			t.Errorf("Len: got %v", r.Len)
		}
	})
}
