// Package ffi provides PureGo bindings to the XMTP libxmtpv3 library
package ffi

import (
	"errors"
	"fmt"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"
)

// Library path
var libPath = "./result/lib/libxmtpv3.so"

// Call status codes
const (
	CallSuccess = 0
	CallError   = 1
	CallPanic   = 2
)

// Future poll results
const (
	FutureReady     = 0
	FutureMaybeReady = 1
)

// rustBuffer matches UniFFI's C layout
type rustBuffer struct {
	capacity uint64
	len      uint64
	data     *byte
}

type rustCallStatus struct {
	code     int8
	errorBuf rustBuffer
}

var lib uintptr
var libLoaded bool
var loadErr error

// FFI function pointers
var (
	// RustBuffer utilities
	ffi_rustbuffer_alloc    func(size uint64, status *rustCallStatus) rustBuffer
	ffi_rustbuffer_from_bytes func(bytes rustBuffer, status *rustCallStatus) rustBuffer
	ffi_rustbuffer_free     func(buf rustBuffer, status *rustCallStatus)
	
	// Async: connect_to_backend
	ffi_connect_to_backend func(v3Host, gatewayHost, clientMode, appVersion, authCallback, authHandle rustBuffer) uint64
	
	// Async: create_client
	ffi_create_client func(api, syncApi uint64, db, inboxId, accountIdentifier rustBuffer, nonce uint64, legacyKey, deviceSyncMode, allowOffline, forkRecovery rustBuffer) uint64
	
	// Sync: client methods
	ffi_client_inbox_id      func(ptr uint64, status *rustCallStatus) rustBuffer
	ffi_client_conversations func(ptr uint64, status *rustCallStatus) uint64
	
	// Sync: conversation methods
	ffi_conversation_id       func(ptr uint64, status *rustCallStatus) rustBuffer
	ffi_conversation_send_text func(ptr uint64, text rustBuffer, status *rustCallStatus) uint64
	
	// Free handles
	ffi_free_client       func(handle uint64, status *rustCallStatus)
	ffi_free_conversation func(handle uint64, status *rustCallStatus)
	
	// Future operations
	ffi_future_poll_u64     func(handle uint64, callback uintptr, data uint64)
	ffi_future_complete_u64 func(handle uint64, status *rustCallStatus) uint64
	ffi_future_free         func(handle uint64)
)

func init() {
	var err error
	lib, err = purego.Dlopen(libPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		loadErr = fmt.Errorf("failed to load %s: %w", libPath, err)
		return
	}
	
	// Register functions
	purego.RegisterLibFunc(&ffi_rustbuffer_alloc, lib, "uniffi_uniffi_rustbuffer_alloc")
	purego.RegisterLibFunc(&ffi_rustbuffer_from_bytes, lib, "uniffi_uniffi_rustbuffer_from_bytes")
	purego.RegisterLibFunc(&ffi_rustbuffer_free, lib, "uniffi_uniffi_rustbuffer_free")
	
	purego.RegisterLibFunc(&ffi_connect_to_backend, lib, "uniffi_xmtpv3_fn_func_connect_to_backend")
	purego.RegisterLibFunc(&ffi_create_client, lib, "uniffi_xmtpv3_fn_func_create_client")
	
	purego.RegisterLibFunc(&ffi_client_inbox_id, lib, "uniffi_xmtpv3_fn_method_ffixmtpclient_inbox_id")
	purego.RegisterLibFunc(&ffi_client_conversations, lib, "uniffi_xmtpv3_fn_method_ffixmtpclient_conversations")
	
	purego.RegisterLibFunc(&ffi_conversation_id, lib, "uniffi_xmtpv3_fn_method_fficonversation_id")
	purego.RegisterLibFunc(&ffi_conversation_send_text, lib, "uniffi_xmtpv3_fn_method_fficonversation_send_text")
	
	purego.RegisterLibFunc(&ffi_free_client, lib, "uniffi_xmtpv3_fn_free_ffixmtpclient")
	purego.RegisterLibFunc(&ffi_free_conversation, lib, "uniffi_xmtpv3_fn_free_fficonversation")
	
	purego.RegisterLibFunc(&ffi_future_poll_u64, lib, "ffi_xmtpv3_rust_future_poll_u64")
	purego.RegisterLibFunc(&ffi_future_complete_u64, lib, "ffi_xmtpv3_rust_future_complete_u64")
	purego.RegisterLibFunc(&ffi_future_free, lib, "ffi_xmtpv3_rust_future_free_u64")
	
	libLoaded = true
}

// SetLibraryPath changes the path to libxmtpv3.so
func SetLibraryPath(path string) {
	libPath = path
}

// IsLoaded returns true if the library was loaded successfully
func IsLoaded() bool {
	return libLoaded
}

// LoadError returns the error from loading the library
func LoadError() error {
	return loadErr
}

// Buffer helpers

func stringToBuffer(s string) (rustBuffer, error) {
	return bytesToBuffer([]byte(s))
}

func bytesToBuffer(b []byte) (rustBuffer, error) {
	if len(b) == 0 {
		return rustBuffer{}, nil
	}
	
	var status rustCallStatus
	buf := ffi_rustbuffer_alloc(uint64(len(b)), &status)
	if status.code != CallSuccess {
		return rustBuffer{}, statusToError(status)
	}
	
	// Copy data
	dst := unsafe.Slice(buf.data, buf.capacity)
	copy(dst, b)
	buf.len = uint64(len(b))
	
	return buf, nil
}

func bufferToBytes(buf rustBuffer) []byte {
	if buf.len == 0 || buf.data == nil {
		return nil
	}
	src := unsafe.Slice(buf.data, buf.len)
	result := make([]byte, buf.len)
	copy(result, src)
	return result
}

func bufferToString(buf rustBuffer) string {
	return string(bufferToBytes(buf))
}

func freeBuffer(buf rustBuffer) {
	if !libLoaded {
		return
	}
	var status rustCallStatus
	ffi_rustbuffer_free(buf, &status)
}

func statusToError(status rustCallStatus) error {
	msg := bufferToString(status.errorBuf)
	freeBuffer(status.errorBuf)
	
	switch status.code {
	case CallError:
		return errors.New(msg)
	case CallPanic:
		return fmt.Errorf("panic: %s", msg)
	default:
		return fmt.Errorf("unknown error: %s", msg)
	}
}

// Async helpers - block on futures

var callbackMu sync.Mutex
var callbackCond = sync.NewCond(&callbackMu)
var callbackResult int8

//go:export futureCallback
func futureCallback(data uint64, result int8) {
	callbackMu.Lock()
	callbackResult = result
	callbackMu.Unlock()
	callbackCond.Broadcast()
}

func awaitFuture(handle uint64) error {
	if !libLoaded {
		return loadErr
	}
	
	callbackMu.Lock()
	defer callbackMu.Unlock()
	
	// Poll with callback
	cb := purego.NewCallback(futureCallback)
	for {
		ffi_future_poll_u64(handle, cb, 0)
		callbackCond.Wait()
		
		if callbackResult == FutureReady {
			break
		}
	}
	
	return nil
}

func awaitFutureU64(handle uint64) (uint64, error) {
	if err := awaitFuture(handle); err != nil {
		return 0, err
	}
	
	var status rustCallStatus
	result := ffi_future_complete_u64(handle, &status)
	if status.code != CallSuccess {
		return 0, statusToError(status)
	}
	
	return result, nil
}
