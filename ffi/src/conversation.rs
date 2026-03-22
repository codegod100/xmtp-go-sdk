//! Conversation FFI functions

use std::ffi::{c_char, c_void, CStr, CString};
use std::ptr;

use crate::types::*;
use crate::error::XmtpFfiError;

/// Free a conversation handle
#[no_mangle]
pub extern "C" fn xmtp_conversation_free(conversation: XmtpConversationHandle) {
    if !conversation.is_null() {
        unsafe {
            drop(Box::from_raw(conversation as *mut ()));
        }
    }
}

/// Get the conversation ID
#[no_mangle]
pub extern "C" fn xmtp_conversation_id(
    conversation: XmtpConversationHandle,
    out_result: *mut XmtpStringResult,
) {
    if conversation.is_null() {
        if !out_result.is_null() {
            unsafe {
                *out_result = XmtpStringResult {
                    value: ptr::null_mut(),
                    error: Box::into_raw(Box::new(XmtpFfiError::InvalidArgument("null conversation".into()).to_ffi())),
                };
            }
        }
        return;
    }
    
    // TODO: Get actual ID
    let id = CString::new("").unwrap().into_raw();
    
    if !out_result.is_null() {
        unsafe {
            *out_result = XmtpStringResult {
                value: id,
                error: ptr::null_mut(),
            };
        }
    }
}

/// Check if the conversation is active
#[no_mangle]
pub extern "C" fn xmtp_conversation_is_active(conversation: XmtpConversationHandle) -> bool {
    if conversation.is_null() {
        return false;
    }
    
    // TODO: Check actual active state
    true
}

/// Get the conversation creation timestamp (nanoseconds)
#[no_mangle]
pub extern "C" fn xmtp_conversation_created_at_ns(conversation: XmtpConversationHandle) -> u64 {
    if conversation.is_null() {
        return 0;
    }
    
    // TODO: Get actual timestamp
    0
}

/// Get the conversation consent state
#[no_mangle]
pub extern "C" fn xmtp_conversation_consent_state(conversation: XmtpConversationHandle) -> XmtpConsentState {
    if conversation.is_null() {
        return XmtpConsentState::Unknown;
    }
    
    // TODO: Get actual consent state
    XmtpConsentState::Unknown
}

/// Update the conversation consent state
#[no_mangle]
pub extern "C" fn xmtp_conversation_update_consent(
    conversation: XmtpConversationHandle,
    state: XmtpConsentState,
) -> XmtpResult {
    if conversation.is_null() {
        return XmtpResult::err("null conversation");
    }
    
    // TODO: Implement actual update
    XmtpResult::ok()
}

/// Sync the conversation from the network
#[no_mangle]
pub extern "C" fn xmtp_conversation_sync(conversation: XmtpConversationHandle) -> XmtpResult {
    if conversation.is_null() {
        return XmtpResult::err("null conversation");
    }
    
    // TODO: Implement actual sync
    XmtpResult::ok()
}

/// Send a text message
#[no_mangle]
pub extern "C" fn xmtp_conversation_send_text(
    conversation: XmtpConversationHandle,
    text: *const c_char,
    optimistic: bool,
    out_result: *mut XmtpStringResult,
) {
    if conversation.is_null() {
        if !out_result.is_null() {
            unsafe {
                *out_result = XmtpStringResult {
                    value: ptr::null_mut(),
                    error: Box::into_raw(Box::new(XmtpFfiError::InvalidArgument("null conversation".into()).to_ffi())),
                };
            }
        }
        return;
    }
    
    if text.is_null() {
        if !out_result.is_null() {
            unsafe {
                *out_result = XmtpStringResult {
                    value: ptr::null_mut(),
                    error: Box::into_raw(Box::new(XmtpFfiError::InvalidArgument("null text".into()).to_ffi())),
                };
            }
        }
        return;
    }
    
    // TODO: Implement actual send
    let message_id = CString::new("").unwrap().into_raw();
    
    if !out_result.is_null() {
        unsafe {
            *out_result = XmtpStringResult {
                value: message_id,
                error: ptr::null_mut(),
            };
        }
    }
}

