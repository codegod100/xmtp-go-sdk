package xmtp

import (
	"context"
	"fmt"
	"sync"

	"github.com/xmtp/go-sdk/internal/ffi"
)

// Client is the main XMTP client
type Client struct {
	mu           sync.RWMutex
	handle       uintptr
	signer       Signer
	env          Env
	options      *ClientOptions
	conversations *Conversations
}

// NewClient creates a new XMTP client with a signer
func NewClient(ctx context.Context, signer Signer, opts ...ClientOption) (*Client, error) {
	_, err := ffi.GetLibrary()
	if err != nil {
		return nil, err
	}

	options := &ClientOptions{
		Env:      EnvDev,
		LogLevel: 2, // Warn
	}
	for _, opt := range opts {
		opt(options)
	}

	// Get identifier from signer
	_, err = signer.GetIdentifier()
	if err != nil {
		return nil, fmt.Errorf("failed to get identifier: %w", err)
	}

	// TODO: Create client via FFI
	// For now, return a skeleton client
	client := &Client{
		handle:  0,
		signer:  signer,
		env:     options.Env,
		options: options,
	}

	// Create conversations manager
	client.conversations = &Conversations{
		client: client,
	}

	return client, nil
}

// Build creates a client from an existing identity (no signer required)
func Build(ctx context.Context, identifier Identifier, opts ...ClientOption) (*Client, error) {
	_, err := ffi.GetLibrary()
	if err != nil {
		return nil, err
	}

	options := &ClientOptions{
		Env:                EnvDev,
		LogLevel:           2,
		DisableAutoRegister: true,
	}
	for _, opt := range opts {
		opt(options)
	}

	// TODO: Build client via FFI
	client := &Client{
		handle:  0,
		env:     options.Env,
		options: options,
	}

	client.conversations = &Conversations{
		client: client,
	}

	return client, nil
}

// Close releases the client resources
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.handle = 0
	return nil
}

// InboxID returns the client's inbox ID
func (c *Client) InboxID() string {
	return "" // TODO: Implement via FFI
}

// InstallationID returns the client's installation ID
func (c *Client) InstallationID() string {
	return "" // TODO: Implement via FFI
}

// IsRegistered returns whether the client is registered
func (c *Client) IsRegistered() bool {
	return false // TODO: Implement via FFI
}

// Conversations returns the conversations manager
func (c *Client) Conversations() *Conversations {
	return c.conversations
}

// Register registers the client with the XMTP network
func (c *Client) Register(ctx context.Context) error {
	if c.signer == nil {
		return ErrSignerUnavailable
	}
	// TODO: Implement via FFI
	return nil
}

// CanMessage checks if the given identifiers can be messaged
func (c *Client) CanMessage(ctx context.Context, identifiers []Identifier) (map[string]bool, error) {
	return make(map[string]bool), nil // TODO: Implement via FFI
}

// FetchInboxIDByIdentifier fetches the inbox ID for an identifier
func (c *Client) FetchInboxIDByIdentifier(ctx context.Context, identifier Identifier) (string, error) {
	return "", nil // TODO: Implement via FFI
}

// LibxmtpVersion returns the libxmtp version
func LibxmtpVersion() string {
	lib, _ := ffi.GetLibrary()
	if lib == nil {
		return ""
	}
	return GoString(lib.XmtpClientLibxmtpVersion())
}

// Handle returns the underlying handle (for internal use)
func (c *Client) Handle() uintptr {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.handle
}

// GoString converts a uintptr C string to Go string
func GoString(s uintptr) string {
	return ffi.GoString(s)
}
