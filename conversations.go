package xmtp

import (
	"context"
	"sync"
)

// Conversations manages all conversations for a client
type Conversations struct {
	client *Client
	mu     sync.RWMutex
	handle uintptr
}

// List returns all conversations
func (c *Conversations) List(ctx context.Context) ([]Conversation, error) {
	return nil, nil // TODO: Implement via FFI
}

// ListGroups returns all group conversations
func (c *Conversations) ListGroups(ctx context.Context) ([]*Group, error) {
	return nil, nil // TODO: Implement via FFI
}

// ListDMs returns all DM conversations
func (c *Conversations) ListDMs(ctx context.Context) ([]*DM, error) {
	return nil, nil // TODO: Implement via FFI
}

// GetByID returns a conversation by ID
func (c *Conversations) GetByID(ctx context.Context, id string) (Conversation, error) {
	return nil, nil // TODO: Implement via FFI
}

// GetDMByInboxID returns a DM by inbox ID
func (c *Conversations) GetDMByInboxID(ctx context.Context, inboxID string) (*DM, error) {
	return nil, nil // TODO: Implement via FFI
}

// CreateGroup creates a new group conversation
func (c *Conversations) CreateGroup(ctx context.Context, inboxIDs []string, opts *CreateGroupOptions) (*Group, error) {
	return nil, nil // TODO: Implement via FFI
}

// CreateDM creates a new DM conversation
func (c *Conversations) CreateDM(ctx context.Context, inboxID string) (*DM, error) {
	return nil, nil // TODO: Implement via FFI
}

// Sync synchronizes conversations from the network
func (c *Conversations) Sync(ctx context.Context) error {
	return nil // TODO: Implement via FFI
}

// Stream returns a channel for new conversations
func (c *Conversations) Stream(ctx context.Context, opts *StreamOptions) (<-chan Conversation, error) {
	return nil, nil // TODO: Implement via FFI
}

// StreamGroups returns a channel for new group conversations
func (c *Conversations) StreamGroups(ctx context.Context, opts *StreamOptions) (<-chan *Group, error) {
	return nil, nil // TODO: Implement via FFI
}

// StreamDMs returns a channel for new DM conversations
func (c *Conversations) StreamDMs(ctx context.Context, opts *StreamOptions) (<-chan *DM, error) {
	return nil, nil // TODO: Implement via FFI
}

// StreamAllMessages returns a channel for all new messages
func (c *Conversations) StreamAllMessages(ctx context.Context, opts *StreamOptions) (<-chan *Message, error) {
	return nil, nil // TODO: Implement via FFI
}
