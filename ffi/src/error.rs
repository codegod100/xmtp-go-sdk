//! Error handling for FFI

use std::ffi::CString;
use thiserror::Error;

/// FFI Error type
#[derive(Error, Debug)]
pub enum XmtpFfiError {
    #[error("Client error: {0}")]
    Client(String),
    
    #[error("Conversation error: {0}")]
    Conversation(String),
    
    #[error("Message error: {0}")]
    Message(String),
    
    #[error("Storage error: {0}")]
    Storage(String),
    
    #[error("Signature error: {0}")]
    Signature(String),
    
    #[error("Network error: {0}")]
    Network(String),
    
    #[error("Invalid argument: {0}")]
    InvalidArgument(String),
    
    #[error("Not found: {0}")]
    NotFound(String),
    
    #[error("Not registered")]
    NotRegistered,
    
    #[error("Signer unavailable")]
    SignerUnavailable,
    
    #[error("Library not initialized")]
    NotInitialized,
    
    #[error("Generic error: {0}")]
    Generic(String),
}

impl XmtpFfiError {
    pub fn to_ffi(self) -> crate::types::XmtpFfiError {
        let code = match &self {
            XmtpFfiError::Client(_) => 1,
            XmtpFfiError::Conversation(_) => 2,
            XmtpFfiError::Message(_) => 3,
            XmtpFfiError::Storage(_) => 4,
            XmtpFfiError::Signature(_) => 5,
            XmtpFfiError::Network(_) => 6,
            XmtpFfiError::InvalidArgument(_) => 7,
            XmtpFfiError::NotFound(_) => 8,
            XmtpFfiError::NotRegistered => 9,
            XmtpFfiError::SignerUnavailable => 10,
            XmtpFfiError::NotInitialized => 11,
            XmtpFfiError::Generic(_) => -1,
        };
        
        let msg = self.to_string();
        let c_msg = CString::new(msg).unwrap_or_default();
        
        crate::types::XmtpFfiError {
            code,
            message: c_msg.into_raw(),
        }
    }
}

impl From<String> for XmtpFfiError {
    fn from(s: String) -> Self {
        XmtpFfiError::Generic(s)
    }
}

impl From<&str> for XmtpFfiError {
    fn from(s: &str) -> Self {
        XmtpFfiError::Generic(s.to_string())
    }
}
