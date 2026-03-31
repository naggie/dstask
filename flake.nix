{
  description = "dstask - single binary terminal-based TODO manager";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    flake-compat = {
      url = "github:edolstra/flake-compat";
      flake = false;
    };
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    let
      version = "1.0.1";
    in
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in {
        packages.default = pkgs.buildGoModule {
          pname = "dstask";
          inherit version;
          src = ./.;

          # To update after changing go.mod/go.sum:
          #   1. Set vendorHash = pkgs.lib.fakeHash;
          #   2. Run `nix build` and let it fail
          #   3. Copy the "got" hash from the error into vendorHash
          vendorHash = "sha256-HSqAbxkkjuMulFymeqApWr/JZ+a7OUTu5EYLGPL/j2U=";

          subPackages = [ "cmd/dstask" "cmd/dstask-import" ];

          ldflags = [
            "-s" "-w"
            "-X github.com/naggie/dstask.GIT_COMMIT=${self.shortRev or "dirty"}"
            "-X github.com/naggie/dstask.VERSION=${version}"
          ];

          meta = with pkgs.lib; {
            description = "Single binary terminal-based TODO manager";
            homepage = "https://github.com/naggie/dstask";
            license = licenses.mit;
            mainProgram = "dstask";
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            golangci-lint
          ];
        };
      }
    );
}
