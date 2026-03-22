package ffi

import (
	"unsafe"
)

// CString creates a C string from a Go string
func CString(s string) uintptr {
	if s == "" {
		return 0
	}
	
	// Append null terminator
	bs := append([]byte(s), 0)
	return uintptr(unsafe.Pointer(&bs[0]))
}

// CStringOwned creates a C string that must be freed
func CStringOwned(s string) uintptr {
	if s == "" {
		return 0
	}
	
	bs := append([]byte(s), 0)
	// Keep a reference to prevent GC
	ptr := uintptr(unsafe.Pointer(&bs[0]))
	_ = ptr
	return ptr
}

// GoString converts a C string to a Go string
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

// CBytes creates a C byte array from a Go byte slice
func CBytes(b []byte) uintptr {
	if len(b) == 0 {
		return 0
	}
	
	return uintptr(unsafe.Pointer(&b[0]))
}

// FreeString frees a string allocated by the library
func FreeString(s uintptr) {
	lib, err := GetLibrary()
	if err != nil {
		return
	}
	lib.XmtpStringFree(s)
}

// FreeBytes frees bytes allocated by the library
func FreeBytes(data uintptr, len uintptr) {
	lib, err := GetLibrary()
	if err != nil {
		return
	}
	lib.XmtpBytesFree(data, len)
}

// IsOk checks if an XmtpResult is successful
func IsOk(result XmtpResult) bool {
	return result.Error == 0
}

// GetError gets the error message from an XmtpResult
func GetError(result XmtpResult) string {
	if result.Error == 0 {
		return ""
	}
	
	err := (*XmtpFfiError)(unsafe.Pointer(result.Error))
	return GoString(err.Message)
}

// StringResultToGo converts an XmtpStringResult to a Go string and frees the C string
func StringResultToGo(result XmtpStringResult) (string, error) {
	if result.Error != 0 {
		err := (*XmtpFfiError)(unsafe.Pointer(result.Error))
		msg := GoString(err.Message)
		return "", NewFfiError(int(err.Code), msg)
	}
	
	if result.Value == 0 {
		return "", nil
	}
	
	s := GoString(result.Value)
	FreeString(result.Value)
	return s, nil
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
