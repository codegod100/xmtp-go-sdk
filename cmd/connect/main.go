package main

import (
	"fmt"
	"log"

	"github.com/xmtp/go-sdk/internal/ffi"
)

func main() {
	if !ffi.IsLoaded() {
		log.Fatal("Library not loaded")
	}

	fmt.Println("Attempting to connect to XMTP dev network...")
	
	// Try to connect to dev network
	handle, err := ffi.ConnectToBackend("dev.xmtp.network:5556", "", "test-app/1.0")
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	
	fmt.Printf("Connected successfully! Handle: %d\n", handle)
	
	// Note: FreeClient has a bug, skipping for now
	// The client handle will be cleaned up when the process exits
	
	fmt.Println("Done!")
}
