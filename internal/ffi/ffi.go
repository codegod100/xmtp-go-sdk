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

// FfiIdentifierKind represents the type of identifier
// Note: UniFFI enum variants are 1-indexed (not 0-indexed)
type FfiIdentifierKind int32

const (
	IdentifierKindEthereum FfiIdentifierKind = 1 // First variant
	IdentifierKindPasskey  FfiIdentifierKind = 2 // Second variant
)

// serializeFfiIdentifier serializes an FfiIdentifier for UniFFI
// FfiIdentifier is a Record with fields:
// - identifier: String (serialized with BE length prefix + UTF8 bytes)
// - identifier_kind: Enum (serialized as i32 variant index)
func serializeFfiIdentifier(identifier string, kind FfiIdentifierKind) (C.RustBuffer, error) {
	if len(identifier) == 0 {
		return C.RustBuffer{}, nil
	}

	// Calculate total size:
	// - String: 4 (BE length) + len(identifier)
	// - Enum: 4 (i32 variant index)
	totalSize := 4 + len(identifier) + 4

	var status C.RustCallStatus
	buf := C.ffi_xmtpv3_rustbuffer_alloc(C.uint64_t(totalSize), &status)
	if status.code != CallSuccess {
		return C.RustBuffer{}, cStatusToError(status)
	}

	dst := unsafe.Slice((*byte)(unsafe.Pointer(buf.data)), totalSize)
	offset := 0

	// Write identifier string (inside Record: uses write() which adds BE length prefix)
	identifierLen := uint32(len(identifier))
	dst[offset] = byte(identifierLen >> 24)
	dst[offset+1] = byte(identifierLen >> 16)
	dst[offset+2] = byte(identifierLen >> 8)
	dst[offset+3] = byte(identifierLen)
	offset += 4
	copy(dst[offset:], identifier)
	offset += len(identifier)

	// Write identifier_kind enum (variant index as i32 BE)
	// Note: UniFFI serializes enums by variant index
	kindInt := int32(kind)
	dst[offset] = byte(kindInt >> 24)
	dst[offset+1] = byte(kindInt >> 16)
	dst[offset+2] = byte(kindInt >> 8)
	dst[offset+3] = byte(kindInt)

	buf.len = C.uint64_t(totalSize)
	return buf, nil
}

// cStringToBuffer creates a RustBuffer from a Go string
// Note: UniFFI String inputs use try_lift which expects raw bytes (no length prefix)
func cStringToBuffer(s string) (C.RustBuffer, error) {
	if len(s) == 0 {
		return C.RustBuffer{}, nil
	}

	var status C.RustCallStatus

	// For strings, UniFFI expects just the raw UTF-8 bytes (no length prefix)
	buf := C.ffi_xmtpv3_rustbuffer_alloc(C.uint64_t(len(s)), &status)
	if status.code != CallSuccess {
		return C.RustBuffer{}, cStatusToError(status)
	}

	// Copy the string bytes directly
	dst := unsafe.Slice((*byte)(unsafe.Pointer(buf.data)), buf.capacity)
	copy(dst, s)
	buf.len = C.uint64_t(len(s))

	return buf, nil
}

// cBytesToBuffer creates a RustBuffer from Go bytes
// UniFFI uses big-endian for the length prefix for Vec<u8>
func cBytesToBuffer(b []byte) (C.RustBuffer, error) {
	if len(b) == 0 {
		return C.RustBuffer{}, nil
	}

	// UniFFI serializes Vec<u8> as: [i32 length (BIG-ENDIAN)] + bytes
	totalLen := 4 + len(b)
	
	var status C.RustCallStatus
	buf := C.ffi_xmtpv3_rustbuffer_alloc(C.uint64_t(totalLen), &status)
	if status.code != CallSuccess {
		return C.RustBuffer{}, cStatusToError(status)
	}

	// Write i32 length prefix in BIG-ENDIAN (bytes crate default)
	dst := unsafe.Slice((*byte)(unsafe.Pointer(buf.data)), buf.capacity)
	dst[0] = byte(len(b) >> 24) // high byte first (big-endian)
	dst[1] = byte(len(b) >> 16)
	dst[2] = byte(len(b) >> 8)
	dst[3] = byte(len(b))
	
	// Write bytes
	copy(dst[4:], b)
	buf.len = C.uint64_t(totalLen)

	return buf, nil
}