/// Send a markdown message
#[no_mangle]
pub extern "C" fn xmtp_conversation_send_markdown(
    conversation: XmtpConversationHandle,
    markdown: *const c_char,
    optimistic: bool,
    out_result: *mut XmtpStringResult,
) {
    if conversation.is_null() || markdown.is_null() {
        if !out_result.is_null() {
            unsafe {
                *out_result = XmtpStringResult {
                    value: ptr::null_mut(),
                    error: Box::into_raw(Box::new(XmtpFfiError::InvalidArgument("null argument".into()).to_ffi())),
                };
            }
        }
        return;
    }
    
    // TODO: Implement actual send
    let message_id = CString::new("").unwrap().into_raw();
    
    if !out_result.is_null() {
        unsafe {
            *out_result = XmtpStringResult {
                value: message_id,
                error: ptr::null_mut(),
            };
        }
    }
}

/// Send a reaction
#[no_mangle]
pub extern "C" fn xmtp_conversation_send_reaction(
    conversation: XmtpConversationHandle,
    reference_message_id: *const c_char,
    action: c_int, // 0 = add, 1 = remove
    schema: c_int, // 0 = unicode, 1 = custom
    content: *const c_char,
    optimistic: bool,
    out_result: *mut XmtpStringResult,
) {
    if conversation.is_null() || reference_message_id.is_null() || content.is_null() {
        if !out_result.is_null() {
            unsafe {
                *out_result = XmtpStringResult {
                    value: ptr::null_mut(),
                    error: Box::into_raw(Box::new(XmtpFfiError::InvalidArgument("null argument".into()).to_ffi())),
                };
            }
        }
        return;
    }
    
    // TODO: Implement actual send
    let message_id = CString::new("").unwrap().into_raw();
    
    if !out_result.is_null() {
        unsafe {
            *out_result = XmtpStringResult {
                value: message_id,
                error: ptr::null_mut(),
            };
        }
    }
}

/// List messages in the conversation
#[no_mangle]
pub extern "C" fn xmtp_conversation_list_messages(
    conversation: XmtpConversationHandle,
    opts: *const XmtpListMessagesOptions,
    out_messages: *mut XmtpMessageHandle,
    out_len: *mut usize,
) -> XmtpResult {
    if conversation.is_null() {
        return XmtpResult::err("null conversation");
    }
    
    let opts = if opts.is_null() {
        XmtpListMessagesOptions::default()
    } else {
        unsafe { *opts }
    };
    
    // TODO: Implement actual listing
    if !out_len.is_null() {
        unsafe { *out_len = 0 };
    }
    
    XmtpResult::ok()
}

/// Get a message by ID
#[no_mangle]
pub extern "C" fn xmtp_conversation_get_message_by_id(
    conversation: XmtpConversationHandle,
    message_id: *const c_char,
    out_message: *mut XmtpMessageHandle,
) -> XmtpResult {
    if conversation.is_null() {
        return XmtpResult::err("null conversation");
    }
    
    if message_id.is_null() {
        return XmtpResult::err("null message_id");
    }
    
    // TODO: Implement actual lookup
    XmtpResult::err("not found")
}

/// Stream messages in the conversation
#[no_mangle]
pub extern "C" fn xmtp_conversation_stream_messages(
    conversation: XmtpConversationHandle,
    callback: XmtpStreamCallback,
    error_callback: XmtpStreamErrorCallback,
    user_data: *mut c_void,
    out_stream: *mut XmtpStreamHandle,
) -> XmtpResult {
    if conversation.is_null() {
        return XmtpResult::err("null conversation");
    }
    
    // TODO: Implement actual streaming
    XmtpResult::err("not implemented")
}

// Group-specific functions

/// Get the group name
#[no_mangle]
pub extern "C" fn xmtp_group_name(conversation: XmtpConversationHandle) -> *mut c_char {
    if conversation.is_null() {
        return ptr::null_mut();
    }
    
    // TODO: Get actual name
    CString::new("").unwrap().into_raw()
}

/// Update the group name
#[no_mangle]
pub extern "C" fn xmtp_group_update_name(
    conversation: XmtpConversationHandle,
    name: *const c_char,
) -> XmtpResult {
    if conversation.is_null() || name.is_null() {
        return XmtpResult::err("null argument");
    }
    
    // TODO: Implement actual update
    XmtpResult::ok()
}

/// Get the group image URL
#[no_mangle]
pub extern "C" fn xmtp_group_image_url(conversation: XmtpConversationHandle) -> *mut c_char {
    if conversation.is_null() {
        return ptr::null_mut();
    }
    
    // TODO: Get actual URL
    CString::new("").unwrap().into_raw()
}

/// Update the group image URL
#[no_mangle]
pub extern "C" fn xmtp_group_update_image_url(
    conversation: XmtpConversationHandle,
    url: *const c_char,
) -> XmtpResult {
    if conversation.is_null() || url.is_null() {
        return XmtpResult::err("null argument");
    }
    
    // TODO: Implement actual update
    XmtpResult::ok()
}

