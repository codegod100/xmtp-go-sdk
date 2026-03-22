{
  description = "XMTP Go SDK - Go bindings for XMTP messaging using PureGo";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    rust-overlay.url = "github:oxalica/rust-overlay";
  };

  outputs = { self, nixpkgs, flake-utils, rust-overlay, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [ (import rust-overlay) ];
        pkgs = import nixpkgs {
          inherit system overlays;
        };
        
        # Rust toolchain for building the FFI library
        rustToolchain = pkgs.rust-bin.stable.latest.default.override {
          extensions = [ "rust-src" "rust-analyzer" ];
          targets = [ "wasm32-unknown-unknown" ];
        };
        
        # Build inputs common to both build and dev
        buildInputs = with pkgs; [
          openssl
        ] ++ lib.optionals stdenv.isLinux [
          systemd
        ] ++ lib.optionals stdenv.isDarwin [
          darwin.apple_sdk.frameworks.Security
          darwin.apple_sdk.frameworks.CoreFoundation
          darwin.apple_sdk.frameworks.SystemConfiguration
        ];
        
        # Native build inputs
        nativeBuildInputs = with pkgs; [
          pkg-config
        ];
        
        lib = nixpkgs.lib;
      in {
        packages = {
          default = self.packages.${system}.xmtp-go-sdk;
          
          xmtp-go-sdk = pkgs.buildGo124Module {
            pname = "xmtp-go-sdk";
            version = "0.1.0";
            src = ./.;
            
            vendorHash = "sha256-7SrehajxKazbFz7m9YbslCOrI3U+NDEzxWroo5Jy8VU=";
            modRoot = ".";
            
            inherit buildInputs nativeBuildInputs;
            
            # Skip tests during build (run separately)
            doCheck = false;
            
            meta = with lib; {
              description = "Go SDK for XMTP messaging";
              homepage = "https://github.com/xmtp/go-sdk";
              license = licenses.mit;
              maintainers = [ ];
            };
          };
          
          # Build the Rust FFI library
          xmtp-ffi = pkgs.rustPlatform.buildRustPackage {
            pname = "xmtp-ffi";
            version = "0.1.0";
            src = ./ffi;
            
            cargoLock = {
              lockFile = ./ffi/Cargo.lock;
            };
            
            inherit buildInputs nativeBuildInputs;
            
            postInstall = ''
              # Copy the shared library to a standard location
              mkdir -p $out/lib
              cp target/*/release/libxmtp_ffi.* $out/lib/ || true
              cp target/release/libxmtp_ffi.* $out/lib/ || true
            '';
            
            meta = with lib; {
              description = "C FFI bindings for libxmtp";
            };
          };
        };
        
        devShells = {
          default = pkgs.mkShell {
            inputsFrom = [ self.packages.${system}.xmtp-go-sdk ];
            
            buildInputs = buildInputs ++ [
              # Go
              pkgs.go
              pkgs.gotools
              pkgs.gopls
              pkgs.go-outline
              
              # Rust
              rustToolchain
              pkgs.cargo-watch
              pkgs.cargo-edit
              
              # Build tools
              pkgs.pkg-config
              pkgs.cmake
              pkgs.mold
              
              # Development tools
              pkgs.gdb
              pkgs.just
            ] ++ nativeBuildInputs;
            
            shellHook = ''
              echo "XMTP Go SDK Development Environment"
              echo "===================================="
              echo ""
              echo "Go version: $(go version)"
              echo "Rust version: $(rustc --version)"
              echo ""
              echo "Commands:"
              echo "  make build      - Build the Go SDK"
              echo "  make test       - Run tests"
              echo "  make ffi        - Build the Rust FFI library"
              echo "  make example    - Build the example"
              echo ""
              
              # Set library path for finding the FFI library
              export LD_LIBRARY_PATH="$PWD:$LD_LIBRARY_PATH"
            '';
          };
        };
        
        checks = {
          # Go tests
          test = pkgs.buildGo124Module {
            pname = "xmtp-go-sdk-test";
            version = "0.1.0";
            src = ./.;
            
            vendorHash = "sha256-7SrehajxKazbFz7m9YbslCOrI3U+NDEzxWroo5Jy8VU=";
            modRoot = ".";
            
            inherit buildInputs nativeBuildInputs;
            
            doCheck = true;
            checkPhase = ''
              runHook preCheck
              go test -v ./...
              runHook postCheck
            '';
          };
        };
      }
    );
}