// cBytesToBufferFromForeign creates a RustBuffer using ForeignBytes
// This is an alternative path that uses rustbuffer_from_bytes
func cBytesToBufferFromForeign(b []byte) (C.RustBuffer, error) {
	if len(b) == 0 {
		return C.RustBuffer{}, nil
	}

	var status C.RustCallStatus
	
	// Allocate C memory for the data
	cdata := C.malloc(C.size_t(len(b)))
	if cdata == nil {
		return C.RustBuffer{}, fmt.Errorf("failed to allocate memory")
	}
	C.memcpy(cdata, unsafe.Pointer(&b[0]), C.size_t(len(b)))
	
	// Create ForeignBytes
	fb := C.make_foreign_bytes((*C.uint8_t)(cdata), C.int32_t(len(b)))
	
	// Convert to RustBuffer
	buf := C.ffi_xmtpv3_rustbuffer_from_bytes(fb, &status)
	
	// Free C memory (Rust has copied the data)
	C.free(cdata)
	
	if status.code != CallSuccess {
		return C.RustBuffer{}, cStatusToError(status)
	}
	
	return buf, nil
}

// -- Buffer conversion for testing --

func cBufferToBytes(buf C.RustBuffer) []byte {
	if buf.len == 0 || buf.data == nil {
		return nil
	}
	
	// Read the buffer data
	data := unsafe.Slice((*byte)(unsafe.Pointer(buf.data)), buf.len)
	
	// UniFFI serializes Vec<u8> returns as [i32 length (BIG-ENDIAN)] + bytes
	if buf.len < 4 {
		return nil
	}
	
	// Read length prefix (BIG-ENDIAN i32)
	length := int(data[0])<<24 | int(data[1])<<16 | int(data[2])<<8 | int(data[3])
	
	// Copy the actual bytes (skip the length prefix)
	result := make([]byte, length)
	copy(result, data[4:4+length])
	
	return result
}

func cFreeBuffer(buf C.RustBuffer) {
	var status C.RustCallStatus
	C.ffi_xmtpv3_rustbuffer_free(buf, &status)
}

func cStatusToError(status C.RustCallStatus) error {
	var msg string
	if status.errorBuf.len > 0 && status.errorBuf.data != nil {
		msg = cBufferToString(status.errorBuf)
	}
	cFreeBuffer(status.errorBuf)

	switch status.code {
	case CallSuccess:
		return nil
	case CallError:
		return fmt.Errorf("error: %s", msg)
	case CallPanic:
		return fmt.Errorf("panic: %s", msg)
	default:
		return fmt.Errorf("unknown status code %d: %s", status.code, msg)
	}
}

