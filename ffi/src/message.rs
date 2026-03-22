//! Message FFI functions

use std::ffi::{c_char, c_void, CStr, CString};
use std::ptr;

use crate::types::*;
use crate::error::XmtpFfiError;

/// Free a message handle
#[no_mangle]
pub extern "C" fn xmtp_message_free(message: XmtpMessageHandle) {
    if !message.is_null() {
        unsafe {
            drop(Box::from_raw(message as *mut XmtpMessageFfi));
        }
    }
}

/// Get the message ID
#[no_mangle]
pub extern "C" fn xmtp_message_id(message: XmtpMessageHandle) -> *mut c_char {
    if message.is_null() {
        return ptr::null_mut();
    }
    
    let msg = unsafe { &*(message as *const XmtpMessageFfi) };
    msg.id
}

/// Get the sender's inbox ID
#[no_mangle]
pub extern "C" fn xmtp_message_sender_inbox_id(message: XmtpMessageHandle) -> *mut c_char {
    if message.is_null() {
        return ptr::null_mut();
    }
    
    let msg = unsafe { &*(message as *const XmtpMessageFfi) };
    msg.sender_inbox_id
}

/// Get the message timestamp (nanoseconds)
#[no_mangle]
pub extern "C" fn xmtp_message_sent_at_ns(message: XmtpMessageHandle) -> u64 {
    if message.is_null() {
        return 0;
    }
    
    let msg = unsafe { &*(message as *const XmtpMessageFfi) };
    msg.sent_at_ns
}

/// Get the message expiry timestamp (nanoseconds)
#[no_mangle]
pub extern "C" fn xmtp_message_expires_at_ns(message: XmtpMessageHandle) -> u64 {
    if message.is_null() {
        return 0;
    }
    
    let msg = unsafe { &*(message as *const XmtpMessageFfi) };
    msg.expires_at_ns
}

/// Get the conversation ID
#[no_mangle]
pub extern "C" fn xmtp_message_conversation_id(message: XmtpMessageHandle) -> *mut c_char {
    if message.is_null() {
        return ptr::null_mut();
    }
    
    let msg = unsafe { &*(message as *const XmtpMessageFfi) };
    msg.conversation_id
}

/// Get the content type
#[no_mangle]
pub extern "C" fn xmtp_message_content_type(message: XmtpMessageHandle) -> XmtpContentType {
    if message.is_null() {
        return XmtpContentType::Text;
    }
    
    let msg = unsafe { &*(message as *const XmtpMessageFfi) };
    msg.content_type
}

/// Get the delivery status
#[no_mangle]
pub extern "C" fn xmtp_message_delivery_status(message: XmtpMessageHandle) -> XmtpDeliveryStatus {
    if message.is_null() {
        return XmtpDeliveryStatus::Unpublished;
    }
    
    let msg = unsafe { &*(message as *const XmtpMessageFfi) };
    msg.delivery_status
}

/// Get the fallback content
#[no_mangle]
pub extern "C" fn xmtp_message_fallback(message: XmtpMessageHandle) -> *mut c_char {
    if message.is_null() {
        return ptr::null_mut();
    }
    
    let msg = unsafe { &*(message as *const XmtpMessageFfi) };
    msg.fallback
}

/// Get text content (if text message)
#[no_mangle]
pub extern "C" fn xmtp_message_content_text(message: XmtpMessageHandle) -> *mut c_char {
    if message.is_null() {
        return ptr::null_mut();
    }
    
    let msg = unsafe { &*(message as *const XmtpMessageFfi) };
    
    if msg.content_type != XmtpContentType::Text {
        return ptr::null_mut();
    }
    
    // TODO: Get actual text content
    CString::new("").unwrap().into_raw()
}

/// Get markdown content (if markdown message)
#[no_mangle]
pub extern "C" fn xmtp_message_content_markdown(message: XmtpMessageHandle) -> *mut c_char {
    if message.is_null() {
        return ptr::null_mut();
    }
    
    let msg = unsafe { &*(message as *const XmtpMessageFfi) };
    
    if msg.content_type != XmtpContentType::Markdown {
        return ptr::null_mut();
    }
    
    // TODO: Get actual markdown content
    CString::new("").unwrap().into_raw()
}

/// Reaction content result
#[repr(C)]
pub struct XmtpReactionContent {
    pub reference_message_id: *mut c_char,
    pub action: c_int,
    pub schema: c_int,
    pub content: *mut c_char,
}