/// Get the group description
#[no_mangle]
pub extern "C" fn xmtp_group_description(conversation: XmtpConversationHandle) -> *mut c_char {
    if conversation.is_null() {
        return ptr::null_mut();
    }
    
    // TODO: Get actual description
    CString::new("").unwrap().into_raw()
}

/// Update the group description
#[no_mangle]
pub extern "C" fn xmtp_group_update_description(
    conversation: XmtpConversationHandle,
    description: *const c_char,
) -> XmtpResult {
    if conversation.is_null() || description.is_null() {
        return XmtpResult::err("null argument");
    }
    
    // TODO: Implement actual update
    XmtpResult::ok()
}

/// List group members
#[no_mangle]
pub extern "C" fn xmtp_group_list_members(
    conversation: XmtpConversationHandle,
    out_members: *mut XmtpGroupMember,
    out_len: *mut usize,
) -> XmtpResult {
    if conversation.is_null() {
        return XmtpResult::err("null conversation");
    }
    
    // TODO: Implement actual listing
    if !out_len.is_null() {
        unsafe { *out_len = 0 };
    }
    
    XmtpResult::ok()
}

/// Add members to the group
#[no_mangle]
pub extern "C" fn xmtp_group_add_members(
    conversation: XmtpConversationHandle,
    inbox_ids: *const *const c_char,
    inbox_ids_len: usize,
) -> XmtpResult {
    if conversation.is_null() {
        return XmtpResult::err("null conversation");
    }
    
    if inbox_ids.is_null() || inbox_ids_len == 0 {
        return XmtpResult::err("no inbox IDs provided");
    }
    
    // TODO: Implement actual add
    XmtpResult::ok()
}

/// Remove members from the group
#[no_mangle]
pub extern "C" fn xmtp_group_remove_members(
    conversation: XmtpConversationHandle,
    inbox_ids: *const *const c_char,
    inbox_ids_len: usize,
) -> XmtpResult {
    if conversation.is_null() {
        return XmtpResult::err("null conversation");
    }
    
    if inbox_ids.is_null() || inbox_ids_len == 0 {
        return XmtpResult::err("no inbox IDs provided");
    }
    
    // TODO: Implement actual remove
    XmtpResult::ok()
}

/// Check if an inbox ID is an admin
#[no_mangle]
pub extern "C" fn xmtp_group_is_admin(
    conversation: XmtpConversationHandle,
    inbox_id: *const c_char,
) -> bool {
    if conversation.is_null() || inbox_id.is_null() {
        return false;
    }
    
    // TODO: Implement actual check
    false
}

/// Check if an inbox ID is a super admin
#[no_mangle]
pub extern "C" fn xmtp_group_is_super_admin(
    conversation: XmtpConversationHandle,
    inbox_id: *const c_char,
) -> bool {
    if conversation.is_null() || inbox_id.is_null() {
        return false;
    }
    
    // TODO: Implement actual check
    false
}

/// Add an admin to the group
#[no_mangle]
pub extern "C" fn xmtp_group_add_admin(
    conversation: XmtpConversationHandle,
    inbox_id: *const c_char,
) -> XmtpResult {
    if conversation.is_null() || inbox_id.is_null() {
        return XmtpResult::err("null argument");
    }
    
    // TODO: Implement actual add
    XmtpResult::ok()
}

/// Remove an admin from the group
#[no_mangle]
pub extern "C" fn xmtp_group_remove_admin(
    conversation: XmtpConversationHandle,
    inbox_id: *const c_char,
) -> XmtpResult {
    if conversation.is_null() || inbox_id.is_null() {
        return XmtpResult::err("null argument");
    }
    
    // TODO: Implement actual remove
    XmtpResult::ok()
}

/// Leave the group
#[no_mangle]
pub extern "C" fn xmtp_group_leave(conversation: XmtpConversationHandle) -> XmtpResult {
    if conversation.is_null() {
        return XmtpResult::err("null conversation");
    }
    
    // TODO: Implement actual leave
    XmtpResult::ok()
}

// DM-specific functions

/// Get the peer's inbox ID for a DM
#[no_mangle]
pub extern "C" fn xmtp_dm_peer_inbox_id(conversation: XmtpConversationHandle) -> *mut c_char {
    if conversation.is_null() {
        return ptr::null_mut();
    }
    
    // TODO: Get actual peer inbox ID
    CString::new("").unwrap().into_raw()
}
