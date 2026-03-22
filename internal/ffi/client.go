package ffi

import (
	"context"
)

// Client represents an XMTP client
type Client struct {
	handle uint64
	api    uint64 // API client handle
}

// Conversation represents an XMTP conversation
type Conversation struct {
	handle uint64
	client *Client
}

// ClientOptions for creating a new XMTP client
type ClientOptions struct {
	V3Host      string
	GatewayHost string // Optional - enables D14n
	AppVersion  string
	DBPath      string // Empty = ephemeral
}

// NewClient creates a new XMTP client
func NewClient(ctx context.Context, opts ClientOptions) (*Client, error) {
	if !libLoaded {
		return nil, loadErr
	}
	
	appVersion := opts.AppVersion
	if appVersion == "" {
		appVersion = "xmtp-go-sdk/0.1.0"
	}
	
	// Convert strings to buffers
	v3HostBuf, err := stringToBuffer(opts.V3Host)
	if err != nil {
		return nil, err
	}
	defer freeBuffer(v3HostBuf)
	
	gatewayHostBuf, err := stringToBuffer(opts.GatewayHost)
	if err != nil {
		return nil, err
	}
	defer freeBuffer(gatewayHostBuf)
	
	clientModeBuf, err := stringToBuffer("ReadWrite")
	if err != nil {
		return nil, err
	}
	defer freeBuffer(clientModeBuf)
	
	appVersionBuf, err := stringToBuffer(appVersion)
	if err != nil {
		return nil, err
	}
	defer freeBuffer(appVersionBuf)
	
	// Connect to backend (async)
	futureHandle := ffi_connect_to_backend(v3HostBuf, gatewayHostBuf, clientModeBuf, appVersionBuf, rustBuffer{}, rustBuffer{})
	apiHandle, err := awaitFutureU64(futureHandle)
	if err != nil {
		return nil, err
	}
	ffi_future_free(futureHandle)
	
	// TODO: Create full client with identity
	// For now, just return the API handle
	return &Client{
		handle: 0, // No client handle yet
		api:    apiHandle,
	}, nil
}

// Close frees the client resources
func (c *Client) Close() {
	if c.handle != 0 && libLoaded {
		var status rustCallStatus
		ffi_free_client(c.handle, &status)
		c.handle = 0
	}
}

// InboxID returns the client's inbox ID
func (c *Client) InboxID() (string, error) {
	if c.handle == 0 {
		return "", nil
	}
	
	var status rustCallStatus
	buf := ffi_client_inbox_id(c.handle, &status)
	if status.code != CallSuccess {
		return "", statusToError(status)
	}
	
	result := bufferToString(buf)
	freeBuffer(buf)
	return result, nil
}

// Conversations returns the list of conversations
func (c *Client) Conversations() ([]*Conversation, error) {
	if c.handle == 0 {
		return nil, nil
	}
	
	var status rustCallStatus
	// TODO: Implement conversations list
	_ = ffi_client_conversations(c.handle, &status)
	if status.code != CallSuccess {
		return nil, statusToError(status)
	}
	
	return nil, nil
}

// SendText sends a text message to a conversation
func (c *Conversation) SendText(text string) error {
	if !libLoaded {
		return loadErr
	}
	
	textBuf, err := stringToBuffer(text)
	if err != nil {
		return err
	}
	defer freeBuffer(textBuf)
	
	var status rustCallStatus
	futureHandle := ffi_conversation_send_text(c.handle, textBuf, &status)
	if status.code != CallSuccess {
		return statusToError(status)
	}
	
	// Wait for send to complete
	_, err = awaitFutureU64(futureHandle)
	ffi_future_free(futureHandle)
	return err
}

// ID returns the conversation ID as bytes
func (c *Conversation) ID() ([]byte, error) {
	if !libLoaded {
		return nil, loadErr
	}
	
	var status rustCallStatus
	buf := ffi_conversation_id(c.handle, &status)
	if status.code != CallSuccess {
		return nil, statusToError(status)
	}
	
	result := bufferToBytes(buf)
	freeBuffer(buf)
	return result, nil
}

// Close frees the conversation handle
func (c *Conversation) Close() {
	if c.handle != 0 && libLoaded {
		var status rustCallStatus
		ffi_free_conversation(c.handle, &status)
		c.handle = 0
	}
}
