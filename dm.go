package xmtp

// DM represents a direct message conversation
type DM struct {
	baseConversation
}

// PeerInboxID returns the inbox ID of the other participant
func (d *DM) PeerInboxID() string {
	return "" // TODO: Implement via FFI
}
