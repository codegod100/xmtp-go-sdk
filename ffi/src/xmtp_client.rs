//! XMTP Client integration with libxmtp

use std::ffi::{c_char, c_void, CString};
use std::ptr;
use std::sync::Arc;
use parking_lot::RwLock;

use xmtp_mls::client::Client as MlsClient;
use xmtp_mls::MlsContext;
use xmtp_db::{NativeDb, EncryptedMessageStore};
use xmtp_id::associations::Identifier;
use xmtp_id::InboxId;
use xmtp_api_grpc::ApiClientWrapper;
use xmtp_api::strategies;

use crate::types::*;
use crate::error::XmtpError;
use crate::signer::FfiSigner;
use crate::{block_on, get_runtime};

/// The actual XMTP client type
pub type XmtpClient = MlsClient<MlsContext>;

/// Internal client wrapper that holds the actual client
pub struct XmtpClientInner {
    client: Arc<XmtpClient>,
    inbox_id: String,
}

impl XmtpClientInner {
    /// Create a new client with a signer callback
    pub fn new(
        signer: FfiSigner,
        identifier: Identifier,
        opts: XmtpClientOptions,
    ) -> Result<Self, XmtpError> {
        block_on(async {
            Self::new_async(signer, identifier, opts).await
        })
    }

    async fn new_async(
        signer: FfiSigner,
        identifier: Identifier,
        opts: XmtpClientOptions,
    ) -> Result<Self, XmtpError> {
        // Build database
        let db = if let Some(path) = unsafe { opts.db_path.as_ref() } {
            let path = std::ffi::CStr::from_ptr(path)
                .to_str()
                .map_err(|e| XmtpError::Generic(e.to_string()))?;
            NativeDb::builder().persistent(path)
        } else {
            NativeDb::builder().ephemeral()
        };

        let db = if let (Some(key), Some(len)) = (
            unsafe { opts.db_encryption_key.as_ref() },
            opts.db_encryption_key_len,
        ) {
            if len != 32 {
                return Err(XmtpError::Generic("Encryption key must be 32 bytes".into()));
            }
            let key_bytes: [u8; 32] = unsafe { std::slice::from_raw_parts(key, 32) }
                .try_into()
                .map_err(|_| XmtpError::Generic("Invalid encryption key".into()))?;
            db.key(key_bytes).build()
        } else {
            db.build_unencrypted()
        }.map_err(|e| XmtpError::Storage(e.to_string()))?;

        let store = EncryptedMessageStore::new(db)
            .map_err(|e| XmtpError::Storage(e.to_string()))?;

        // Get environment config
        let env = match opts.env {
            XmtpEnv::Local => xmtp_configuration::Configuration::local(),
            XmtpEnv::Dev => xmtp_configuration::Configuration::dev(),
            XmtpEnv::Production => xmtp_configuration::Configuration::production(),
            XmtpEnv::TestnetStaging => xmtp_configuration::Configuration::testnet_staging(),
            XmtpEnv::TestnetDev => xmtp_configuration::Configuration::testnet_dev(),
            XmtpEnv::Testnet => xmtp_configuration::Configuration::testnet(),
            XmtpEnv::Mainnet => xmtp_configuration::Configuration::mainnet(),
        };

        let host = env.api_grpc_host();
        let app_version = unsafe {
            opts.app_version.as_ref()
                .map(|v| std::ffi::CStr::from_ptr(v).to_string_lossy().into_owned())
                .unwrap_or_default()
        };

        // Create API client
        let api_client = xmtp_api_grpc::grpc_client::GrpcClient::new(&host, &app_version)
            .map_err(|e| XmtpError::Network(e.to_string()))?;

        // Create identity strategy
        let inbox_id = generate_inbox_id(&identifier, 0); // TODO: proper nonce handling
        let identity_strategy = xmtp_mls::identity::IdentityStrategy::new(
            inbox_id.clone(),
            identifier.clone(),
            0,
            None,
        );

        // Build client
        let client = xmtp_mls::Client::builder(identity_strategy)
            .api_client(api_client)
            .store(store)
            .default_mls_store()
            .map_err(|e| XmtpError::Client(e.to_string()))?
            .build()
            .await
            .map_err(|e| XmtpError::Client(e.to_string()))?;

        Ok(XmtpClientInner {
            client: Arc::new(client),
            inbox_id,
        })
    }

    /// Get inbox ID
    pub fn inbox_id(&self) -> &str {
        &self.inbox_id
    }

    /// Check if registered
    pub fn is_registered(&self) -> bool {
        block_on(async {
            self.client.inbox_state(false).await.is_ok()
        })
    }

    /// Get installation ID
    pub fn installation_id(&self) -> Vec<u8> {
        self.client.installation_public_key().to_vec()
    }

    /// Register the client
    pub fn register(&self, signer: &FfiSigner) -> Result<(), XmtpError> {
        // TODO: Implement signature request with signer callback
        Err(XmtpError::Generic("Registration not yet implemented".into()))
    }

