{ pkgs ? import <nixpkgs> {}}:

pkgs.mkShell {
  packages = [ 
    pkgs.entr
    pkgs.fd
    pkgs.go
    pkgs.golangci-lint
    pkgs.just 
    pkgs.podman
  ];
}
