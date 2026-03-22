// Package xmtp provides a Go SDK for XMTP messaging
package xmtp

import (
	"errors"
	"time"
)

// Environment types
type Env int

const (
	EnvLocal Env = iota
	EnvDev
	EnvProduction
	EnvTestnetStaging
	EnvTestnetDev
	EnvTestnet
	EnvMainnet
)

// ConsentState represents the consent state of a conversation
type ConsentState int

const (
	ConsentUnknown ConsentState = iota
	ConsentAllowed
	ConsentDenied
)

// ConversationType represents the type of conversation
type ConversationType int

const (
	ConversationTypeDM ConversationType = iota
	ConversationTypeGroup
)

// DeliveryStatus represents the delivery status of a message
type DeliveryStatus int

const (
	DeliveryStatusUnpublished DeliveryStatus = iota
	DeliveryStatusPublished
	DeliveryStatusFailed
)

// ContentType represents the type of message content
type ContentType int

const (
	ContentTypeText ContentType = iota
	ContentTypeMarkdown
	ContentTypeReply
	ContentTypeReaction
	ContentTypeAttachment
	ContentTypeRemoteAttachment
	ContentTypeMultiRemoteAttachment
	ContentTypeTransactionReference
	ContentTypeGroupUpdated
	ContentTypeReadReceipt
	ContentTypeLeaveRequest
	ContentTypeWalletSendCalls
	ContentTypeActions
	ContentTypeIntent
	ContentTypeDeletedMessage
	ContentTypeCustom
)

// Common errors
var (
	ErrClientNotInitialized = errors.New("client not initialized")
	ErrSignerUnavailable    = errors.New("signer unavailable")
	ErrConversationNotFound = errors.New("conversation not found")
	ErrMessageNotFound      = errors.New("message not found")
	ErrNotRegistered        = errors.New("client not registered")
	ErrLibraryNotLoaded     = errors.New("libxmtp_ffi library not loaded")
)

// Identifier represents an XMTP identifier (e.g., Ethereum address)
type Identifier struct {
	Kind       int    // 0 = Ethereum
	Identifier string // The address or identifier string
}

// Signer is the interface for signing messages
type Signer interface {
	// GetIdentifier returns the identifier for this signer
	GetIdentifier() (Identifier, error)

	// SignMessage signs a message and returns the signature
	SignMessage(message []byte) ([]byte, error)
}

// SCWSigner is a signer for Smart Contract Wallets
type SCWSigner interface {
	Signer

	// GetChainId returns the chain ID
	GetChainId() (uint64, error)

	// GetBlockNumber returns the current block number (optional)
	GetBlockNumber() (uint64, error)
}

// ClientOption is a function that configures the client
type ClientOption func(*ClientOptions)

// ClientOptions holds configuration for the XMTP client
type ClientOptions struct {
	Env                Env
	DbPath             string
	DbEncryptionKey    []byte
	AppVersion         string
	DisableAutoRegister bool
	StructuredLogging  bool
	LogLevel           int // 0=off, 1=error, 2=warn, 3=info, 4=debug, 5=trace
}

// WithEnv sets the XMTP environment
func WithEnv(env Env) ClientOption {
	return func(o *ClientOptions) {
		o.Env = env
	}
}

// WithDbPath sets the database path
func WithDbPath(path string) ClientOption {
	return func(o *ClientOptions) {
		o.DbPath = path
	}
}

// WithDbEncryptionKey sets the database encryption key (32 bytes)
func WithDbEncryptionKey(key []byte) ClientOption {
	return func(o *ClientOptions) {
		o.DbEncryptionKey = key
	}
}

// WithAppVersion sets the app version
func WithAppVersion(version string) ClientOption {
	return func(o *ClientOptions) {
		o.AppVersion = version
	}
}

// WithDisableAutoRegister disables automatic registration
func WithDisableAutoRegister() ClientOption {
	return func(o *ClientOptions) {
		o.DisableAutoRegister = true
	}
}

// WithStructuredLogging enables structured JSON logging
func WithStructuredLogging() ClientOption {
	return func(o *ClientOptions) {
		o.StructuredLogging = true
	}
}

// WithLogLevel sets the logging level
func WithLogLevel(level int) ClientOption {
	return func(o *ClientOptions) {
		o.LogLevel = level
	}
}

// Message represents a decoded XMTP message
type Message struct {
	ID              string
	SenderInboxID   string
	ConversationID  string
	SentAt          time.Time
	ExpiresAt       *time.Time
	ContentType     ContentType
	DeliveryStatus  DeliveryStatus
	Fallback        string
	Content         any // The decoded content (string, Reaction, Reply, etc.)
}

// Reaction represents a reaction to a message
type Reaction struct {
	ReferenceID string
	Action      int // 0 = add, 1 = remove
	Schema      int // 0 = unicode, 1 = custom
	Content     string
}

// Reply represents a reply to a message
type Reply struct {
	ReferenceID string
	ContentType ContentType
	Content     any
}

// Attachment represents a file attachment
type Attachment struct {
	Filename string
	MimeType string
	Data     []byte
}

// GroupMember represents a member of a group
type GroupMember struct {
	InboxID            string
	AccountIdentifiers []Identifier
	PermissionLevel    int // 0=member, 1=admin, 2=super_admin
}

// CreateGroupOptions holds options for creating a new group
type CreateGroupOptions struct {
	Name        string
	ImageURL    string
	Description string
	Permissions int // Permission policy set
}

// ListMessagesOptions holds options for listing messages
type ListMessagesOptions struct {
	Limit     int
	Before    *time.Time
	After     *time.Time
	Ascending bool
}

// StreamOptions holds options for streaming
type StreamOptions struct {
	OnEnd   func()
	OnError func(error)
}
