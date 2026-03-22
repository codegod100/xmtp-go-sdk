package ffi

import (
	"unsafe"
)

// CString creates a C string from a Go string and returns a uintptr
func CString(s string) uintptr {
	if s == "" {
		return 0
	}

	// Allocate memory for string + null terminator using Go's allocator
	bs := append([]byte(s), 0)
	return uintptr(unsafe.Pointer(&bs[0]))
}

// GoString converts a C string (uintptr) to a Go string
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

// CFree is a no-op in PureGo (memory is managed by Go)
func CFree(ptr uintptr) {
	// No-op: Go's garbage collector will clean up
}
