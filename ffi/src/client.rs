//! Client FFI functions

use std::ffi::{c_char, c_void, CString};
use std::ptr;
use std::sync::Arc;

use crate::types::*;
use crate::error::XmtpError;
use crate::signer::FfiSigner;

#[cfg(feature = "libxmtp")]
use crate::xmtp_client::{XmtpClientInner, XmtpConversationsInner, XmtpConversationInner, XmtpMessageInner};

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
    let result: Result<(), XmtpError> = (|| {
        let opts = if opts.is_null() {
            XmtpClientOptions::default()
        } else {
            unsafe { *opts }
        };
        
        let signer = FfiSigner::new(signer_callback, signer_user_data);
        
        #[cfg(feature = "libxmtp")]
        {
            let ident = identifier_to_libxmtp(&identifier)?;
            let client = XmtpClientInner::new(signer, ident, opts)?;
            let handle = Box::into_raw(Box::new(client)) as XmtpClientHandle;
            
            if !out_client.is_null() {
                unsafe { *out_client = handle };
            }
        }
        
        #[cfg(not(feature = "libxmtp"))]
        {
            let _ = (signer, opts);
            return Err(XmtpError::Generic("libxmtp feature not enabled".into()));
        }
        
        Ok(())
    })();
    
    match result {
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
    let result: Result<(), XmtpError> = (|| {
        #[cfg(feature = "libxmtp")]
        {
            let opts = if opts.is_null() {
                XmtpClientOptions::default()
            } else {
                unsafe { *opts }
            };
            
            // Create a dummy signer for building
            let signer = FfiSigner::empty();
            let ident = identifier_to_libxmtp(&identifier)?;
            let client = XmtpClientInner::new(signer, ident, opts)?;
            let handle = Box::into_raw(Box::new(client)) as XmtpClientHandle;
            
            if !out_client.is_null() {
                unsafe { *out_client = handle };
            }
        }
        
        #[cfg(not(feature = "libxmtp"))]
        {
            let _ = (identifier, opts, out_client);
            return Err(XmtpError::Generic("libxmtp feature not enabled".into()));
        }
        
        Ok(())
    })();
    
    match result {
        Ok(()) => XmtpResult::ok(),
        Err(e) => XmtpResult::err(e.to_string()),
    }
}

/// Free a client handle
#[no_mangle]
pub extern "C" fn xmtp_client_free(client: XmtpClientHandle) {
    if !client.is_null() {
        #[cfg(feature = "libxmtp")]
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
                    error: Box::into_raw(Box::new(XmtpError::InvalidArgument("null client".into()).to_ffi())),
                };
            }
        }
        return;
    }
    
    #[cfg(feature = "libxmtp")]
    {
        let client = unsafe { &*(client as *const XmtpClientInner) };
        let inbox_id = client.inbox_id().to_string();
        
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
    
    #[cfg(not(feature = "libxmtp"))]
    {
        if !out_result.is_null() {
            unsafe {
                *out_result = XmtpStringResult {
                    value: ptr::null_mut(),
                    error: Box::into_raw(Box::new(XmtpError::Generic("libxmtp not enabled".into()).to_ffi())),
                };
            }
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
    
    #[cfg(feature = "libxmtp")]
    {
        let client = unsafe { &*(client as *const XmtpClientInner) };
        let installation_id = client.installation_id();
        
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
    
    #[cfg(not(feature = "libxmtp"))]
    {
        XmtpResult::err("libxmtp not enabled")
    }
}

/// Check if the client is registered
#[no_mangle]
pub extern "C" fn xmtp_client_is_registered(client: XmtpClientHandle) -> bool {
    if client.is_null() {
        return false;
    }
    
    #[cfg(feature = "libxmtp")]
    {
        let client = unsafe { &*(client as *const XmtpClientInner) };
        client.is_registered()
    }
    
    #[cfg(not(feature = "libxmtp"))]
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
    
    #[cfg(feature = "libxmtp")]
    {
        let client = unsafe { &*(client as *const XmtpClientInner) };
        let conversations = client.conversations();
        let handle = Arc::into_raw(conversations) as XmtpConversationsHandle;
        
        if !out_conversations.is_null() {
            unsafe { *out_conversations = handle };
        }
        
        XmtpResult::ok()
    }
    
    #[cfg(not(feature = "libxmtp"))]
    {
        XmtpResult::err("libxmtp not enabled")
    }
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
    
    #[cfg(feature = "libxmtp")]
    {
        let client = unsafe { &*(client as *const XmtpClientInner) };
        let signer = FfiSigner::new(signer_callback, signer_user_data);
        
        match client.register(&signer) {
            Ok(()) => XmtpResult::ok(),
            Err(e) => XmtpResult::err(e.to_string()),
        }
    }
    
    #[cfg(not(feature = "libxmtp"))]
    {
        XmtpResult::err("libxmtp not enabled")
    }
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
    
    // TODO: Implement with libxmtp
    if !out_len.is_null() {
        unsafe { *out_len = identifiers_len };
    }
    
    if !out_results.is_null() {
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
                    error: Box::into_raw(Box::new(XmtpError::InvalidArgument("null client".into()).to_ffi())),
                };
            }
        }
        return;
    }
    
    // TODO: Implement with libxmtp
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

#[cfg(feature = "libxmtp")]
fn identifier_to_libxmtp(ident: &XmtpIdentifier) -> Result<xmtp_id::associations::Identifier, XmtpError> {
    let identifier_str = unsafe {
        if ident.identifier.is_null() {
            return Err(XmtpError::InvalidArgument("null identifier".into()));
        }
        std::ffi::CStr::from_ptr(ident.identifier)
            .to_str()
            .map_err(|e| XmtpError::InvalidArgument(e.to_string()))?
    };
    
    match ident.kind {
        XmtpIdentifierKind::Ethereum => {
            xmtp_id::associations::Identifier::parse_ethereum(identifier_str)
                .map_err(|e| XmtpError::InvalidArgument(e.to_string()))
        }
        XmtpIdentifierKind::Passkey => {
            Err(XmtpError::Generic("Passkey not yet supported".into()))
        }
    }
}
