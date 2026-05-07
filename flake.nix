{
  description = "Development environment for vipsgen";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };

        vipsPackages = [
          pkgs.vips
          pkgs.glib
          pkgs.gobject-introspection
        ];
      in
      {
        devShells.default = pkgs.mkShell {
          packages = [
            pkgs.go_1_25
            pkgs.golangci-lint
            pkgs.gnumake
            pkgs.pkg-config
          ] ++ vipsPackages;

          shellHook = ''
            export CGO_CFLAGS_ALLOW="-Xpreprocessor"
            export GI_TYPELIB_PATH="${pkgs.lib.makeSearchPath "lib/girepository-1.0" vipsPackages}:''${GI_TYPELIB_PATH:-}"
            export XDG_DATA_DIRS="${pkgs.lib.makeSearchPath "share" vipsPackages}:''${XDG_DATA_DIRS:-}"
          '';
        };

        formatter = pkgs.nixfmt-rfc-style;
      }
    );
}
