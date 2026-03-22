package xmtp

import (
	"context"
	"testing"
	"time"
)

// TestClientOptions tests client option configuration
func TestClientOptions(t *testing.T) {
	tests := []struct {
		name     string
		opts     []ClientOption
		expected ClientOptions
	}{
		{
			name:     "default options",
			opts:     nil,
			expected: ClientOptions{Env: EnvDev, LogLevel: 2},
		},
		{
			name: "production environment",
			opts: []ClientOption{WithEnv(EnvProduction)},
			expected: ClientOptions{Env: EnvProduction, LogLevel: 2},
		},
		{
			name: "custom db path",
			opts: []ClientOption{WithDbPath("/tmp/xmtp.db")},
			expected: ClientOptions{Env: EnvDev, DbPath: "/tmp/xmtp.db", LogLevel: 2},
		},
		{
			name: "with encryption key",
			opts: []ClientOption{WithDbEncryptionKey(make([]byte, 32))},
			expected: ClientOptions{Env: EnvDev, DbEncryptionKey: make([]byte, 32), LogLevel: 2},
		},
		{
			name: "with app version",
			opts: []ClientOption{WithAppVersion("test/1.0.0")},
			expected: ClientOptions{Env: EnvDev, AppVersion: "test/1.0.0", LogLevel: 2},
		},
		{
			name: "disable auto register",
			opts: []ClientOption{WithDisableAutoRegister()},
			expected: ClientOptions{Env: EnvDev, DisableAutoRegister: true, LogLevel: 2},
		},
		{
			name: "structured logging",
			opts: []ClientOption{WithStructuredLogging()},
			expected: ClientOptions{Env: EnvDev, StructuredLogging: true, LogLevel: 2},
		},
		{
			name: "custom log level",
			opts: []ClientOption{WithLogLevel(4)},
			expected: ClientOptions{Env: EnvDev, LogLevel: 4},
		},
		{
			name: "multiple options",
			opts: []ClientOption{
				WithEnv(EnvTestnet),
				WithDbPath("./test.db"),
				WithAppVersion("app/1.0"),
				WithLogLevel(3),
			},
			expected: ClientOptions{
				Env:       EnvTestnet,
				DbPath:    "./test.db",
				AppVersion: "app/1.0",
				LogLevel:  3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			options := &ClientOptions{
				Env:      EnvDev,
				LogLevel: 2,
			}
			for _, opt := range tt.opts {
				opt(options)
			}

			if options.Env != tt.expected.Env {
				t.Errorf("Env: got %v, want %v", options.Env, tt.expected.Env)
			}
			if options.DbPath != tt.expected.DbPath {
				t.Errorf("DbPath: got %v, want %v", options.DbPath, tt.expected.DbPath)
			}
			if options.AppVersion != tt.expected.AppVersion {
				t.Errorf("AppVersion: got %v, want %v", options.AppVersion, tt.expected.AppVersion)
			}
			if options.DisableAutoRegister != tt.expected.DisableAutoRegister {
				t.Errorf("DisableAutoRegister: got %v, want %v", options.DisableAutoRegister, tt.expected.DisableAutoRegister)
			}
			if options.StructuredLogging != tt.expected.StructuredLogging {
				t.Errorf("StructuredLogging: got %v, want %v", options.StructuredLogging, tt.expected.StructuredLogging)
			}
			if options.LogLevel != tt.expected.LogLevel {
				t.Errorf("LogLevel: got %v, want %v", options.LogLevel, tt.expected.LogLevel)
			}
		})
	}
}

// TestIdentifier tests identifier creation
func TestIdentifier(t *testing.T) {
	ident := Identifier{
		Kind:       0, // Ethereum
		Identifier: "0x1234567890123456789012345678901234567890",
	}

	if ident.Kind != 0 {
		t.Errorf("Kind: got %v, want 0", ident.Kind)
	}
	if ident.Identifier != "0x1234567890123456789012345678901234567890" {
		t.Errorf("Identifier: got %v", ident.Identifier)
	}
}

// TestMessage tests message type methods
func TestMessage(t *testing.T) {
	now := time.Now()
	
	t.Run("text message", func(t *testing.T) {
		msg := &Message{
			ID:             "msg-1",
			SenderInboxID:  "inbox-1",
			ConversationID: "conv-1",
			SentAt:         now,
			ContentType:    ContentTypeText,
			DeliveryStatus: DeliveryStatusPublished,
			Content:        "Hello, World!",
		}

		if !msg.IsText() {
			t.Error("IsText() should return true")
		}
		if msg.Text() != "Hello, World!" {
			t.Errorf("Text(): got %v, want 'Hello, World!'", msg.Text())
		}
	})

	t.Run("markdown message", func(t *testing.T) {
		msg := &Message{
			ID:          "msg-2",
			ContentType: ContentTypeMarkdown,
			Content:     "**Bold text**",
		}

		if !msg.IsMarkdown() {
			t.Error("IsMarkdown() should return true")
		}
		if msg.Text() != "**Bold text**" {
			t.Errorf("Text(): got %v", msg.Text())
		}
	})

	t.Run("reaction message", func(t *testing.T) {
		msg := &Message{
			ID:          "msg-3",
			ContentType: ContentTypeReaction,
			Content: &Reaction{
				ReferenceID: "ref-1",
				Action:      0, // add
				Schema:      0, // unicode
				Content:     "👍",
			},
		}

		if !msg.IsReaction() {
			t.Error("IsReaction() should return true")
		}
		if msg.GetReaction() == nil {
			t.Error("GetReaction() should not return nil")
		}
		if msg.GetReaction().Content != "👍" {
			t.Errorf("Reaction content: got %v", msg.GetReaction().Content)
		}
	})

	t.Run("reply message", func(t *testing.T) {
		msg := &Message{
			ID:          "msg-4",
			ContentType: ContentTypeReply,
			Content: &Reply{
				ReferenceID: "ref-1",
				ContentType: ContentTypeText,
				Content:     "Reply text",
			},
		}

		if !msg.IsReply() {
			t.Error("IsReply() should return true")
		}
		if msg.GetReply() == nil {
			t.Error("GetReply() should not return nil")
		}
	})

	t.Run("attachment message", func(t *testing.T) {
		msg := &Message{
			ID:          "msg-5",
			ContentType: ContentTypeAttachment,
			Content: &Attachment{
				Filename: "test.pdf",
				MimeType: "application/pdf",
				Data:     []byte{1, 2, 3},
			},
		}

		if !msg.IsAttachment() {
			t.Error("IsAttachment() should return true")
		}
		if msg.GetAttachment() == nil {
			t.Error("GetAttachment() should not return nil")
		}
		if msg.GetAttachment().Filename != "test.pdf" {
			t.Errorf("Filename: got %v", msg.GetAttachment().Filename)
		}
	})

	t.Run("with expiry", func(t *testing.T) {
		expiry := now.Add(24 * time.Hour)
		msg := &Message{
			ID:        "msg-6",
			SentAt:    now,
			ExpiresAt: &expiry,
		}

		if msg.ExpiresAt == nil {
			t.Error("ExpiresAt should not be nil")
		}
		if !msg.ExpiresAt.After(msg.SentAt) {
			t.Error("ExpiresAt should be after SentAt")
		}
	})
}