// cBufferToString extracts a string from a RustBuffer
// For error messages, UniFFI serializes String as [i32 length (BIG-ENDIAN)] + utf8 bytes
func cBufferToString(buf C.RustBuffer) string {
	if buf.len == 0 || buf.data == nil {
		return ""
	}
	
	// Read the buffer data - String is returned as raw UTF-8 bytes
	// (String uses try_lift which expects RustBuffer::from_vec directly)
	data := unsafe.Slice((*byte)(unsafe.Pointer(buf.data)), buf.len)
	return string(data)
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

func GetVersionInfo() (string, error) {
	if !libLoaded {
		return "", loadErr
	}

	var status C.RustCallStatus
	result := C.uniffi_xmtpv3_fn_func_get_version_info(&status)
	if status.code != CallSuccess {
		return "", cStatusToError(status)
	}
	defer cFreeBuffer(result)

	// String returns use try_lift which expects raw bytes (no length prefix)
	return cBufferToString(result), nil
}

// EthereumGeneratePublicKey generates a public key from a private key
func EthereumGeneratePublicKey(privateKey []byte) ([]byte, error) {
	if !libLoaded {
		return nil, loadErr
	}
	if len(privateKey) != 32 {
		return nil, fmt.Errorf("private key must be 32 bytes, got %d", len(privateKey))
	}

	// UniFFI serializes Vec<u8> as [i32 length][bytes]
	keyBuf, err := cBytesToBuffer(privateKey)
	if err != nil {
		return nil, err
	}
	// Note: keyBuf is consumed by the function, don't free it

	var status C.RustCallStatus
	result := C.uniffi_xmtpv3_fn_func_ethereum_generate_public_key(keyBuf, &status)
	if status.code != CallSuccess {
		return nil, cStatusToError(status)
	}
	defer cFreeBuffer(result)

	return cBufferToBytes(result), nil
}

// EthereumHashPersonal hashes a message using Ethereum's personal_sign prefix
func EthereumHashPersonal(message string) ([]byte, error) {
	if !libLoaded {
		return nil, loadErr
	}

	msgBuf, err := cBytesToBuffer([]byte(message))
	if err != nil {
		return nil, err
	}
	// Note: msgBuf is consumed by the function

	var status C.RustCallStatus
	result := C.uniffi_xmtpv3_fn_func_ethereum_hash_personal(msgBuf, &status)
	if status.code != CallSuccess {
		return nil, cStatusToError(status)
	}
	defer cFreeBuffer(result)

	return cBufferToBytes(result), nil
}

// EthereumSignRecoverable signs a message with a private key (recoverable signature)
// hashing: 0 = no hashing, 1 = ethereum personal hash prefix
func EthereumSignRecoverable(message, privateKey []byte, hashing int8) ([]byte, error) {
	if !libLoaded {
		return nil, loadErr
	}
	if len(privateKey) != 32 {
		return nil, fmt.Errorf("private key must be 32 bytes, got %d", len(privateKey))
	}

	msgBuf, err := cBytesToBuffer(message)
	if err != nil {
		return nil, err
	}
	// Note: msgBuf is consumed by the function

	keyBuf, err := cBytesToBuffer(privateKey)
	if err != nil {
		return nil, err
	}
	// Note: keyBuf is consumed by the function

	var status C.RustCallStatus
	result := C.uniffi_xmtpv3_fn_func_ethereum_sign_recoverable(msgBuf, keyBuf, C.int8_t(hashing), &status)
	if status.code != CallSuccess {
		return nil, cStatusToError(status)
	}
	defer cFreeBuffer(result)

	return cBufferToBytes(result), nil
}

// EthereumAddressFromPublicKey derives an Ethereum address from a public key
// Note: The public key buffer is consumed by this call
func EthereumAddressFromPublicKey(publicKey []byte) (string, error) {
	if !libLoaded {
		return "", loadErr
	}

	keyBuf, err := cBytesToBuffer(publicKey)
	if err != nil {
		return "", err
	}
	// Note: keyBuf is consumed by the function

	var status C.RustCallStatus
	result := C.uniffi_xmtpv3_fn_func_ethereum_address_from_pubkey(keyBuf, &status)
	if status.code != CallSuccess {
		return "", cStatusToError(status)
	}
	defer cFreeBuffer(result)

	// String returns use try_lift which expects raw bytes (no length prefix)
	return cBufferToString(result), nil
}

// GenerateInboxID generates an inbox ID from an account identifier and nonce
// Note: The account identifier buffer is consumed by this call
func GenerateInboxID(accountIdentifier string, nonce uint64) (string, error) {
	return GenerateInboxIDWithKind(accountIdentifier, IdentifierKindEthereum, nonce)
}

// GenerateInboxIDWithKind generates an inbox ID with a specific identifier kind
func GenerateInboxIDWithKind(accountIdentifier string, kind FfiIdentifierKind, nonce uint64) (string, error) {
	if !libLoaded {
		return "", loadErr
	}

	// FfiIdentifier is a Record, so we need to serialize it properly
	idBuf, err := serializeFfiIdentifier(accountIdentifier, kind)
	if err != nil {
		return "", err
	}
	// Note: idBuf is consumed by the function

	var status C.RustCallStatus
	result := C.uniffi_xmtpv3_fn_func_generate_inbox_id(idBuf, C.uint64_t(nonce), &status)
	if status.code != CallSuccess {
		return "", cStatusToError(status)
	}
	defer cFreeBuffer(result)

	// String returns use try_lift which expects raw bytes (no length prefix)
	return cBufferToString(result), nil
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