//! C FFI bindings for XMTP MLS library
//!
//! This crate provides C-compatible FFI bindings for the XMTP MLS library,
//! designed to be used with PureGo for Go applications.

#![allow(clippy::missing_safety_doc)]
#![allow(clippy::not_unsafe_ptr_arg_deref)]

mod client;
mod conversations;
mod conversation;
mod message;
mod types;
mod error;
mod signer;

pub use types::*;
pub use error::*;

use std::ffi::{c_char, c_int, c_void, CStr, CString};
use std::ptr;

/// Initialize the FFI library
/// Must be called before any other functions
#[no_mangle]
pub extern "C" fn xmtp_init() -> c_int {
    // Initialize tokio runtime (lazy static)
    let _ = get_runtime();
    0
}

/// Get the library version
/// Returns a heap-allocated string that must be freed with xmtp_string_free
#[no_mangle]
pub extern "C" fn xmtp_version() -> *mut c_char {
    let version = env!("CARGO_PKG_VERSION").to_string();
    match CString::new(version) {
        Ok(s) => s.into_raw(),
        Err(_) => ptr::null_mut(),
    }
}

/// Free a string allocated by the library
#[no_mangle]
pub extern "C" fn xmtp_string_free(s: *mut c_char) {
    if !s.is_null() {
        unsafe {
            drop(CString::from_raw(s));
        }
    }
}

/// Free a byte array allocated by the library
#[no_mangle]
pub extern "C" fn xmtp_bytes_free(data: *mut u8, len: usize) {
    if !data.is_null() && len > 0 {
        unsafe {
            drop(Vec::from_raw_parts(data, len, len));
        }
    }
}

// Lazy static tokio runtime
fn get_runtime() -> &'static tokio::runtime::Runtime {
    use std::sync::OnceLock;
    static RUNTIME: OnceLock<tokio::runtime::Runtime> = OnceLock::new();
    RUNTIME.get_or_init(|| {
        tokio::runtime::Builder::new_multi_thread()
            .enable_all()
            .build()
            .expect("Failed to create tokio runtime")
    })
}

// Helper to convert Result to XmtpResult
fn result_to_ffi<T, F>(result: Result<T, XmtpFfiError>, out: *mut T, map_err: F) -> XmtpResult
where
    F: Fn(XmtpFfiError) -> XmtpFfiError,
{
    match result {
        Ok(value) => {
            if !out.is_null() {
                unsafe { *out = value };
            }
            XmtpResult { error: ptr::null_mut() }
        }
        Err(e) => XmtpResult { error: Box::into_raw(Box::new(map_err(e))) },
    }
}
