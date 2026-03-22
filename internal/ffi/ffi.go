package ffi

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/ebitengine/purego"
)

var (
	ErrLibraryNotLoaded = errors.New("libxmtp_ffi library not loaded")
)

// Library paths to search
var libraryPaths = []string{
	"libxmtp_ffi.so",
	"libxmtp_ffi.dylib",
	"xmtp_ffi.dll",
	"/usr/local/lib/libxmtp_ffi.so",
	"/usr/local/lib/libxmtp_ffi.dylib",
	"/usr/lib/libxmtp_ffi.so",
	"./libxmtp_ffi.so",
	"./libxmtp_ffi.dylib",
}

// Library holds the loaded dynamic library and its symbols
type Library struct {
	handle uintptr

	// Initialization
	XmtpInit    func() int32
	XmtpVersion func() uintptr
	XmtpStringFree func(s uintptr)
	XmtpBytesFree  func(data uintptr, len uintptr)

	// Client functions
	XmtpClientCreate              func(signerCB uintptr, signerData uintptr, identifier XmtpIdentifier, opts uintptr, outClient *uintptr) XmtpResult
	XmtpClientBuild               func(identifier XmtpIdentifier, opts uintptr, outClient *uintptr) XmtpResult
	XmtpClientFree                func(client uintptr)
	XmtpClientInboxId             func(client uintptr, outResult *XmtpStringResult)
	XmtpClientInstallationId      func(client uintptr, outData uintptr, outLen *uintptr) XmtpResult
	XmtpClientIsRegistered        func(client uintptr) bool
	XmtpClientConversations       func(client uintptr, outConversations *uintptr) XmtpResult
	XmtpClientRegister            func(client uintptr, signerCB uintptr, signerData uintptr) XmtpResult
	XmtpClientCanMessage          func(client uintptr, identifiers uintptr, identifiersLen uintptr, outResults uintptr, outLen *uintptr) XmtpResult
	XmtpClientGetInboxIdByIdentifier func(client uintptr, identifier XmtpIdentifier, outResult *XmtpStringResult)
	XmtpClientLibxmtpVersion      func() uintptr

	// Conversations functions
	XmtpConversationsFree            func(conversations uintptr)
	XmtpConversationsList            func(conversations uintptr, convType int32, outConversations uintptr, outLen *uintptr) XmtpResult
	XmtpConversationsGetById         func(conversations uintptr, id uintptr, outConversation *uintptr) XmtpResult
	XmtpConversationsGetDmByInboxId  func(conversations uintptr, inboxId uintptr, outConversation *uintptr) XmtpResult
	XmtpConversationsCreateGroup     func(conversations uintptr, inboxIds uintptr, inboxIdsLen uintptr, name uintptr, imageUrl uintptr, description uintptr, outGroup *uintptr) XmtpResult
	XmtpConversationsCreateDm        func(conversations uintptr, inboxId uintptr, outDm *uintptr) XmtpResult
	XmtpConversationsSync            func(conversations uintptr) XmtpResult
	XmtpConversationsStream          func(conversations uintptr, convType int32, callback uintptr, errorCallback uintptr, userData uintptr, outStream *uintptr) XmtpResult
	XmtpStreamEnd                    func(stream uintptr)

	// Conversation functions
	XmtpConversationFree            func(conversation uintptr)
	XmtpConversationId              func(conversation uintptr, outResult *XmtpStringResult)
	XmtpConversationIsActive        func(conversation uintptr) bool
	XmtpConversationCreatedAtNs     func(conversation uintptr) uint64
	XmtpConversationConsentState    func(conversation uintptr) int32
	XmtpConversationUpdateConsent   func(conversation uintptr, state int32) XmtpResult
	XmtpConversationSync            func(conversation uintptr) XmtpResult
	XmtpConversationSendText        func(conversation uintptr, text uintptr, optimistic bool, outResult *XmtpStringResult)
	XmtpConversationSendMarkdown    func(conversation uintptr, markdown uintptr, optimistic bool, outResult *XmtpStringResult)
	XmtpConversationSendReaction    func(conversation uintptr, refId uintptr, action int32, schema int32, content uintptr, optimistic bool, outResult *XmtpStringResult)
	XmtpConversationListMessages    func(conversation uintptr, opts uintptr, outMessages uintptr, outLen *uintptr) XmtpResult
	XmtpConversationGetMessageById  func(conversation uintptr, messageId uintptr, outMessage *uintptr) XmtpResult
	XmtpConversationStreamMessages  func(conversation uintptr, callback uintptr, errorCallback uintptr, userData uintptr, outStream *uintptr) XmtpResult

	// Group functions
	XmtpGroupName              func(conversation uintptr) uintptr
	XmtpGroupUpdateName        func(conversation uintptr, name uintptr) XmtpResult
	XmtpGroupImageUrl          func(conversation uintptr) uintptr
	XmtpGroupUpdateImageUrl    func(conversation uintptr, url uintptr) XmtpResult
	XmtpGroupDescription       func(conversation uintptr) uintptr
	XmtpGroupUpdateDescription func(conversation uintptr, description uintptr) XmtpResult
	XmtpGroupListMembers       func(conversation uintptr, outMembers uintptr, outLen *uintptr) XmtpResult
	XmtpGroupAddMembers        func(conversation uintptr, inboxIds uintptr, inboxIdsLen uintptr) XmtpResult
	XmtpGroupRemoveMembers     func(conversation uintptr, inboxIds uintptr, inboxIdsLen uintptr) XmtpResult
	XmtpGroupIsAdmin           func(conversation uintptr, inboxId uintptr) bool
	XmtpGroupIsSuperAdmin      func(conversation uintptr, inboxId uintptr) bool
	XmtpGroupAddAdmin          func(conversation uintptr, inboxId uintptr) XmtpResult
	XmtpGroupRemoveAdmin       func(conversation uintptr, inboxId uintptr) XmtpResult
	XmtpGroupLeave             func(conversation uintptr) XmtpResult

	// DM functions
	XmtpDmPeerInboxId func(conversation uintptr) uintptr

	// Message functions
	XmtpMessageFree              func(message uintptr)
	XmtpMessageId                func(message uintptr) uintptr
	XmtpMessageSenderInboxId     func(message uintptr) uintptr
	XmtpMessageSentAtNs          func(message uintptr) uint64
	XmtpMessageExpiresAtNs       func(message uintptr) uint64
	XmtpMessageConversationId    func(message uintptr) uintptr
	XmtpMessageContentType       func(message uintptr) int32
	XmtpMessageDeliveryStatus    func(message uintptr) int32
	XmtpMessageFallback          func(message uintptr) uintptr
	XmtpMessageContentText       func(message uintptr) uintptr
	XmtpMessageContentMarkdown   func(message uintptr) uintptr
	XmtpMessageContentReaction   func(message uintptr, outReaction *XmtpReactionContent)
	XmtpMessageContentReply      func(message uintptr, outReply *XmtpReplyContent)
	XmtpMessageContentAttachment func(message uintptr, outAttachment *XmtpAttachmentContent)
	XmtpMessageContentBytes      func(message uintptr, outData uintptr, outLen *uintptr) XmtpResult

	// Content cleanup
	XmtpReactionContentFree   func(reaction uintptr)
	XmtpReplyContentFree      func(reply uintptr)
	XmtpAttachmentContentFree func(attachment uintptr)
}

