package ffi

import (
	"unsafe"
)

// XmtpEnv represents the XMTP environment
type XmtpEnv int32

const (
	XmtpEnvLocal         XmtpEnv = 0
	XmtpEnvDev           XmtpEnv = 1
	XmtpEnvProduction    XmtpEnv = 2
	XmtpEnvTestnetStaging XmtpEnv = 3
	XmtpEnvTestnetDev    XmtpEnv = 4
	XmtpEnvTestnet       XmtpEnv = 5
	XmtpEnvMainnet       XmtpEnv = 6
)

// XmtpConversationType represents the type of conversation
type XmtpConversationType int32

const (
	XmtpConversationTypeDm    XmtpConversationType = 0
	XmtpConversationTypeGroup XmtpConversationType = 1
)

// XmtpConsentState represents the consent state
type XmtpConsentState int32

const (
	XmtpConsentStateUnknown XmtpConsentState = 0
	XmtpConsentStateAllowed XmtpConsentState = 1
	XmtpConsentStateDenied  XmtpConsentState = 2
)

// XmtpContentType represents the content type
type XmtpContentType int32

const (
	XmtpContentTypeText               XmtpContentType = 0
	XmtpContentTypeMarkdown           XmtpContentType = 1
	XmtpContentTypeReply              XmtpContentType = 2
	XmtpContentTypeReaction           XmtpContentType = 3
	XmtpContentTypeAttachment         XmtpContentType = 4
	XmtpContentTypeRemoteAttachment   XmtpContentType = 5
	XmtpContentTypeMultiRemoteAttachment XmtpContentType = 6
	XmtpContentTypeTransactionReference XmtpContentType = 7
	XmtpContentTypeGroupUpdated       XmtpContentType = 8
	XmtpContentTypeReadReceipt        XmtpContentType = 9
	XmtpContentTypeLeaveRequest       XmtpContentType = 10
	XmtpContentTypeWalletSendCalls    XmtpContentType = 11
	XmtpContentTypeActions            XmtpContentType = 12
	XmtpContentTypeIntent             XmtpContentType = 13
	XmtpContentTypeDeletedMessage     XmtpContentType = 14
	XmtpContentTypeCustom             XmtpContentType = 15
)

// XmtpDeliveryStatus represents the delivery status
type XmtpDeliveryStatus int32

const (
	XmtpDeliveryStatusUnpublished XmtpDeliveryStatus = 0
	XmtpDeliveryStatusPublished   XmtpDeliveryStatus = 1
	XmtpDeliveryStatusFailed      XmtpDeliveryStatus = 2
)

// XmtpIdentifierKind represents the identifier type
type XmtpIdentifierKind int32

const (
	XmtpIdentifierKindEthereum XmtpIdentifierKind = 0
	XmtpIdentifierKindPasskey  XmtpIdentifierKind = 1
)

// XmtpResult represents an FFI result
type XmtpResult struct {
	Error uintptr // *XmtpFfiError
}

// IsOk returns true if the result has no error
func (r XmtpResult) IsOk() bool {
	return r.Error == 0
}

// XmtpFfiError represents an FFI error
type XmtpFfiError struct {
	Code    int32
	Message uintptr
}

// XmtpStringResult represents a string result
type XmtpStringResult struct {
	Value uintptr
	Error uintptr
}

// IsOk returns true if the result has no error
func (r XmtpStringResult) IsOk() bool {
	return r.Error == 0
}

// XmtpBytesResult represents a bytes result
type XmtpBytesResult struct {
	Data  uintptr
	Len   uintptr
	Error uintptr
}

// XmtpBoolResult represents a bool result
type XmtpBoolResult struct {
	Value bool
	Error uintptr
}

// XmtpIntResult represents an int result
type XmtpIntResult struct {
	Value int32
	Error uintptr
}

// XmtpIdentifier represents an XMTP identifier
type XmtpIdentifier struct {
	Kind       XmtpIdentifierKind
	Identifier uintptr
}

// XmtpClientOptions represents client options
type XmtpClientOptions struct {
	Env                 XmtpEnv
	DbPath              uintptr
	DbEncryptionKey     uintptr
	DbEncryptionKeyLen  uintptr
	AppVersion          uintptr
	DisableAutoRegister bool
	StructuredLogging   bool
	LogLevel            int32
}

