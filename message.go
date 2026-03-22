package xmtp

// DecodeMessage decodes a message from an FFI handle
func DecodeMessage(handle uintptr) *Message {
	if handle == 0 {
		return nil
	}
	// TODO: Implement via FFI
	return &Message{}
}

// Text returns the text content if this is a text message
func (m *Message) Text() string {
	if s, ok := m.Content.(string); ok {
		return s
	}
	return ""
}

// IsText returns true if this is a text message
func (m *Message) IsText() bool {
	return m.ContentType == ContentTypeText
}

// IsMarkdown returns true if this is a markdown message
func (m *Message) IsMarkdown() bool {
	return m.ContentType == ContentTypeMarkdown
}

// IsReaction returns true if this is a reaction
func (m *Message) IsReaction() bool {
	return m.ContentType == ContentTypeReaction
}

// IsReply returns true if this is a reply
func (m *Message) IsReply() bool {
	return m.ContentType == ContentTypeReply
}

// IsAttachment returns true if this is an attachment
func (m *Message) IsAttachment() bool {
	return m.ContentType == ContentTypeAttachment
}

// GetReaction returns the reaction if this is a reaction message
func (m *Message) GetReaction() *Reaction {
	if r, ok := m.Content.(*Reaction); ok {
		return r
	}
	return nil
}

// GetReply returns the reply if this is a reply message
func (m *Message) GetReply() *Reply {
	if r, ok := m.Content.(*Reply); ok {
		return r
	}
	return nil
}

// GetAttachment returns the attachment if this is an attachment message
func (m *Message) GetAttachment() *Attachment {
	if a, ok := m.Content.(*Attachment); ok {
		return a
	}
	return nil
}
