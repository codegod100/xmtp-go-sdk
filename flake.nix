{
  description = "XMTP Go SDK - Go bindings for XMTP messaging using PureGo";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    rust-overlay.url = "github:oxalica/rust-overlay";
    crane.url = "github:ipetkov/crane";
  };

  outputs = { self, nixpkgs, flake-utils, rust-overlay, crane, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        overlays = [ (import rust-overlay) ];
        pkgs = import nixpkgs {
          inherit system overlays;
        };

        lib = nixpkgs.lib;

        # Rust toolchain
        rustToolchain = pkgs.rust-bin.stable.latest.default;

        # Crane library
        craneLib = crane.mkLib pkgs;

        # Build inputs
        buildInputs = with pkgs; [
          openssl
        ] ++ lib.optionals stdenv.isLinux [
          systemd
        ] ++ lib.optionals stdenv.isDarwin [
          darwin.apple_sdk.frameworks.Security
          darwin.apple_sdk.frameworks.CoreFoundation
          darwin.apple_sdk.frameworks.SystemConfiguration
        ];

        nativeBuildInputs = with pkgs; [
          pkg-config
        ];

        # Fetch libxmtp source
        libxmtpSrc = pkgs.fetchFromGitHub {
          owner = "xmtp";
          repo = "libxmtp";
          rev = "main";
          hash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="; # Will be fixed
        };

        # Common args for crane builds
        commonArgs = {
          inherit buildInputs nativeBuildInputs;
        };

      in {
        packages = {
          default = self.packages.${system}.xmtp-go-sdk;

          # Go SDK package
          xmtp-go-sdk = pkgs.buildGo124Module {
            pname = "xmtp-go-sdk";
            version = "0.1.0";
            src = ./.;

            vendorHash = "sha256-7SrehajxKazbFz7m9YbslCOrI3U+NDEzxWroo5Jy8VU=";
            modRoot = ".";

            inherit buildInputs nativeBuildInputs;
            doCheck = false;

            meta = with lib; {
              description = "Go SDK for XMTP messaging";
              homepage = "https://github.com/xmtp/go-sdk";
              license = licenses.mit;
            };
          };

          # Build stub FFI library (no libxmtp dependency, for testing)
          xmtp-ffi-stub = craneLib.buildPackage (commonArgs // {
            pname = "xmtp-ffi-stub";
            version = "0.1.0";
            src = ./ffi;

            cargoToml = ./ffi/Cargo.toml;
            cargoLock = ./ffi/Cargo.lock;

            postInstall = ''
              mkdir -p $out/lib
              cp target/release/libxmtp_ffi.so $out/lib/ 2>/dev/null || true
              cp target/release/libxmtp_ffi.a $out/lib/ 2>/dev/null || true
            '';

            meta.description = "C FFI stubs for libxmtp";
          });

          # Build FFI with real libxmtp (workspace build)
          xmtp-ffi = pkgs.stdenv.mkDerivation {
            pname = "xmtp-ffi";
            version = "0.1.0";

            # Use libxmtp source as base
            src = pkgs.fetchFromGitHub {
              owner = "xmtp";
              repo = "libxmtp";
              rev = "2f6afaa729b11c8c2a25d7d470f6d7f852e78553";
              hash = "sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";
            };

            nativeBuildInputs = nativeBuildInputs ++ [ rustToolchain ];
            buildInputs = buildInputs;

            # Copy our FFI crate into the workspace
            preBuild = ''
              # Add our FFI crate to the workspace
              cp -r ${./ffi} bindings/xmtp-ffi
              chmod -R u+w bindings/xmtp-ffi

              # Update workspace members
              sed -i 's/members = \[/members = ["bindings\/xmtp-ffi",/' Cargo.toml
            '';

            buildPhase = ''
              runHook preBuild
              cargo build -p xmtp-ffi --release
              runHook postBuild
            '';

            installPhase = ''
              runHook preInstall
              mkdir -p $out/lib
              cp target/release/libxmtp_ffi.so $out/lib/
              cp target/release/libxmtp_ffi.a $out/lib/
              runHook postInstall
            '';

            meta.description = "C FFI bindings for libxmtp with full implementation";
          };
        };

        devShells.default = pkgs.mkShell {
          inputsFrom = [ self.packages.${system}.xmtp-go-sdk ];

          buildInputs = buildInputs ++ [
            pkgs.go
            pkgs.gopls
            pkgs.gotools
            rustToolchain
            pkgs.cargo-watch
          ] ++ nativeBuildInputs;

          shellHook = ''
            echo "XMTP Go SDK Development"
            echo "  Go: $(go version)"
            echo "  Rust: $(rustc --version)"
            echo ""
            echo "Commands: make build | make test | make ffi"
            export LD_LIBRARY_PATH="$PWD/result/lib:$LD_LIBRARY_PATH"
          '';
        };

        checks.test = pkgs.buildGo124Module {
          pname = "xmtp-go-sdk-test";
          version = "0.1.0";
          src = ./.;
          vendorHash = "sha256-7SrehajxKazbFz7m9YbslCOrI3U+NDEzxWroo5Jy8VU=";
          modRoot = ".";
          inherit buildInputs nativeBuildInputs;
          doCheck = true;
          checkPhase = "go test -v ./...";
        };
      }
    );
}
