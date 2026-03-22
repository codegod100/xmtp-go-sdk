{
  description = "XMTP Go SDK - Go bindings for XMTP messaging using PureGo";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };

        lib = nixpkgs.lib;

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
          rev = "2f6afaa729b11c8c2a25d7d470f6d7f852e78553";
          hash = "sha256-kz4dzh6Xm5k6yHjhUUp/k6chv2QGuSAc9h8CZmnROWQ=";
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

          # Build mobile bindings with UniFFI and generate C headers
          xmtp-ffi = pkgs.stdenv.mkDerivation {
            pname = "xmtp-ffi";
            version = "0.1.0";

            src = libxmtpSrc;

            nativeBuildInputs = nativeBuildInputs ++ [ 
              pkgs.rustup
              pkgs.stdenv.cc
              pkgs.git
              pkgs.cacert
              pkgs.perl
            ];
            inherit buildInputs;

            # Relax sandbox so cargo can access network
            __noChroot = true;

            buildPhase = ''
              export CARGO_HOME=$TMPDIR/cargo-home
              export RUSTUP_HOME=$TMPDIR/rustup-home
              export CARGO_NET_GIT_FETCH_WITH_CLI=true
              export SSL_CERT_FILE=${pkgs.cacert}/etc/ssl/certs/ca-bundle.crt
              export OPENSSL_NO_VENDOR=1
              export OPENSSL_LIB_DIR=${pkgs.openssl.out}/lib
              export OPENSSL_INCLUDE_DIR=${pkgs.openssl.dev}/include
              export RUSTFLAGS="-C link-arg=-fuse-ld=gold"
              
              rustup default stable
              
              # Build the mobile crate (produces libxmtpv3.so)
              cargo build -p xmtpv3 --release
              
              # Generate Swift headers - module.modulemap has C declarations
              cargo run --features "uniffi/cli" --bin ffi-uniffi-bindgen -- generate --library target/release/libxmtpv3.so --language swift --out-dir $TMPDIR/swift_headers
              
              # Copy outputs
              mkdir -p $out/lib $out/include
              cp target/release/libxmtpv3.so $out/lib/
              cp target/release/libxmtpv3.a $out/lib/ 2>/dev/null || true
              
              # Copy modulemap if it exists (C-compatible declarations)
              cp $TMPDIR/swift_headers/*.h $out/include/ 2>/dev/null || true
              cp $TMPDIR/swift_headers/module.modulemap $out/include/ 2>/dev/null || true
            '';

            installPhase = ''
              # Already done in buildPhase
            '';

            meta.description = "XMTP mobile bindings with C headers via UniFFI";
          };

        };

        devShells.default = pkgs.mkShell {
          inputsFrom = [ self.packages.${system}.xmtp-go-sdk ];

          buildInputs = buildInputs ++ [
            pkgs.go
            pkgs.gopls
            pkgs.gotools
            pkgs.rustup
            pkgs.cargo-watch
          ] ++ nativeBuildInputs;

          shellHook = ''
            echo "XMTP Go SDK Development"
            echo "  Go: $(go version)"
            if command -v rustc &> /dev/null; then
              echo "  Rust: $(rustc --version)"
            else
              echo "  Rust: run 'rustup default stable' to install"
            fi
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
