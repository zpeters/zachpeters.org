{
  description = "zachpeters.org website";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    { nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        # nixpkgs has go 1.27.1; 1.27.2 has security fixes. Remove this
        # override when nixpkgs catches up.
        go = pkgs.go_1_27.overrideAttrs (_: rec {
          version = "1.27.2";
          src = pkgs.fetchurl {
            url = "https://go.dev/dl/go${version}.src.tar.gz";
            hash = "sha256-A0ldorpkiU1A9cSZLklFT6eLUGkGBP+Stq//UIG3bmI=";
          };
        });
      in
      {
        devShells.default = pkgs.mkShell {
          packages = with pkgs; [
            entr
            fd
            go
            golangci-lint
            govulncheck
            just
            podman
            mitmproxy
            tmux
          ];
        };
      }
    );
}
