package ffi

import (
	"unsafe"
)

// CString creates a C-style null-terminated string from a Go string
func CString(s string) uintptr {
	if s == "" {
		return 0
	}
	
	// Append null terminator
	bs := append([]byte(s), 0)
	return uintptr(unsafe.Pointer(&bs[0]))
}

// GoString converts a C-style string to a Go string
func GoString(s uintptr) string {
	if s == 0 {
		return ""
	}
	
	// Find null terminator
	var i int
	for {
		b := *(*byte)(unsafe.Pointer(uintptr(s) + uintptr(i)))
		if b == 0 {
			break
		}
		i++
		if i > 1000000 { // Safety limit
			break
		}
	}
	
	if i == 0 {
		return ""
	}
	
	// Create Go string from the bytes
	bytes := make([]byte, i)
	for j := 0; j < i; j++ {
		bytes[j] = *(*byte)(unsafe.Pointer(uintptr(s) + uintptr(j)))
	}
	return string(bytes)
}

// GoBytes converts a C byte array to a Go byte slice
func GoBytes(data uintptr, len int) []byte {
	if data == 0 || len == 0 {
		return nil
	}
	
	result := make([]byte, len)
	for i := 0; i < len; i++ {
		result[i] = *(*byte)(unsafe.Pointer(uintptr(data) + uintptr(i)))
	}
	return result
}

// CBytes creates a C byte array pointer from a Go byte slice
func CBytes(b []byte) uintptr {
	if len(b) == 0 {
		return 0
	}
	
	return uintptr(unsafe.Pointer(&b[0]))
}

// FfiError represents an error from the FFI layer
type FfiError struct {
	Code    int
	Message string
}

// NewFfiError creates a new FFI error
func NewFfiError(code int, message string) *FfiError {
	return &FfiError{Code: code, Message: message}
}

// Error implements the error interface
func (e *FfiError) Error() string {
	return e.Message
}
