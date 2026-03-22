package ffi

import (
	"context"
)

// Client represents an XMTP client
type Client struct {
	handle uint64
	api    uint64
}

// Conversation represents an XMTP conversation
type Conversation struct {
	handle uint64
	client *Client
}

// ClientOptions for creating a new XMTP client
type ClientOptions struct {
	V3Host      string
	GatewayHost string
	AppVersion  string
	DBPath      string
}

// NewClient creates a new XMTP client
func NewClient(ctx context.Context, opts ClientOptions) (*Client, error) {
	if !IsLoaded() {
		return nil, LoadError()
	}

	appVersion := opts.AppVersion
	if appVersion == "" {
		appVersion = "xmtp-go-sdk/0.1.0"
	}

	apiHandle, err := ConnectToBackend(opts.V3Host, opts.GatewayHost, appVersion)
	if err != nil {
		return nil, err
	}

	return &Client{
		handle: 0,
		api:    apiHandle,
	}, nil
}

// Close frees the client resources
func (c *Client) Close() {
	if c.handle != 0 {
		FreeClient(c.handle)
		c.handle = 0	}
}

// InboxID returns the client's inbox ID
func (c *Client) InboxID() (string, error) {
	if c.handle == 0 {
		return "", nil
	}
	//TODO: implement
	return "", nil
}

// Conversations returns the list of conversations
func (c *Client) Conversations() ([]*Conversation, error) {
	if c.handle == 0 {
		return nil, nil
	}
	// TODO: implement
	return nil, nil
}

// Close frees the conversation handle
func (c *Conversation) Close() {
	if c.handle != 0 {
		FreeConversation(c.handle)
		c.handle = 0
	}
}