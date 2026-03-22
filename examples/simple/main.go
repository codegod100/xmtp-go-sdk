package main

import (
	"context"
	"fmt"
	"log"

	"github.com/xmtp/go-sdk"
)

func main() {
	// Create a client connected to XMTP dev network
	// Note: This requires an EthSigner implementation
	// For now, this demonstrates the API structure

	fmt.Println("XMTP Go SDK Example")
	fmt.Println("===================")

	// The client requires a signer for identity
	// You would implement EthSigner or use a provided implementation
	//
	// Example usage:
	//   client, err := xmtp.NewClient(ctx, signer,
	//       xmtp.WithEnv(xmtp.EnvDev),
	//       xmtp.WithAppVersion("my-app/1.0"),
	//   )

	// For testing the FFI layer directly:
	fmt.Println("To test the FFI layer, run:")
	fmt.Println("  go test -v ./internal/ffi/...")
	fmt.Println()
	fmt.Println("Available environments:")
	fmt.Println("  EnvLocal          - localhost:5556")
	fmt.Println("  EnvDev            - dev.xmtp.network:5556")
	fmt.Println("  EnvTestnet        - testnet.xmtp.network:5556")
	fmt.Println("  EnvProduction     - production.xmtp.network:5556")
	fmt.Println()
	fmt.Println("Client options:")
	fmt.Println("  WithEnv(env)              - Set XMTP environment")
	fmt.Println("  WithDbPath(path)          - Set database path")
	fmt.Println("  WithAppVersion(version)   - Set app version")
	fmt.Println("  WithLogLevel(level)       - Set log level (0-5)")

	_ = context.Background()
	_ = xmtp.EnvDev

	// Placeholder
	log.Println("See internal/ffi/ for working FFI tests")
}