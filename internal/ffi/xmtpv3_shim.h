// xmtpv3_shim.h - Minimal shim for CGO compatibility
// Removes _Nonnull/_Nullable annotations that CGO can't parse

#ifndef XMTPV3_SHIM_H
#define XMTPV3_SHIM_H

#include <stdint.h>
#include <stdbool.h>
#include <stdlib.h>
#include <string.h>

// Define away the annotations
#define _Nonnull
#define _Nullable
#define _Null_unspecified

// RustBuffer matches UniFFI's C layout
typedef struct RustBuffer {
    uint64_t capacity;
    uint64_t len;
    uint8_t *data;
} RustBuffer;

typedef struct RustCallStatus {
    int8_t code;
    RustBuffer errorBuf;
} RustCallStatus;

typedef struct ForeignBytes {
    int32_t len;
    const uint8_t *data;
} ForeignBytes;

// Call status codes
#define UNIFFI_CALL_SUCCESS 0
#define UNIFFI_CALL_ERROR 1
#define UNIFFI_CALL_PANIC 2

// Helper to create ForeignBytes from raw bytes (properly aligned)
static inline ForeignBytes make_foreign_bytes(const uint8_t *data, int32_t len) {
    ForeignBytes fb;
    fb.len = len;
    fb.data = data;
    return fb;
}

// Buffer functions - use ffi_xmtpv3_ prefix
extern RustBuffer ffi_xmtpv3_rustbuffer_alloc(uint64_t size, RustCallStatus *out_status);
extern RustBuffer ffi_xmtpv3_rustbuffer_from_bytes(ForeignBytes bytes, RustCallStatus *out_status);
extern void ffi_xmtpv3_rustbuffer_free(RustBuffer buf, RustCallStatus *out_status);
extern RustBuffer ffi_xmtpv3_rustbuffer_reserve(RustBuffer buf, uint64_t additional, RustCallStatus *out_status);

// Ethereum utilities
extern RustBuffer uniffi_xmtpv3_fn_func_get_version_info(RustCallStatus *out_status);
extern RustBuffer uniffi_xmtpv3_fn_func_ethereum_generate_public_key(RustBuffer private_key, RustCallStatus *out_status);
extern RustBuffer uniffi_xmtpv3_fn_func_ethereum_hash_personal(RustBuffer message, RustCallStatus *out_status);
extern RustBuffer uniffi_xmtpv3_fn_func_ethereum_sign_recoverable(RustBuffer msg, RustBuffer private_key, int8_t hashing, RustCallStatus *out_status);
extern RustBuffer uniffi_xmtpv3_fn_func_ethereum_address_from_pubkey(RustBuffer pubkey, RustCallStatus *out_status);
extern RustBuffer uniffi_xmtpv3_fn_func_generate_inbox_id(RustBuffer account_identifier, uint64_t nonce, RustCallStatus *out_status);

// Connect to backend
extern uint64_t uniffi_xmtpv3_fn_func_connect_to_backend(
    RustBuffer v3_host, RustBuffer gateway_host, RustBuffer client_mode,
    RustBuffer app_version, RustBuffer auth_callback, RustBuffer auth_handle,
    RustCallStatus *out_status);

// Client methods
extern RustBuffer uniffi_xmtpv3_fn_method_ffixmtpclient_inbox_id(uint64_t ptr, RustCallStatus *out_status);
extern uint64_t uniffi_xmtpv3_fn_method_ffixmtpclient_conversations(uint64_t ptr, RustCallStatus *out_status);
extern void uniffi_xmtpv3_fn_free_ffixmtpclient(uint64_t handle, RustCallStatus *out_status);

// Conversation methods
extern RustBuffer uniffi_xmtpv3_fn_method_fficonversation_id(uint64_t ptr, RustCallStatus *out_status);
extern uint64_t uniffi_xmtpv3_fn_method_fficonversation_send_text(uint64_t ptr, RustBuffer text, RustCallStatus *out_status);
extern void uniffi_xmtpv3_fn_free_fficonversation(uint64_t handle, RustCallStatus *out_status);

// Future operations
extern void ffi_xmtpv3_rust_future_poll_u64(uint64_t handle, void *callback, uint64_t callback_data);
extern uint64_t ffi_xmtpv3_rust_future_complete_u64(uint64_t handle, RustCallStatus *out_status);
extern void ffi_xmtpv3_rust_future_free_u64(uint64_t handle);

#endif // XMTPV3_SHIM_H
