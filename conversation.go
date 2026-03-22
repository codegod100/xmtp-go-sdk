package xmtp

import (
	"context"
	"sync"
	"time"
)

// Conversation is the interface for all conversation types
type Conversation interface {
	ID() string
	IsActive() bool
	CreatedAt() time.Time
	ConsentState() ConsentState
	UpdateConsent(ctx context.Context, state ConsentState) error
	Sync(ctx context.Context) error
	SendText(ctx context.Context, text string) (string, error)
	SendMarkdown(ctx context.Context, markdown string) (string, error)
	SendReaction(ctx context.Context, referenceID string, content string) (string, error)
	Messages(ctx context.Context, opts *ListMessagesOptions) ([]*Message, error)
	StreamMessages(ctx context.Context, opts *StreamOptions) (<-chan *Message, error)
}

// baseConversation provides common functionality for all conversation types
type baseConversation struct {
	handle uintptr
	client *Client
	mu     sync.RWMutex
}

// ID returns the conversation ID
func (c *baseConversation) ID() string {
	return "" // TODO: Implement via FFI
}

// IsActive returns whether the conversation is active
func (c *baseConversation) IsActive() bool {
	return false // TODO: Implement via FFI
}

// CreatedAt returns the creation timestamp
func (c *baseConversation) CreatedAt() time.Time {
	return time.Time{} // TODO: Implement via FFI
}

// ConsentState returns the consent state
func (c *baseConversation) ConsentState() ConsentState {
	return ConsentUnknown // TODO: Implement via FFI
}

// UpdateConsent updates the consent state
func (c *baseConversation) UpdateConsent(ctx context.Context, state ConsentState) error {
	return nil // TODO: Implement via FFI
}

// Sync synchronizes the conversation from the network
func (c *baseConversation) Sync(ctx context.Context) error {
	return nil // TODO: Implement via FFI
}

// SendText sends a text message
func (c *baseConversation) SendText(ctx context.Context, text string) (string, error) {
	return "", nil // TODO: Implement via FFI
}

// SendMarkdown sends a markdown message
func (c *baseConversation) SendMarkdown(ctx context.Context, markdown string) (string, error) {
	return "", nil // TODO: Implement via FFI
}

// SendReaction sends a reaction to a message
func (c *baseConversation) SendReaction(ctx context.Context, referenceID string, content string) (string, error) {
	return "", nil // TODO: Implement via FFI
}

// Messages lists messages in the conversation
func (c *baseConversation) Messages(ctx context.Context, opts *ListMessagesOptions) ([]*Message, error) {
	return nil, nil // TODO: Implement via FFI
}

// StreamMessages returns a channel for new messages
func (c *baseConversation) StreamMessages(ctx context.Context, opts *StreamOptions) (<-chan *Message, error) {
	return nil, nil // TODO: Implement via FFI
}