// NewXmtpClientOptions creates default client options
func NewXmtpClientOptions() XmtpClientOptions {
	return XmtpClientOptions{
		Env:                 XmtpEnvDev,
		DbPath:              0,
		DbEncryptionKey:     0,
		DbEncryptionKeyLen:  0,
		AppVersion:          0,
		DisableAutoRegister: false,
		StructuredLogging:   false,
		LogLevel:            2, // Warn
	}
}

// XmtpListMessagesOptions represents options for listing messages
type XmtpListMessagesOptions struct {
	Limit     uintptr
	BeforeNs  uint64
	AfterNs   uint64
	Ascending bool
}

// NewXmtpListMessagesOptions creates default list messages options
func NewXmtpListMessagesOptions() XmtpListMessagesOptions {
	return XmtpListMessagesOptions{
		Limit:     100,
		BeforeNs:  0,
		AfterNs:   0,
		Ascending: false,
	}
}

// XmtpCreateGroupOptions represents options for creating a group
type XmtpCreateGroupOptions struct {
	Name        uintptr
	ImageUrl    uintptr
	Description uintptr
	Permissions int32
}

// XmtpGroupMember represents a group member
type XmtpGroupMember struct {
	InboxId        uintptr
	PermissionLevel int32
}

// XmtpMessageFfi represents a message in FFI
type XmtpMessageFfi struct {
	Id             uintptr
	SenderInboxId  uintptr
	ConversationId uintptr
	SentAtNs       uint64
	ExpiresAtNs    uint64
	ContentType    XmtpContentType
	DeliveryStatus XmtpDeliveryStatus
	Fallback       uintptr
	ContentData    uintptr
	ContentLen     uintptr
}

// XmtpReactionContent represents reaction content
type XmtpReactionContent struct {
	ReferenceMessageId uintptr
	Action             int32
	Schema             int32
	Content            uintptr
}

// XmtpReplyContent represents reply content
type XmtpReplyContent struct {
	ReferenceMessageId uintptr
	ContentType        XmtpContentType
	Content            uintptr
}

// XmtpAttachmentContent represents attachment content
type XmtpAttachmentContent struct {
	Filename uintptr
	MimeType uintptr
	Data     uintptr
	DataLen  uintptr
}

// XmtpSignerCallback is the type for signer callbacks
type XmtpSignerCallback func(message uintptr, messageLen uintptr, userData uintptr, outSignature uintptr, outSignatureLen *uintptr) int32

// XmtpStreamCallback is the type for stream callbacks
type XmtpStreamCallback func(data uintptr, dataLen uintptr, userData uintptr)

// XmtpStreamErrorCallback is the type for stream error callbacks
type XmtpStreamErrorCallback func(error uintptr, userData uintptr)

// Helper functions for creating FFI types

// NewIdentifier creates a new XmtpIdentifier
func NewIdentifier(kind XmtpIdentifierKind, identifier string) XmtpIdentifier {
	return XmtpIdentifier{
		Kind:       kind,
		Identifier: CString(identifier),
	}
}

// NewClientOptions creates XmtpClientOptions from Go values
func NewClientOptions(env XmtpEnv, dbPath string, appVersion string, disableAutoRegister bool, logLevel int32) XmtpClientOptions {
	opts := NewXmtpClientOptions()
	opts.Env = env
	opts.DisableAutoRegister = disableAutoRegister
	opts.LogLevel = logLevel
	
	if dbPath != "" {
		opts.DbPath = CString(dbPath)
	}
	if appVersion != "" {
		opts.AppVersion = CString(appVersion)
	}
	
	return opts
}

// CStringArray creates a C array of strings from a Go slice
func CStringArray(strs []string) (uintptr, func()) {
	if len(strs) == 0 {
		return 0, func() {}
	}
	
	// Allocate memory for pointers
	ptrs := make([]uintptr, len(strs))
	for i, s := range strs {
		ptrs[i] = CString(s)
	}
	
	// Return pointer to first element and cleanup function
	return uintptr(unsafe.Pointer(&ptrs[0])), func() {
		// Cleanup is handled by Go's GC for now
	}
}
