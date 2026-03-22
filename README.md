# XMTP Go SDK

A Go SDK for XMTP messaging, using [PureGo](https://github.com/ebitengine/purego) to call into libxmtp without CGO.

## Features

- **No CGO required** - Uses PureGo for dynamic library loading
- **Cross-platform** - Works on Linux, macOS, and Windows
- **Full XMTP v3 support** - Groups, DMs, MLS encryption
- **Idiomatic Go API** - Context support, channels for streaming

## Requirements

- Go 1.22+
- libxmtp_ffi shared library

## Installation

```bash
go get github.com/xmtp/go-sdk
```

## Quick Start

### With Ethereum Signer

```go
package main

import (
    "context"
    "fmt"
    "log"

    xmtp "github.com/xmtp/go-sdk"
)

func main() {
    // Create client from private key
    client, signer, err := xmtp.QuickClient(
        context.Background(),
        "your-private-key-hex",
        xmtp.WithEnv(xmtp.EnvDev),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    fmt.Printf("Inbox ID: %s\n", client.InboxID())

    // Register if needed
    if !client.IsRegistered() {
        if err := client.Register(context.Background()); err != nil {
            log.Fatal(err)
        }
    }

    // Create a DM
    dm, err := client.Conversations().CreateDM(
        context.Background(),
        "peer-inbox-id",
    )
    if err != nil {
        log.Fatal(err)
    }

    // Send a message
    msgID, err := dm.SendText(context.Background(), "Hello!")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Sent message: %s\n", msgID)
}
```

### Creating a Group

```go
// Create a group with members
group, err := client.Conversations().CreateGroup(
    context.Background(),
    []string{"inbox-id-1", "inbox-id-2"},
    &xmtp.CreateGroupOptions{
        Name:        "My Group",
        Description: "A group chat",
    },
)
if err != nil {
    log.Fatal(err)
}

// Update group properties
group.UpdateName(context.Background(), "New Name")
group.AddMembers(context.Background(), []string{"inbox-id-3"})

// Send a message
msgID, err := group.SendText(context.Background(), "Hello group!")
```

### Streaming Messages

```go
// Stream all new messages
msgChan, err := client.Conversations().StreamAllMessages(
    context.Background(),
    &xmtp.StreamOptions{
        OnError: func(err error) {
            log.Printf("Stream error: %v", err)
        },
    },
)
if err != nil {
    log.Fatal(err)
}

for msg := range msgChan {
    fmt.Printf("New message from %s: %s\n", msg.SenderInboxID, msg.Text())
}
```

## Building the FFI Library

The Go SDK requires the `libxmtp_ffi` shared library. Build it from the Rust source:

```bash
# Clone libxmtp (if not already cloned)
git clone https://github.com/xmtp/libxmtp.git
cd libxmtp

# Copy the ffi crate from this SDK
cp -r /path/to/xmtp-go-sdk/ffi ./bindings/

# Build the FFI library
cd bindings/ffi
cargo build --release
```

This produces:
- Linux: `target/release/libxmtp_ffi.so`
- macOS: `target/release/libxmtp_ffi.dylib`
- Windows: `target/release/xmtp_ffi.dll`

Copy the library to a location where the SDK can find it:

```bash
# Option 1: Current directory
cp target/release/libxmtp_ffi.so .

# Option 2: System library path (Linux/macOS)
sudo cp target/release/libxmtp_ffi.so /usr/local/lib/
sudo ldconfig  # Linux only

# Option 3: Set environment variable
export XMTP_FFI_PATH=/path/to/libxmtp_ffi.so
```

## Configuration

### Client Options

```go
client, err := xmtp.NewClient(ctx, signer,
    xmtp.WithEnv(xmtp.EnvProduction),       // Environment
    xmtp.WithDbPath("./xmtp.db"),           // Database path
    xmtp.WithDbEncryptionKey(key32bytes),   // Encryption key (32 bytes)
    xmtp.WithAppVersion("myapp/1.0.0"),     // App version
    xmtp.WithDisableAutoRegister(),         // Disable auto-registration
    xmtp.WithStructuredLogging(),           // JSON logging
    xmtp.WithLogLevel(4),                   // 0-5: off, error, warn, info, debug, trace
)
```

### Environments

| Env | Description |
|-----|-------------|
| `EnvLocal` | Local development (localhost) |
| `EnvDev` | XMTP dev network |
| `EnvProduction` | XMTP production network |
| `EnvTestnet` | XMTP testnet |
| `EnvMainnet` | XMTP mainnet |

## API Reference

### Client

```go
// Create client
client, err := xmtp.NewClient(ctx, signer, opts...)

// Client properties
client.InboxID()           // string
client.InstallationID()    // string
client.IsRegistered()      // bool
client.Conversations()     // *Conversations

// Methods
client.Register(ctx)                              // error
client.CanMessage(ctx, identifiers)               // (map[string]bool, error)
client.FetchInboxIDByIdentifier(ctx, identifier)  // (string, error)
client.Close()                                    // error
```

### Conversations

```go
convos := client.Conversations()

// List conversations
all, _ := convos.List(ctx)           // []Conversation
groups, _ := convos.ListGroups(ctx)  // []*Group
dms, _ := convos.ListDMs(ctx)        // []*DM

// Get by ID
conv, _ := convos.GetByID(ctx, "id") // Conversation
dm, _ := convos.GetDMByInboxID(ctx, "inbox-id") // *DM

// Create new conversations
group, _ := convos.CreateGroup(ctx, []string{"inbox-id"}, &xmtp.CreateGroupOptions{...})
dm, _ := convos.CreateDM(ctx, "inbox-id")

// Sync
convos.Sync(ctx)

// Stream
msgChan, _ := convos.StreamAllMessages(ctx, opts)
```

### Group

```go
// Properties
group.ID()            // string
group.Name()          // string
group.ImageURL()      // string
group.Description()   // string
group.IsActive()      // bool
group.ConsentState()  // ConsentState

// Update properties
group.UpdateName(ctx, "New Name")
group.UpdateImageURL(ctx, "https://...")
group.UpdateDescription(ctx, "Description")

// Members
members, _ := group.ListMembers(ctx)
group.AddMembers(ctx, []string{"inbox-id"})
group.RemoveMembers(ctx, []string{"inbox-id"})

// Admin management
group.IsAdmin("inbox-id")       // bool
group.IsSuperAdmin("inbox-id")  // bool
group.AddAdmin(ctx, "inbox-id")
group.RemoveAdmin(ctx, "inbox-id")

// Messages
group.SendText(ctx, "Hello!")
group.SendMarkdown(ctx, "**Hello!**")
group.SendReaction(ctx, "msg-id", "👍")
messages, _ := group.Messages(ctx, &xmtp.ListMessagesOptions{Limit: 50})

// Leave
group.Leave(ctx)
```

### DM

```go
// Properties (inherits from Conversation)
dm.ID()
dm.IsActive()
dm.ConsentState()
dm.PeerInboxID()      // string - the other participant

// Messages (same as Group)
dm.SendText(ctx, "Hello!")
dm.SendMarkdown(ctx, "**Hello!**")
dm.SendReaction(ctx, "msg-id", "👍")
messages, _ := dm.Messages(ctx, opts)
```

### Message

```go
msg.ID()              // string
msg.SenderInboxID     // string
msg.ConversationID    // string
msg.SentAt            // time.Time
msg.ExpiresAt         // *time.Time
msg.ContentType       // ContentType
msg.DeliveryStatus    // DeliveryStatus
msg.Content           // any (string, *Reaction, *Reply, *Attachment, etc.)

// Type checks
msg.IsText()          // bool
msg.IsMarkdown()      // bool
msg.IsReaction()      // bool
msg.IsReply()         // bool

// Content access
msg.Text()            // string (if text/markdown)
msg.GetReaction()     // *Reaction (if reaction)
msg.GetReply()        // *Reply (if reply)
msg.GetAttachment()   // *Attachment (if attachment)
```

## Custom Content Types

The SDK supports custom content types through codecs:

```go
// TODO: Document custom codec registration
```

## Error Handling

```go
import "errors"

// Common errors
errors.Is(err, xmtp.ErrClientNotInitialized)
errors.Is(err, xmtp.ErrSignerUnavailable)
errors.Is(err, xmtp.ErrConversationNotFound)
errors.Is(err, xmtp.ErrMessageNotFound)
errors.Is(err, xmtp.ErrLibraryNotLoaded)
```

## Architecture

```
┌─────────────────────────────────────────┐
│           Go Application                 │
├─────────────────────────────────────────┤
│          XMTP Go SDK                     │
│  (Client, Conversations, Group, DM...)   │
├─────────────────────────────────────────┤
│         internal/ffi (PureGo)            │
│    (Dynamic library loading via          │
│     dlopen/dlsym - no CGO)               │
├─────────────────────────────────────────┤
│        libxmtp_ffi.so/.dylib/.dll        │
│    (Rust FFI layer with C ABI)           │
├─────────────────────────────────────────┤
│            libxmtp (Rust)                │
│    (Core XMTP protocol implementation)   │
└─────────────────────────────────────────┘
```

## Testing

```bash
# Run tests
go test ./...

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Contributing

Contributions are welcome! Please read the contributing guidelines before submitting PRs.

## License

MIT