// TestConsentState tests consent state constants
func TestConsentState(t *testing.T) {
	states := []ConsentState{
		ConsentUnknown,
		ConsentAllowed,
		ConsentDenied,
	}

	for i, state := range states {
		if int(state) != i {
			t.Errorf("ConsentState %d: got %v", i, state)
		}
	}
}

// TestContentType tests content type constants
func TestContentType(t *testing.T) {
	types := []ContentType{
		ContentTypeText,
		ContentTypeMarkdown,
		ContentTypeReply,
		ContentTypeReaction,
		ContentTypeAttachment,
		ContentTypeRemoteAttachment,
		ContentTypeMultiRemoteAttachment,
		ContentTypeTransactionReference,
		ContentTypeGroupUpdated,
		ContentTypeReadReceipt,
		ContentTypeLeaveRequest,
		ContentTypeWalletSendCalls,
		ContentTypeActions,
		ContentTypeIntent,
		ContentTypeDeletedMessage,
		ContentTypeCustom,
	}

	for i, ct := range types {
		if int(ct) != i {
			t.Errorf("ContentType %d: got %v", i, ct)
		}
	}
}

// TestDeliveryStatus tests delivery status constants
func TestDeliveryStatus(t *testing.T) {
	statuses := []DeliveryStatus{
		DeliveryStatusUnpublished,
		DeliveryStatusPublished,
		DeliveryStatusFailed,
	}

	for i, status := range statuses {
		if int(status) != i {
			t.Errorf("DeliveryStatus %d: got %v", i, status)
		}
	}
}

// TestEnv tests environment constants
func TestEnv(t *testing.T) {
	envs := []Env{
		EnvLocal,
		EnvDev,
		EnvProduction,
		EnvTestnetStaging,
		EnvTestnetDev,
		EnvTestnet,
		EnvMainnet,
	}

	for i, env := range envs {
		if int(env) != i {
			t.Errorf("Env %d: got %v", i, env)
		}
	}
}

// TestErrors tests error values
func TestErrors(t *testing.T) {
	errors := []error{
		ErrClientNotInitialized,
		ErrSignerUnavailable,
		ErrConversationNotFound,
		ErrMessageNotFound,
		ErrNotRegistered,
		ErrLibraryNotLoaded,
	}

	for _, err := range errors {
		if err == nil {
			t.Error("Error should not be nil")
		}
		if err.Error() == "" {
			t.Error("Error message should not be empty")
		}
	}
}

// TestStream tests stream functionality
func TestStream(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream := NewStream[int](ctx, 10)

	// Test pushing values
	stream.Push(1)
	stream.Push(2)
	stream.Push(3)

	// Test collecting values
	go func() {
		time.Sleep(10 * time.Millisecond)
		stream.Close()
	}()

	values, err := stream.Collect(ctx)
	if err != nil {
		t.Errorf("Collect error: %v", err)
	}
	if len(values) != 3 {
		t.Errorf("Expected 3 values, got %d", len(values))
	}
}

// TestCreateGroupOptions tests create group options
func TestCreateGroupOptions(t *testing.T) {
	opts := &CreateGroupOptions{
		Name:        "Test Group",
		ImageURL:    "https://example.com/image.png",
		Description: "A test group",
		Permissions: 0,
	}

	if opts.Name != "Test Group" {
		t.Errorf("Name: got %v", opts.Name)
	}
	if opts.ImageURL != "https://example.com/image.png" {
		t.Errorf("ImageURL: got %v", opts.ImageURL)
	}
	if opts.Description != "A test group" {
		t.Errorf("Description: got %v", opts.Description)
	}
}

// TestListMessagesOptions tests list messages options
func TestListMessagesOptions(t *testing.T) {
	before := time.Now().Add(-24 * time.Hour)
	after := time.Now().Add(-48 * time.Hour)

	opts := &ListMessagesOptions{
		Limit:     50,
		Before:    &before,
		After:     &after,
		Ascending: true,
	}

	if opts.Limit != 50 {
		t.Errorf("Limit: got %v", opts.Limit)
	}
	if !opts.Ascending {
		t.Error("Ascending should be true")
	}
	if opts.Before == nil || opts.After == nil {
		t.Error("Before and After should not be nil")
	}
}