    /// Get conversations
    pub fn conversations(&self) -> Arc<XmtpConversationsInner> {
        Arc::new(XmtpConversationsInner {
            client: self.client.clone(),
        })
    }
}

fn generate_inbox_id(identifier: &Identifier, nonce: u64) -> String {
    // Use libxmtp's inbox ID generation
    xmtp_id::generate_inbox_id(identifier, nonce)
        .to_string()
}

/// Conversations manager wrapper
pub struct XmtpConversationsInner {
    client: Arc<XmtpClient>,
}

impl XmtpConversationsInner {
    /// List conversations
    pub fn list(&self, conv_type: XmtpConversationType) -> Result<Vec<XmtpConversationInner>, XmtpError> {
        block_on(async {
            let query = xmtp_db::group::GroupQueryArgs::default();
            let groups = self.client.find_groups(query)
                .await
                .map_err(|e| XmtpError::Conversation(e.to_string()))?;

            Ok(groups
                .into_iter()
                .filter(|g| {
                    match conv_type {
                        XmtpConversationType::Dm => g.conversation_type() == xmtp_db::group::ConversationType::Dm,
                        XmtpConversationType::Group => g.conversation_type() == xmtp_db::group::ConversationType::Group,
                    }
                })
                .map(|g| XmtpConversationInner { group: Arc::new(g) })
                .collect())
        })
    }

    /// Create a DM
    pub fn create_dm(&self, target_inbox_id: &str) -> Result<XmtpConversationInner, XmtpError> {
        block_on(async {
            let group = self.client.dm_group_from_target_inbox(target_inbox_id)
                .map_err(|e| XmtpError::Conversation(e.to_string()))?;
            Ok(XmtpConversationInner { group: Arc::new(group) })
        })
    }

    /// Create a group
    pub fn create_group(&self, name: Option<&str>, description: Option<&str>) -> Result<XmtpConversationInner, XmtpError> {
        block_on(async {
            let metadata_opts = xmtp_mls::mls_common::group::GroupMetadataOptions {
                name: name.map(String::from),
                description: description.map(String::from),
                ..Default::default()
            };

            let group = self.client.create_group(None, Some(metadata_opts))
                .map_err(|e| XmtpError::Conversation(e.to_string()))?;

            Ok(XmtpConversationInner { group: Arc::new(group) })
        })
    }
}

/// Conversation wrapper
pub struct XmtpConversationInner {
    group: Arc<xmtp_mls::groups::MlsGroup<MlsContext>>,
}

impl XmtpConversationInner {
    /// Get conversation ID
    pub fn id(&self) -> Vec<u8> {
        self.group.group_id().to_vec()
    }

    /// Check if active
    pub fn is_active(&self) -> bool {
        self.group.is_active()
    }

    /// Get created timestamp
    pub fn created_at_ns(&self) -> u64 {
        self.group.created_at_ns() as u64
    }

    /// Send a text message
    pub fn send_text(&self, text: &str) -> Result<Vec<u8>, XmtpError> {
        block_on(async {
            let message_id = self.group.send_text(text)
                .await
                .map_err(|e| XmtpError::Message(e.to_string()))?;
            Ok(message_id.to_vec())
        })
    }

    /// List messages
    pub fn list_messages(&self, limit: usize, before_ns: u64, after_ns: u64) -> Result<Vec<XmtpMessageInner>, XmtpError> {
        block_on(async {
            let args = xmtp_db::group_message::MsgQueryArgs {
                limit: Some(limit as i64),
                kind: Some(xmtp_db::group_message::GroupMessageKind::Application),
                ..Default::default()
            };

            let messages = self.group.find_messages(&args)
                .await
                .map_err(|e| XmtpError::Message(e.to_string()))?;

            Ok(messages.into_iter().map(XmtpMessageInner::from).collect())
        })
    }
}

/// Message wrapper
pub struct XmtpMessageInner {
    pub id: Vec<u8>,
    pub sender_inbox_id: String,
    pub sent_at_ns: u64,
    pub content: Vec<u8>,
    pub content_type: XmtpContentType,
}

impl From<xmtp_db::group_message::StoredGroupMessage> for XmtpMessageInner {
    fn from(msg: xmtp_db::group_message::StoredGroupMessage) -> Self {
        let content_type = match msg.content_type {
            xmtp_db::group_message::ContentType::Text => XmtpContentType::Text,
            xmtp_db::group_message::ContentType::Attachment => XmtpContentType::Attachment,
            xmtp_db::group_message::ContentType::Reaction => XmtpContentType::Reaction,
            xmtp_db::group_message::ContentType::Reply => XmtpContentType::Reply,
            _ => XmtpContentType::Custom,
        };

        XmtpMessageInner {
            id: msg.id,
            sender_inbox_id: msg.sender_inbox_id,
            sent_at_ns: msg.sent_at_ns as u64,
            content: msg.decrypted_message_bytes,
            content_type,
        }
    }
}
