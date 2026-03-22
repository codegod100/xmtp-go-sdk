package xmtp

import (
	"context"
	"errors"
	"time"
)

// Preferences manages user preferences and consent
type Preferences struct {
	client *Client
}

// InboxState represents the state of an inbox
type InboxState struct {
	InboxID            string
	RecoveryIdentifier *Identifier
	Identifiers        []Identifier
	Installations      []Installation
}

// Installation represents a client installation
type Installation struct {
	ID        string
	CreatedAt time.Time
}

// InboxState returns the current inbox state from local database
func (p *Preferences) InboxState(ctx context.Context) (*InboxState, error) {
	return nil, errors.New("not implemented")
}

// FetchInboxState fetches the latest inbox state from the network
func (p *Preferences) FetchInboxState(ctx context.Context) (*InboxState, error) {
	return nil, errors.New("not implemented")
}

// SetConsentStates updates consent states for multiple entities
func (p *Preferences) SetConsentStates(ctx context.Context, states []ConsentRecord) error {
	return errors.New("not implemented")
}

// GetConsentState gets the consent state for an entity
func (p *Preferences) GetConsentState(ctx context.Context, entityType int, entity string) (ConsentState, error) {
	return ConsentUnknown, errors.New("not implemented")
}

// StreamConsent streams consent state updates
func (p *Preferences) StreamConsent(ctx context.Context, opts *StreamOptions) (<-chan []ConsentRecord, error) {
	return nil, errors.New("not implemented")
}

// ConsentRecord represents a consent record
type ConsentRecord struct {
	EntityType int
	Entity     string
	State      ConsentState
}

// DeviceSync manages device synchronization
type DeviceSync struct {
	client *Client
}

// SendSyncRequest sends a sync request to other devices
func (d *DeviceSync) SendSyncRequest(ctx context.Context, opts *SyncOptions) error {
	return errors.New("not implemented")
}

// SyncOptions represents options for sync operations
type SyncOptions struct {
	Elements                   []SyncElement
	ExcludeDisappearingMessages bool
}

// SyncElement represents a sync element type
type SyncElement int

const (
	SyncElementConsent SyncElement = iota
	SyncElementMessages
)

// DebugInformation provides debugging utilities
type DebugInformation struct {
	client *Client
}

// ConversationDebugInfo represents debug info for a conversation
type ConversationDebugInfo struct {
	ID                 string
	CreatedAtNs        int64
	LastMessageAtNs    int64
	NumMessages        int
	NumPendingMessages int
}
