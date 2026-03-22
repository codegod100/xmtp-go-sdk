//! Conversations FFI functions

use std::ffi::{c_char, c_void, CStr, CString};
use std::ptr;

use crate::types::*;
use crate::error::XmtpFfiError;

/// Free a conversations handle
#[no_mangle]
pub extern "C" fn xmtp_conversations_free(conversations: XmtpConversationsHandle) {
    if !conversations.is_null() {
        unsafe {
            drop(Box::from_raw(conversations as *mut ()));
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
    
    // TODO: Implement actual listing
    if !out_len.is_null() {
        unsafe { *out_len = 0 };
    }
    
    XmtpResult::ok()
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
    
    // TODO: Implement actual lookup
    XmtpResult::err("not found")
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
    
    // TODO: Implement actual lookup
    XmtpResult::err("not found")
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
    
    // TODO: Implement actual group creation
    XmtpResult::err("not implemented")
}

/// Create a new DM
#[no_mangle]
pub extern "C" fn xmtp_conversations_create_dm(
    conversations: XmtpConversationsHandle,
    inbox_id: *const c_char,
    out_dm: *mut XmtpConversationHandle,
) -> XmtpResult {
    if conversations.is_null() {
        return XmtpResult::err("null conversations");
    }
    
    if inbox_id.is_null() {
        return XmtpResult::err("null inbox_id");
    }
    
    // TODO: Implement actual DM creation
    XmtpResult::err("not implemented")
}

/// Sync conversations from the network
#[no_mangle]
pub extern "C" fn xmtp_conversations_sync(
    conversations: XmtpConversationsHandle,
) -> XmtpResult {
    if conversations.is_null() {
        return XmtpResult::err("null conversations");
    }
    
    // TODO: Implement actual sync
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
    
    // TODO: Implement actual streaming
    XmtpResult::err("not implemented")
}

/// End a stream
#[no_mangle]
pub extern "C" fn xmtp_stream_end(stream: XmtpStreamHandle) {
    if !stream.is_null() {
        unsafe {
            drop(Box::from_raw(stream as *mut ()));
        }
    }
}
