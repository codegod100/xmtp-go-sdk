//! FFI type definitions

use std::ffi::{c_char, c_int, c_void};
use std::ptr;

/// Opaque handle to an XMTP Client
pub type XmtpClientHandle = *mut c_void;

/// Opaque handle to a Conversations manager
pub type XmtpConversationsHandle = *mut c_void;

/// Opaque handle to a Conversation
pub type XmtpConversationHandle = *mut c_void;

/// Opaque handle to a Message
pub type XmtpMessageHandle = *mut c_void;

/// Opaque handle to a Stream
pub type XmtpStreamHandle = *mut c_void;

/// Opaque handle to a Signature Request
pub type XmtpSignatureRequestHandle = *mut c_void;

/// Environment type
#[repr(C)]
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum XmtpEnv {
    Local = 0,
    Dev = 1,
    Production = 2,
    TestnetStaging = 3,
    TestnetDev = 4,
    Testnet = 5,
    Mainnet = 6,
}

impl Default for XmtpEnv {
    fn default() -> Self {
        XmtpEnv::Dev
    }
}

/// Conversation type
#[repr(C)]
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum XmtpConversationType {
    Dm = 0,
    Group = 1,
}

/// Consent state
#[repr(C)]
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum XmtpConsentState {
    Unknown = 0,
    Allowed = 1,
    Denied = 2,
}

/// Content type
#[repr(C)]
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum XmtpContentType {
    Text = 0,
    Markdown = 1,
    Reply = 2,
    Reaction = 3,
    Attachment = 4,
    RemoteAttachment = 5,
    MultiRemoteAttachment = 6,
    TransactionReference = 7,
    GroupUpdated = 8,
    ReadReceipt = 9,
    LeaveRequest = 10,
    WalletSendCalls = 11,
    Actions = 12,
    Intent = 13,
    DeletedMessage = 14,
    Custom = 15,
}

/// Delivery status
#[repr(C)]
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum XmtpDeliveryStatus {
    Unpublished = 0,
    Published = 1,
    Failed = 2,
}

/// Identifier type
#[repr(C)]
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub enum XmtpIdentifierKind {
    Ethereum = 0,
    Passkey = 1,
}

/// FFI result wrapper
#[repr(C)]
pub struct XmtpResult {
    pub error: *mut XmtpFfiError,
}

impl XmtpResult {
    pub fn ok() -> Self {
        XmtpResult { error: ptr::null_mut() }
    }

    pub fn err(msg: impl Into<String>) -> Self {
        XmtpResult {
            error: Box::into_raw(Box::new(XmtpFfiError::Generic(msg.into()))),
        }
    }

    pub fn is_ok(&self) -> bool {
        self.error.is_null()
    }
}

impl Drop for XmtpResult {
    fn drop(&mut self) {
        if !self.error.is_null() {
            unsafe {
                drop(Box::from_raw(self.error));
            }
        }
    }
}

/// FFI Error type
#[repr(C)]
pub struct XmtpFfiError {
    pub code: c_int,
    pub message: *mut c_char,
}

impl XmtpFfiError {
    pub fn new(code: c_int, message: impl Into<String>) -> Self {
        let msg = message.into();
        let c_msg = std::ffi::CString::new(msg).unwrap_or_default();
        XmtpFfiError {
            code,
            message: c_msg.into_raw(),
        }
    }

    pub fn generic(message: impl Into<String>) -> Self {
        Self::new(-1, message)
    }
}

impl Drop for XmtpFfiError {
    fn drop(&mut self) {
        if !self.message.is_null() {
            unsafe {
                drop(std::ffi::CString::from_raw(self.message));
            }
        }
    }
}

/// Client options
#[repr(C)]
pub struct XmtpClientOptions {
    pub env: XmtpEnv,
    pub db_path: *const c_char,
    pub db_encryption_key: *const u8,
    pub db_encryption_key_len: usize,
    pub app_version: *const c_char,
    pub disable_auto_register: bool,
    pub structured_logging: bool,
    pub log_level: c_int,
}

impl Default for XmtpClientOptions {
    fn default() -> Self {
        XmtpClientOptions {
            env: XmtpEnv::Dev,
            db_path: ptr::null(),
            db_encryption_key: ptr::null(),
            db_encryption_key_len: 0,
            app_version: ptr::null(),
            disable_auto_register: false,
            structured_logging: false,
            log_level: 2, // Warn
        }
    }
}

/// String result
#[repr(C)]
pub struct XmtpStringResult {
    pub value: *mut c_char,
    pub error: *mut XmtpFfiError,
}

/// Bytes result
#[repr(C)]
pub struct XmtpBytesResult {
    pub data: *mut u8,
    pub len: usize,
    pub error: *mut XmtpFfiError,
}

/// Bool result
#[repr(C)]
pub struct XmtpBoolResult {
    pub value: bool,
    pub error: *mut XmtpFfiError,
}

/// Int result
#[repr(C)]
pub struct XmtpIntResult {
    pub value: c_int,
    pub error: *mut XmtpFfiError,
}

/// Identifier
#[repr(C)]
pub struct XmtpIdentifier {
    pub kind: XmtpIdentifierKind,
    pub identifier: *const c_char,
}

/// Signer callback type
pub type XmtpSignerCallback = extern "C" fn(
    message: *const u8,
    message_len: usize,
    user_data: *mut c_void,
    out_signature: *mut u8,
    out_signature_len: *mut usize,
) -> c_int;

/// Signer callback data
#[repr(C)]
pub struct XmtpSignerCallbackData {
    pub callback: XmtpSignerCallback,
    pub user_data: *mut c_void,
}

/// List messages options
#[repr(C)]
pub struct XmtpListMessagesOptions {
    pub limit: usize,
    pub before_ns: u64,
    pub after_ns: u64,
    pub ascending: bool,
}

impl Default for XmtpListMessagesOptions {
    fn default() -> Self {
        XmtpListMessagesOptions {
            limit: 100,
            before_ns: 0,
            after_ns: 0,
            ascending: false,
        }
    }
}

/// Create group options
#[repr(C)]
pub struct XmtpCreateGroupOptions {
    pub name: *const c_char,
    pub image_url: *const c_char,
    pub description: *const c_char,
    pub permissions: c_int,
}

impl Default for XmtpCreateGroupOptions {
    fn default() -> Self {
        XmtpCreateGroupOptions {
            name: ptr::null(),
            image_url: ptr::null(),
            description: ptr::null(),
            permissions: 0,
        }
    }
}

/// Group member
#[repr(C)]
pub struct XmtpGroupMember {
    pub inbox_id: *mut c_char,
    pub permission_level: c_int, // 0=member, 1=admin, 2=super_admin
}

/// Message FFI representation
#[repr(C)]
pub struct XmtpMessageFfi {
    pub id: *mut c_char,
    pub sender_inbox_id: *mut c_char,
    pub conversation_id: *mut c_char,
    pub sent_at_ns: u64,
    pub expires_at_ns: u64,
    pub content_type: XmtpContentType,
    pub delivery_status: XmtpDeliveryStatus,
    pub fallback: *mut c_char,
    pub content_data: *mut u8,
    pub content_len: usize,
}

/// Stream callback type
pub type XmtpStreamCallback = extern "C" fn(
    data: *mut c_void,
    data_len: usize,
    user_data: *mut c_void,
);

/// Stream error callback type
pub type XmtpStreamErrorCallback = extern "C" fn(
    error: *const XmtpFfiError,
    user_data: *mut c_void,
);
