//! Client FFI functions

use std::ffi::{c_char, c_void, CStr, CString};
use std::ptr;
use std::sync::Arc;
use parking_lot::RwLock;

use crate::types::*;
use crate::error::XmtpFfiError;
use crate::signer::FfiSigner;

/// Internal client wrapper
pub struct XmtpClientInner {
    // TODO: Add actual client when integrated with libxmtp
    // client: Arc<RustXmtpClient>,
    inbox_id: String,
    signer: Option<Arc<FfiSigner>>,
    options: XmtpClientOptions,
}

impl XmtpClientInner {
    pub fn new(
        signer: Option<Arc<FfiSigner>>,
        identifier: Option<XmtpIdentifier>,
        options: XmtpClientOptions,
    ) -> Result<Self, XmtpFfiError> {
        // TODO: Create actual client
        // For now, return a mock client
        Ok(XmtpClientInner {
            inbox_id: String::new(),
            signer,
            options,
        })
    }
}

/// Create a new XMTP client with a signer
/// 
/// # Safety
/// The signer_callback must remain valid for the lifetime of the client.
#[no_mangle]
pub extern "C" fn xmtp_client_create(
    signer_callback: XmtpSignerCallback,
    signer_user_data: *mut c_void,
    identifier: XmtpIdentifier,
    opts: *const XmtpClientOptions,
    out_client: *mut XmtpClientHandle,
) -> XmtpResult {
    let result = || {
        let opts = if opts.is_null() {
            XmtpClientOptions::default()
        } else {
            unsafe { *opts }
        };
        
        let signer = Arc::new(FfiSigner::new(signer_callback, signer_user_data));
        let identifier = identifier;
        
        let client = XmtpClientInner::new(
            Some(signer),
            Some(identifier),
            opts,
        )?;
        
        let handle = Box::into_raw(Box::new(client)) as XmtpClientHandle;
        
        if !out_client.is_null() {
            unsafe { *out_client = handle };
        }
        
        Ok(())
    };
    
    match result() {
        Ok(()) => XmtpResult::ok(),
        Err(e) => XmtpResult::err(e.to_string()),
    }
}

/// Build a client from an existing identity (no signer required)
#[no_mangle]
pub extern "C" fn xmtp_client_build(
    identifier: XmtpIdentifier,
    opts: *const XmtpClientOptions,
    out_client: *mut XmtpClientHandle,
) -> XmtpResult {
    let result = || {
        let opts = if opts.is_null() {
            XmtpClientOptions::default()
        } else {
            unsafe { *opts }
        };
        
        let client = XmtpClientInner::new(None, Some(identifier), opts)?;
        let handle = Box::into_raw(Box::new(client)) as XmtpClientHandle;
        
        if !out_client.is_null() {
            unsafe { *out_client = handle };
        }
        
        Ok(())
    };
    
    match result() {
        Ok(()) => XmtpResult::ok(),
        Err(e) => XmtpResult::err(e.to_string()),
    }
}

/// Free a client handle
#[no_mangle]
pub extern "C" fn xmtp_client_free(client: XmtpClientHandle) {
    if !client.is_null() {
        unsafe {
            drop(Box::from_raw(client as *mut XmtpClientInner));
        }
    }
}

/// Get the client's inbox ID
#[no_mangle]
pub extern "C" fn xmtp_client_inbox_id(
    client: XmtpClientHandle,
    out_result: *mut XmtpStringResult,
) {
    if client.is_null() {
        if !out_result.is_null() {
            unsafe {
                *out_result = XmtpStringResult {
                    value: ptr::null_mut(),
                    error: Box::into_raw(Box::new(XmtpFfiError::InvalidArgument("null client".into()).to_ffi())),
                };
            }
        }
        return;
    }
    
    let client = unsafe { &*(client as *const XmtpClientInner) };
    let inbox_id = client.inbox_id.clone();
    
    let value = match CString::new(inbox_id) {
        Ok(s) => s.into_raw(),
        Err(_) => ptr::null_mut(),
    };
    
    if !out_result.is_null() {
        unsafe {
            *out_result = XmtpStringResult {
                value,
                error: ptr::null_mut(),
            };
        }
    }
}