/// Get reaction content (if reaction message)
#[no_mangle]
pub extern "C" fn xmtp_message_content_reaction(
    message: XmtpMessageHandle,
    out_reaction: *mut XmtpReactionContent,
) {
    if message.is_null() || out_reaction.is_null() {
        return;
    }
    
    let msg = unsafe { &*(message as *const XmtpMessageFfi) };
    
    if msg.content_type != XmtpContentType::Reaction {
        return;
    }
    
    // TODO: Get actual reaction content
    unsafe {
        *out_reaction = XmtpReactionContent {
            reference_message_id: CString::new("").unwrap().into_raw(),
            action: 0,
            schema: 0,
            content: CString::new("").unwrap().into_raw(),
        };
    }
}

/// Free reaction content
#[no_mangle]
pub extern "C" fn xmtp_reaction_content_free(reaction: *mut XmtpReactionContent) {
    if !reaction.is_null() {
        unsafe {
            let r = &mut *reaction;
            if !r.reference_message_id.is_null() {
                drop(CString::from_raw(r.reference_message_id));
            }
            if !r.content.is_null() {
                drop(CString::from_raw(r.content));
            }
        }
    }
}

/// Reply content result
#[repr(C)]
pub struct XmtpReplyContent {
    pub reference_message_id: *mut c_char,
    pub content_type: XmtpContentType,
    pub content: *mut c_char,
}

/// Get reply content (if reply message)
#[no_mangle]
pub extern "C" fn xmtp_message_content_reply(
    message: XmtpMessageHandle,
    out_reply: *mut XmtpReplyContent,
) {
    if message.is_null() || out_reply.is_null() {
        return;
    }
    
    let msg = unsafe { &*(message as *const XmtpMessageFfi) };
    
    if msg.content_type != XmtpContentType::Reply {
        return;
    }
    
    // TODO: Get actual reply content
    unsafe {
        *out_reply = XmtpReplyContent {
            reference_message_id: CString::new("").unwrap().into_raw(),
            content_type: XmtpContentType::Text,
            content: CString::new("").unwrap().into_raw(),
        };
    }
}

/// Free reply content
#[no_mangle]
pub extern "C" fn xmtp_reply_content_free(reply: *mut XmtpReplyContent) {
    if !reply.is_null() {
        unsafe {
            let r = &mut *reply;
            if !r.reference_message_id.is_null() {
                drop(CString::from_raw(r.reference_message_id));
            }
            if !r.content.is_null() {
                drop(CString::from_raw(r.content));
            }
        }
    }
}

/// Attachment content result
#[repr(C)]
pub struct XmtpAttachmentContent {
    pub filename: *mut c_char,
    pub mime_type: *mut c_char,
    pub data: *mut u8,
    pub data_len: usize,
}

/// Get attachment content (if attachment message)
#[no_mangle]
pub extern "C" fn xmtp_message_content_attachment(
    message: XmtpMessageHandle,
    out_attachment: *mut XmtpAttachmentContent,
) {
    if message.is_null() || out_attachment.is_null() {
        return;
    }
    
    let msg = unsafe { &*(message as *const XmtpMessageFfi) };
    
    if msg.content_type != XmtpContentType::Attachment {
        return;
    }
    
    // TODO: Get actual attachment content
    unsafe {
        *out_attachment = XmtpAttachmentContent {
            filename: CString::new("").unwrap().into_raw(),
            mime_type: CString::new("").unwrap().into_raw(),
            data: ptr::null_mut(),
            data_len: 0,
        };
    }
}

/// Free attachment content
#[no_mangle]
pub extern "C" fn xmtp_attachment_content_free(attachment: *mut XmtpAttachmentContent) {
    if !attachment.is_null() {
        unsafe {
            let a = &mut *attachment;
            if !a.filename.is_null() {
                drop(CString::from_raw(a.filename));
            }
            if !a.mime_type.is_null() {
                drop(CString::from_raw(a.mime_type));
            }
            if !a.data.is_null() && a.data_len > 0 {
                drop(Vec::from_raw_parts(a.data, a.data_len, a.data_len));
            }
        }
    }
}

/// Get raw content bytes
#[no_mangle]
pub extern "C" fn xmtp_message_content_bytes(
    message: XmtpMessageHandle,
    out_data: *mut u8,
    out_len: *mut usize,
) -> XmtpResult {
    if message.is_null() {
        return XmtpResult::err("null message");
    }
    
    let msg = unsafe { &*(message as *const XmtpMessageFfi) };
    
    if !out_len.is_null() {
        unsafe { *out_len = msg.content_len };
    }
    
    if !out_data.is_null() && msg.content_len > 0 && !msg.content_data.is_null() {
        unsafe {
            std::ptr::copy_nonoverlapping(
                msg.content_data,
                out_data,
                msg.content_len,
            );
        }
    }
    
    XmtpResult::ok()
}
