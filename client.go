package xmtp

import (
	"context"
	"fmt"

	"github.com/xmtp/go-sdk/internal/ffi"
)

// Client represents an XMTP client
type Client struct {
	apiHandle uint64
	opts      ClientOptions
}

// NewClient creates a new XMTP client
func NewClient(ctx context.Context, signer Signer, opts ...ClientOption) (*Client, error) {
	// Apply default options
	options := ClientOptions{
		Env:        EnvDev,
		AppVersion: "xmtp-go-sdk/0.1.0",
	}
	for _, opt := range opts {
		opt(&options)
	}

	// Get host based on environment
	v3Host := getHostForEnv(options.Env)

	apiHandle, err := ffi.ConnectToBackend(v3Host, "", options.AppVersion)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to backend: %w", err)
	}

	return &Client{
		apiHandle: apiHandle,
		opts:      options,
	}, nil
}

// Close frees client resources
func (c *Client) Close() error {
	if c.apiHandle != 0 {
		if err := ffi.FreeClient(c.apiHandle); err != nil {
			return err
		}
		c.apiHandle = 0
	}
	return nil
}

func getHostForEnv(env Env) string {
	switch env {
	case EnvLocal:
		return "localhost:5556"
	case EnvDev:
		return "dev.xmtp.network:5556"
	case EnvProduction, EnvMainnet:
		return "production.xmtp.network:5556"
	case EnvTestnetStaging:
		return "staging.testnet.xmtp.network:5556"
	case EnvTestnetDev:
		return "dev.testnet.xmtp.network:5556"
	case EnvTestnet:
		return "testnet.xmtp.network:5556"
	default:
		return "dev.xmtp.network:5556"
	}
}