/// Get the client's installation ID as bytes
#[no_mangle]
pub extern "C" fn xmtp_client_installation_id(
    client: XmtpClientHandle,
    out_data: *mut u8,
    out_len: *mut usize,
) -> XmtpResult {
    if client.is_null() {
        return XmtpResult::err("null client");
    }
    
    // TODO: Get actual installation ID
    let installation_id: [u8; 32] = [0u8; 32];
    
    if !out_len.is_null() {
        unsafe { *out_len = installation_id.len() };
    }
    
    if !out_data.is_null() {
        unsafe {
            std::ptr::copy_nonoverlapping(
                installation_id.as_ptr(),
                out_data,
                installation_id.len(),
            );
        }
    }
    
    XmtpResult::ok()
}

/// Check if the client is registered
#[no_mangle]
pub extern "C" fn xmtp_client_is_registered(client: XmtpClientHandle) -> bool {
    if client.is_null() {
        return false;
    }
    
    let _client = unsafe { &*(client as *const XmtpClientInner) };
    // TODO: Check actual registration status
    false
}

/// Get the conversations manager
#[no_mangle]
pub extern "C" fn xmtp_client_conversations(
    client: XmtpClientHandle,
    out_conversations: *mut XmtpConversationsHandle,
) -> XmtpResult {
    if client.is_null() {
        return XmtpResult::err("null client");
    }
    
    // TODO: Create actual conversations manager
    let conversations = Box::new(());
    let handle = Box::into_raw(conversations) as XmtpConversationsHandle;
    
    if !out_conversations.is_null() {
        unsafe { *out_conversations = handle };
    }
    
    XmtpResult::ok()
}

/// Register the client with the XMTP network
#[no_mangle]
pub extern "C" fn xmtp_client_register(
    client: XmtpClientHandle,
    signer_callback: XmtpSignerCallback,
    signer_user_data: *mut c_void,
) -> XmtpResult {
    if client.is_null() {
        return XmtpResult::err("null client");
    }
    
    let _client = unsafe { &*(client as *const XmtpClientInner) };
    
    // TODO: Implement actual registration
    XmtpResult::ok()
}

/// Check if identifiers can be messaged
#[no_mangle]
pub extern "C" fn xmtp_client_can_message(
    client: XmtpClientHandle,
    identifiers: *const XmtpIdentifier,
    identifiers_len: usize,
    out_results: *mut bool,
    out_len: *mut usize,
) -> XmtpResult {
    if client.is_null() {
        return XmtpResult::err("null client");
    }
    
    if identifiers.is_null() || identifiers_len == 0 {
        if !out_len.is_null() {
            unsafe { *out_len = 0 };
        }
        return XmtpResult::ok();
    }
    
    // TODO: Implement actual can_message check
    let _identifiers = unsafe { std::slice::from_raw_parts(identifiers, identifiers_len) };
    
    if !out_len.is_null() {
        unsafe { *out_len = identifiers_len };
    }
    
    if !out_results.is_null() {
        // For now, return all true
        for i in 0..identifiers_len {
            unsafe { *out_results.add(i) = true };
        }
    }
    
    XmtpResult::ok()
}

/// Get inbox ID by identifier
#[no_mangle]
pub extern "C" fn xmtp_client_get_inbox_id_by_identifier(
    client: XmtpClientHandle,
    identifier: XmtpIdentifier,
    out_result: *mut XmtpStringResult,
) {
    if client.is_null() {
        if !out_result.is_null() {
            unsafe {
                *out_result = XmtpStringResult {
                    value: ptr::null_mut(),
                    error: Box::into_raw(Box::new(XmtpFfiError::InvalidArgument("null client".into()).to_ffi())),
                };
            }
        }
        return;
    }
    
    // TODO: Implement actual lookup
    let value = CString::new("").unwrap().into_raw();
    
    if !out_result.is_null() {
        unsafe {
            *out_result = XmtpStringResult {
                value,
                error: ptr::null_mut(),
            };
        }
    }
}

/// Get libxmtp version
#[no_mangle]
pub extern "C" fn xmtp_client_libxmtp_version() -> *mut c_char {
    match CString::new(env!("CARGO_PKG_VERSION")) {
        Ok(s) => s.into_raw(),
        Err(_) => ptr::null_mut(),
    }
}
