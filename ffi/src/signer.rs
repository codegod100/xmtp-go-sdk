//! Signer callback handling

use std::ffi::c_void;
use crate::types::*;

/// FFI Signer wrapper
pub struct FfiSigner {
    callback: Option<XmtpSignerCallback>,
    user_data: *mut c_void,
}

impl FfiSigner {
    pub fn new(callback: XmtpSignerCallback, user_data: *mut c_void) -> Self {
        FfiSigner { 
            callback: Some(callback), 
            user_data 
        }
    }
    
    /// Create an empty signer (for building without signing)
    pub fn empty() -> Self {
        FfiSigner {
            callback: None,
            user_data: std::ptr::null_mut(),
        }
    }
    
    pub fn sign(&self, message: &[u8]) -> Result<Vec<u8>, String> {
        let callback = match self.callback {
            Some(cb) => cb,
            None => return Err("No signer callback configured".into()),
        };
        
        let mut signature = vec![0u8; 65]; // Ethereum signature size
        let mut signature_len = signature.len();
        
        let result = (callback)(
            message.as_ptr(),
            message.len(),
            self.user_data,
            signature.as_mut_ptr(),
            &mut signature_len,
        );
        
        if result != 0 {
            return Err(format!("Signer callback returned error: {}", result));
        }
        
        signature.truncate(signature_len);
        Ok(signature)
    }
}

// Safety: The callback must remain valid for the lifetime of this struct
unsafe impl Send for FfiSigner {}
unsafe impl Sync for FfiSigner {}
