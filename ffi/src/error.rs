//! Error handling for FFI

use std::ffi::CString;
use thiserror::Error;

/// FFI Error enum (internal)
#[derive(Error, Debug)]
pub enum XmtpError {
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

impl XmtpError {
    pub fn to_ffi(self) -> crate::types::XmtpFfiError {
        let code = match &self {
            XmtpError::Client(_) => 1,
            XmtpError::Conversation(_) => 2,
            XmtpError::Message(_) => 3,
            XmtpError::Storage(_) => 4,
            XmtpError::Signature(_) => 5,
            XmtpError::Network(_) => 6,
            XmtpError::InvalidArgument(_) => 7,
            XmtpError::NotFound(_) => 8,
            XmtpError::NotRegistered => 9,
            XmtpError::SignerUnavailable => 10,
            XmtpError::NotInitialized => 11,
            XmtpError::Generic(_) => -1,
        };
        
        let msg = self.to_string();
        let c_msg = CString::new(msg).unwrap_or_default();
        
        crate::types::XmtpFfiError {
            code,
            message: c_msg.into_raw(),
        }
    }
}

impl From<String> for XmtpError {
    fn from(s: String) -> Self {
        XmtpError::Generic(s)
    }
}

impl From<&str> for XmtpError {
    fn from(s: &str) -> Self {
        XmtpError::Generic(s.to_string())
    }
}
