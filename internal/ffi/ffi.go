// Package ffi provides bindings to the XMTP libxmtpv3 library.
// We use CGO for buffer operations (struct returns not supported by PureGo on Linux)
// and PureGo for all other FFI calls.
package ffi

/*
#cgo CFLAGS: -I/home/nandi/code/chat/xmtp-go-sdk/result/lib
#cgo LDFLAGS: -L/home/nandi/code/chat/xmtp-go-sdk/result/lib -lxmtpv3 -ldl -lm -lpthread

#include "xmtpv3_shim.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"os"
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
	FutureReady        = 0
	FutureMaybeReady   = 1
	FuturePollComplete = 0
)

var lib uintptr
var libLoaded bool
var loadErr error

// FFI function pointers (PureGo for non-struct returns)
var (
	ffi_free_client         func(handle uint64, status *C.RustCallStatus)
	ffi_free_conversation   func(handle uint64, status *C.RustCallStatus)
	ffi_future_poll_u64     func(handle uint64, callback uintptr, data uint64)
	ffi_future_complete_u64 func(handle uint64, status *C.RustCallStatus) uint64
	ffi_future_free         func(handle uint64)
)

func init() {
	// Check for env var override first
	libPath = "./result/lib/libxmtpv3.so" // default
	if p := os.Getenv("XMTP_LIB_PATH"); p != "" {
		libPath = p
	}

	// Debug: print what we're loading
	fmt.Fprintf(os.Stderr, "DEBUG: Loading library from: %s\n", libPath)

	var err error
	lib, err = purego.Dlopen(libPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		loadErr = fmt.Errorf("failed to load %s: %w", libPath, err)
		return
	}

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

// -- Buffer operations (CGO-only, not exposed outside package) --

func cBytesToBuffer(b []byte) (C.RustBuffer, error) {
	if len(b) == 0 {
		return C.RustBuffer{}, nil
	}

	var status C.RustCallStatus
	cdata := C.CBytes(b)
	defer C.free(cdata)

	fb := C.ForeignBytes{
		len:  C.int32_t(len(b)),
		data: (*C.uint8_t)(cdata),
	}
	buf := C.ffi_xmtpv3_rustbuffer_from_bytes(fb, &status)
	if status.code != CallSuccess {
		return C.RustBuffer{}, cStatusToError(status)
	}
	return buf, nil
}

func cBufferToBytes(buf C.RustBuffer) []byte {
	if buf.len == 0 || buf.data == nil {
		return nil
	}
	src := unsafe.Slice((*byte)(unsafe.Pointer(buf.data)), buf.len)
	result := make([]byte, buf.len)
	copy(result, src)
	return result
}

func cFreeBuffer(buf C.RustBuffer) {
	var status C.RustCallStatus
	C.ffi_xmtpv3_rustbuffer_free(buf, &status)
}

func cStatusToError(status C.RustCallStatus) error {
	msg := string(cBufferToBytes(status.errorBuf))
	cFreeBuffer(status.errorBuf)

	switch status.code {
	case CallError:
		return errors.New(msg)
	case CallPanic:
		return fmt.Errorf("panic: %s", msg)
	default:
		return fmt.Errorf("unknown error: %s", msg)
	}
}

// -- Exported functions for testing --

// BufferRoundTrip tests the buffer conversion functions
func BufferRoundTrip(input []byte) ([]byte, error) {
	buf, err := cBytesToBuffer(input)
	if err != nil {
		return nil, err
	}
	defer cFreeBuffer(buf)
	return cBufferToBytes(buf), nil
}

// -- Async helpers --

var callbackMu sync.Mutex
var callbackCond = sync.NewCond(&callbackMu)
var callbackResult int8

//export futureCallback
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

	var status C.RustCallStatus
	result := ffi_future_complete_u64(handle, &status)
	if status.code != CallSuccess {
		return 0, cStatusToError(status)
	}

	return result, nil
}

func freeFuture(handle uint64) {
	ffi_future_free(handle)
}

// -- Client functions --

// ConnectToBackend connects to the XMTP backend
func ConnectToBackend(v3Host, gatewayHost, appVersion string) (apiHandle uint64, err error) {
	if !libLoaded {
		return 0, loadErr
	}

	v3HostBuf, err := cBytesToBuffer([]byte(v3Host))
	if err != nil {
		return 0, err
	}
	defer cFreeBuffer(v3HostBuf)

	gatewayHostBuf, err := cBytesToBuffer([]byte(gatewayHost))
	if err != nil {
		return 0, err
	}
	defer cFreeBuffer(gatewayHostBuf)

	clientModeBuf, err := cBytesToBuffer([]byte("ReadWrite"))
	if err != nil {
		return 0, err
	}
	defer cFreeBuffer(clientModeBuf)

	appVersionBuf, err := cBytesToBuffer([]byte(appVersion))
	if err != nil {
		return 0, err
	}
	defer cFreeBuffer(appVersionBuf)

	var status C.RustCallStatus
	futureHandle := C.uniffi_xmtpv3_fn_func_connect_to_backend(
		v3HostBuf, gatewayHostBuf, clientModeBuf, appVersionBuf,
		C.RustBuffer{}, C.RustBuffer{},
		&status,
	)
	if status.code != CallSuccess {
		return 0, cStatusToError(status)
	}

	apiHandle, err = awaitFutureU64(uint64(futureHandle))
	freeFuture(uint64(futureHandle))
	return apiHandle, err
}

// FreeClient frees a client handle
func FreeClient(handle uint64) error {
	if !libLoaded {
		return loadErr
	}
	var status C.RustCallStatus
	ffi_free_client(handle, &status)
	if status.code != CallSuccess {
		return cStatusToError(status)
	}
	return nil
}

// FreeConversation frees a conversation handle
func FreeConversation(handle uint64) error {
	if !libLoaded {
		return loadErr
	}
	var status C.RustCallStatus
	ffi_free_conversation(handle, &status)
	if status.code != CallSuccess {
		return cStatusToError(status)
	}
	return nil
}