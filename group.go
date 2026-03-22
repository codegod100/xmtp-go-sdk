package xmtp

import (
	"context"
)

// Group represents a group conversation
type Group struct {
	baseConversation
}

// Name returns the group name
func (g *Group) Name() string {
	return "" // TODO: Implement via FFI
}

// UpdateName updates the group name
func (g *Group) UpdateName(ctx context.Context, name string) error {
	return nil // TODO: Implement via FFI
}

// ImageURL returns the group image URL
func (g *Group) ImageURL() string {
	return "" // TODO: Implement via FFI
}

// UpdateImageURL updates the group image URL
func (g *Group) UpdateImageURL(ctx context.Context, url string) error {
	return nil // TODO: Implement via FFI
}

// Description returns the group description
func (g *Group) Description() string {
	return "" // TODO: Implement via FFI
}

// UpdateDescription updates the group description
func (g *Group) UpdateDescription(ctx context.Context, description string) error {
	return nil // TODO: Implement via FFI
}

// ListMembers returns the group members
func (g *Group) ListMembers(ctx context.Context) ([]GroupMember, error) {
	return nil, nil // TODO: Implement via FFI
}

// AddMembers adds members to the group
func (g *Group) AddMembers(ctx context.Context, inboxIDs []string) error {
	return nil // TODO: Implement via FFI
}

// RemoveMembers removes members from the group
func (g *Group) RemoveMembers(ctx context.Context, inboxIDs []string) error {
	return nil // TODO: Implement via FFI
}

// IsAdmin checks if an inbox ID is an admin
func (g *Group) IsAdmin(inboxID string) bool {
	return false // TODO: Implement via FFI
}

// IsSuperAdmin checks if an inbox ID is a super admin
func (g *Group) IsSuperAdmin(inboxID string) bool {
	return false // TODO: Implement via FFI
}

// AddAdmin promotes a member to admin
func (g *Group) AddAdmin(ctx context.Context, inboxID string) error {
	return nil // TODO: Implement via FFI
}

// RemoveAdmin demotes an admin to member
func (g *Group) RemoveAdmin(ctx context.Context, inboxID string) error {
	return nil // TODO: Implement via FFI
}

// Leave leaves the group
func (g *Group) Leave(ctx context.Context) error {
	return nil // TODO: Implement via FFI
}
