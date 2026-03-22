//! Conversations FFI functions

use std::ffi::{c_char, c_void, CString};
use std::ptr;
use std::sync::Arc;

use crate::types::*;
use crate::error::XmtpError;

#[cfg(feature = "libxmtp")]
use crate::xmtp_client::{XmtpConversationsInner, XmtpConversationInner};

/// Free a conversations handle
#[no_mangle]
pub extern "C" fn xmtp_conversations_free(conversations: XmtpConversationsHandle) {
    if !conversations.is_null() {
        #[cfg(feature = "libxmtp")]
        unsafe {
            drop(Arc::from_raw(conversations as *const XmtpConversationsInner));
        }
    }
}

/// List all conversations
#[no_mangle]
pub extern "C" fn xmtp_conversations_list(
    conversations: XmtpConversationsHandle,
    conv_type: XmtpConversationType,
    out_conversations: *mut XmtpConversationHandle,
    out_len: *mut usize,
) -> XmtpResult {
    if conversations.is_null() {
        return XmtpResult::err("null conversations");
    }
    
    #[cfg(feature = "libxmtp")]
    {
        let conversations = unsafe { &*(conversations as *const XmtpConversationsInner) };
        
        match conversations.list(conv_type) {
            Ok(list) => {
                if !out_len.is_null() {
                    unsafe { *out_len = list.len() };
                }
                
                if !out_conversations.is_null() {
                    let handles: Vec<XmtpConversationHandle> = list
                        .into_iter()
                        .map(|c| Arc::into_raw(Arc::new(c)) as XmtpConversationHandle)
                        .collect();
                    
                    unsafe {
                        std::ptr::copy_nonoverlapping(
                            handles.as_ptr(),
                            out_conversations,
                            handles.len(),
                        );
                    }
                }
                
                XmtpResult::ok()
            }
            Err(e) => XmtpResult::err(e.to_string()),
        }
    }
    
    #[cfg(not(feature = "libxmtp"))]
    {
        if !out_len.is_null() {
            unsafe { *out_len = 0 };
        }
        XmtpResult::err("libxmtp not enabled")
    }
}

/// Get a conversation by ID
#[no_mangle]
pub extern "C" fn xmtp_conversations_get_by_id(
    conversations: XmtpConversationsHandle,
    id: *const c_char,
    out_conversation: *mut XmtpConversationHandle,
) -> XmtpResult {
    if conversations.is_null() {
        return XmtpResult::err("null conversations");
    }
    
    if id.is_null() {
        return XmtpResult::err("null id");
    }
    
    // TODO: Implement
    XmtpResult::err("not implemented")
}

/// Get a DM by inbox ID
#[no_mangle]
pub extern "C" fn xmtp_conversations_get_dm_by_inbox_id(
    conversations: XmtpConversationsHandle,
    inbox_id: *const c_char,
    out_conversation: *mut XmtpConversationHandle,
) -> XmtpResult {
    if conversations.is_null() {
        return XmtpResult::err("null conversations");
    }
    
    if inbox_id.is_null() {
        return XmtpResult::err("null inbox_id");
    }
    
    #[cfg(feature = "libxmtp")]
    {
        let conversations = unsafe { &*(conversations as *const XmtpConversationsInner) };
        let inbox_id = unsafe { std::ffi::CStr::from_ptr(inbox_id) }
            .to_str()
            .map_err(|e| XmtpError::InvalidArgument(e.to_string()));
        
        match inbox_id {
            Ok(id) => match conversations.create_dm(id) {
                Ok(conv) => {
                    if !out_conversation.is_null() {
                        unsafe { 
                            *out_conversation = Arc::into_raw(Arc::new(conv)) as XmtpConversationHandle;
                        };
                    }
                    XmtpResult::ok()
                }
                Err(e) => XmtpResult::err(e.to_string()),
            },
            Err(e) => XmtpResult::err(e.to_string()),
        }
    }
    
    #[cfg(not(feature = "libxmtp"))]
    {
        XmtpResult::err("libxmtp not enabled")
    }
}

/// Create a new group
#[no_mangle]
pub extern "C" fn xmtp_conversations_create_group(
    conversations: XmtpConversationsHandle,
    inbox_ids: *const *const c_char,
    inbox_ids_len: usize,
    name: *const c_char,
    image_url: *const c_char,
    description: *const c_char,
    out_group: *mut XmtpConversationHandle,
) -> XmtpResult {
    if conversations.is_null() {
        return XmtpResult::err("null conversations");
    }
    
    #[cfg(feature = "libxmtp")]
    {
        let conversations = unsafe { &*(conversations as *const XmtpConversationsInner) };
        
        let name = if name.is_null() {
            None
        } else {
            unsafe { std::ffi::CStr::from_ptr(name).to_str().ok() }
        };
        
        let description = if description.is_null() {
            None
        } else {
            unsafe { std::ffi::CStr::from_ptr(description).to_str().ok() }
        };
        
        match conversations.create_group(name, description) {
            Ok(group) => {
                if !out_group.is_null() {
                    unsafe {
                        *out_group = Arc::into_raw(Arc::new(group)) as XmtpConversationHandle;
                    }
                }
                XmtpResult::ok()
            }
            Err(e) => XmtpResult::err(e.to_string()),
        }
    }
    
    #[cfg(not(feature = "libxmtp"))]
    {
        XmtpResult::err("libxmtp not enabled")
    }
}

/// Create a new DM
#[no_mangle]
pub extern "C" fn xmtp_conversations_create_dm(
    conversations: XmtpConversationsHandle,
    inbox_id: *const c_char,
    out_dm: *mut XmtpConversationHandle,
) -> XmtpResult {
    xmtp_conversations_get_dm_by_inbox_id(conversations, inbox_id, out_dm)
}

/// Sync conversations from the network
#[no_mangle]
pub extern "C" fn xmtp_conversations_sync(
    conversations: XmtpConversationsHandle,
) -> XmtpResult {
    if conversations.is_null() {
        return XmtpResult::err("null conversations");
    }
    
    // TODO: Implement sync
    XmtpResult::ok()
}

/// Stream new conversations
#[no_mangle]
pub extern "C" fn xmtp_conversations_stream(
    conversations: XmtpConversationsHandle,
    conv_type: XmtpConversationType,
    callback: XmtpStreamCallback,
    error_callback: XmtpStreamErrorCallback,
    user_data: *mut c_void,
    out_stream: *mut XmtpStreamHandle,
) -> XmtpResult {
    if conversations.is_null() {
        return XmtpResult::err("null conversations");
    }
    
    // TODO: Implement streaming
    XmtpResult::err("streaming not yet implemented")
}

/// End a stream
#[no_mangle]
pub extern "C" fn xmtp_stream_end(stream: XmtpStreamHandle) {
    if !stream.is_null() {
        // TODO: Implement stream cleanup
    }
}