var (
	library     *Library
	libraryOnce sync.Once
	libraryErr  error
)

// LoadLibrary loads the libxmtp_ffi shared library
func LoadLibrary() error {
	libraryOnce.Do(func() {
		library, libraryErr = loadLibrary()
	})
	return libraryErr
}

// LoadLibraryFromPath loads the library from a specific path
func LoadLibraryFromPath(path string) error {
	libraryOnce.Do(func() {
		library, libraryErr = loadLibraryFromPath(path)
	})
	return libraryErr
}

func loadLibrary() (*Library, error) {
	// Check for environment variable override
	if path := os.Getenv("XMTP_FFI_PATH"); path != "" {
		return loadLibraryFromPath(path)
	}

	// Try each library path
	for _, path := range libraryPaths {
		lib, err := loadLibraryFromPath(path)
		if err == nil {
			return lib, nil
		}
	}

	return nil, fmt.Errorf("could not find libxmtp_ffi, tried: %v", libraryPaths)
}

func loadLibraryFromPath(path string) (*Library, error) {
	handle, err := purego.Dlopen(path, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		return nil, fmt.Errorf("failed to load %s: %w", path, err)
	}

	lib := &Library{handle: handle}

	// Register all FFI symbols
	// Use panic recovery since RegisterLibFunc panics on error
	defer func() {
		if r := recover(); r != nil {
			libraryErr = fmt.Errorf("failed to register symbol: %v", r)
		}
	}()

	// Initialization
	purego.RegisterLibFunc(&lib.XmtpInit, handle, "xmtp_init")
	purego.RegisterLibFunc(&lib.XmtpVersion, handle, "xmtp_version")
	purego.RegisterLibFunc(&lib.XmtpStringFree, handle, "xmtp_string_free")
	purego.RegisterLibFunc(&lib.XmtpBytesFree, handle, "xmtp_bytes_free")

	// Client
	purego.RegisterLibFunc(&lib.XmtpClientCreate, handle, "xmtp_client_create")
	purego.RegisterLibFunc(&lib.XmtpClientBuild, handle, "xmtp_client_build")
	purego.RegisterLibFunc(&lib.XmtpClientFree, handle, "xmtp_client_free")
	purego.RegisterLibFunc(&lib.XmtpClientInboxId, handle, "xmtp_client_inbox_id")
	purego.RegisterLibFunc(&lib.XmtpClientInstallationId, handle, "xmtp_client_installation_id")
	purego.RegisterLibFunc(&lib.XmtpClientIsRegistered, handle, "xmtp_client_is_registered")
	purego.RegisterLibFunc(&lib.XmtpClientConversations, handle, "xmtp_client_conversations")
	purego.RegisterLibFunc(&lib.XmtpClientRegister, handle, "xmtp_client_register")
	purego.RegisterLibFunc(&lib.XmtpClientCanMessage, handle, "xmtp_client_can_message")
	purego.RegisterLibFunc(&lib.XmtpClientGetInboxIdByIdentifier, handle, "xmtp_client_get_inbox_id_by_identifier")
	purego.RegisterLibFunc(&lib.XmtpClientLibxmtpVersion, handle, "xmtp_client_libxmtp_version")

	// Conversations
	purego.RegisterLibFunc(&lib.XmtpConversationsFree, handle, "xmtp_conversations_free")
	purego.RegisterLibFunc(&lib.XmtpConversationsList, handle, "xmtp_conversations_list")
	purego.RegisterLibFunc(&lib.XmtpConversationsGetById, handle, "xmtp_conversations_get_by_id")
	purego.RegisterLibFunc(&lib.XmtpConversationsGetDmByInboxId, handle, "xmtp_conversations_get_dm_by_inbox_id")
	purego.RegisterLibFunc(&lib.XmtpConversationsCreateGroup, handle, "xmtp_conversations_create_group")
	purego.RegisterLibFunc(&lib.XmtpConversationsCreateDm, handle, "xmtp_conversations_create_dm")
	purego.RegisterLibFunc(&lib.XmtpConversationsSync, handle, "xmtp_conversations_sync")
	purego.RegisterLibFunc(&lib.XmtpConversationsStream, handle, "xmtp_conversations_stream")
	purego.RegisterLibFunc(&lib.XmtpStreamEnd, handle, "xmtp_stream_end")

	// Conversation
	purego.RegisterLibFunc(&lib.XmtpConversationFree, handle, "xmtp_conversation_free")
	purego.RegisterLibFunc(&lib.XmtpConversationId, handle, "xmtp_conversation_id")
	purego.RegisterLibFunc(&lib.XmtpConversationIsActive, handle, "xmtp_conversation_is_active")
	purego.RegisterLibFunc(&lib.XmtpConversationCreatedAtNs, handle, "xmtp_conversation_created_at_ns")
	purego.RegisterLibFunc(&lib.XmtpConversationConsentState, handle, "xmtp_conversation_consent_state")
	purego.RegisterLibFunc(&lib.XmtpConversationUpdateConsent, handle, "xmtp_conversation_update_consent")
	purego.RegisterLibFunc(&lib.XmtpConversationSync, handle, "xmtp_conversation_sync")
	purego.RegisterLibFunc(&lib.XmtpConversationSendText, handle, "xmtp_conversation_send_text")
	purego.RegisterLibFunc(&lib.XmtpConversationSendMarkdown, handle, "xmtp_conversation_send_markdown")
	purego.RegisterLibFunc(&lib.XmtpConversationSendReaction, handle, "xmtp_conversation_send_reaction")
	purego.RegisterLibFunc(&lib.XmtpConversationListMessages, handle, "xmtp_conversation_list_messages")
	purego.RegisterLibFunc(&lib.XmtpConversationGetMessageById, handle, "xmtp_conversation_get_message_by_id")
	purego.RegisterLibFunc(&lib.XmtpConversationStreamMessages, handle, "xmtp_conversation_stream_messages")

	// Group
	purego.RegisterLibFunc(&lib.XmtpGroupName, handle, "xmtp_group_name")
	purego.RegisterLibFunc(&lib.XmtpGroupUpdateName, handle, "xmtp_group_update_name")
	purego.RegisterLibFunc(&lib.XmtpGroupImageUrl, handle, "xmtp_group_image_url")
	purego.RegisterLibFunc(&lib.XmtpGroupUpdateImageUrl, handle, "xmtp_group_update_image_url")
	purego.RegisterLibFunc(&lib.XmtpGroupDescription, handle, "xmtp_group_description")
	purego.RegisterLibFunc(&lib.XmtpGroupUpdateDescription, handle, "xmtp_group_update_description")
	purego.RegisterLibFunc(&lib.XmtpGroupListMembers, handle, "xmtp_group_list_members")
	purego.RegisterLibFunc(&lib.XmtpGroupAddMembers, handle, "xmtp_group_add_members")
	purego.RegisterLibFunc(&lib.XmtpGroupRemoveMembers, handle, "xmtp_group_remove_members")
	purego.RegisterLibFunc(&lib.XmtpGroupIsAdmin, handle, "xmtp_group_is_admin")
	purego.RegisterLibFunc(&lib.XmtpGroupIsSuperAdmin, handle, "xmtp_group_is_super_admin")
	purego.RegisterLibFunc(&lib.XmtpGroupAddAdmin, handle, "xmtp_group_add_admin")
	purego.RegisterLibFunc(&lib.XmtpGroupRemoveAdmin, handle, "xmtp_group_remove_admin")
	purego.RegisterLibFunc(&lib.XmtpGroupLeave, handle, "xmtp_group_leave")

	// DM
	purego.RegisterLibFunc(&lib.XmtpDmPeerInboxId, handle, "xmtp_dm_peer_inbox_id")

	// Message
	purego.RegisterLibFunc(&lib.XmtpMessageFree, handle, "xmtp_message_free")
	purego.RegisterLibFunc(&lib.XmtpMessageId, handle, "xmtp_message_id")
	purego.RegisterLibFunc(&lib.XmtpMessageSenderInboxId, handle, "xmtp_message_sender_inbox_id")
	purego.RegisterLibFunc(&lib.XmtpMessageSentAtNs, handle, "xmtp_message_sent_at_ns")
	purego.RegisterLibFunc(&lib.XmtpMessageExpiresAtNs, handle, "xmtp_message_expires_at_ns")
	purego.RegisterLibFunc(&lib.XmtpMessageConversationId, handle, "xmtp_message_conversation_id")
	purego.RegisterLibFunc(&lib.XmtpMessageContentType, handle, "xmtp_message_content_type")
	purego.RegisterLibFunc(&lib.XmtpMessageDeliveryStatus, handle, "xmtp_message_delivery_status")
	purego.RegisterLibFunc(&lib.XmtpMessageFallback, handle, "xmtp_message_fallback")
	purego.RegisterLibFunc(&lib.XmtpMessageContentText, handle, "xmtp_message_content_text")
	purego.RegisterLibFunc(&lib.XmtpMessageContentMarkdown, handle, "xmtp_message_content_markdown")
	purego.RegisterLibFunc(&lib.XmtpMessageContentReaction, handle, "xmtp_message_content_reaction")
	purego.RegisterLibFunc(&lib.XmtpMessageContentReply, handle, "xmtp_message_content_reply")
	purego.RegisterLibFunc(&lib.XmtpMessageContentAttachment, handle, "xmtp_message_content_attachment")
	purego.RegisterLibFunc(&lib.XmtpMessageContentBytes, handle, "xmtp_message_content_bytes")

	// Content cleanup
	purego.RegisterLibFunc(&lib.XmtpReactionContentFree, handle, "xmtp_reaction_content_free")
	purego.RegisterLibFunc(&lib.XmtpReplyContentFree, handle, "xmtp_reply_content_free")
	purego.RegisterLibFunc(&lib.XmtpAttachmentContentFree, handle, "xmtp_attachment_content_free")

	// Initialize the library
	lib.XmtpInit()

	return lib, nil
}

// GetLibrary returns the loaded library
func GetLibrary() (*Library, error) {
	if library == nil {
		return nil, ErrLibraryNotLoaded
	}
	return library, nil
}

// Version returns the library version
func Version() string {
	lib, err := GetLibrary()
	if err != nil {
		return ""
	}
	return GoString(lib.XmtpVersion())
}
