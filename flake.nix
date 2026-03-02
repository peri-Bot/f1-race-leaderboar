{
  description = "F1 Race Leaderboard — Serverless AWS API";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
            go-tools           # staticcheck
            opentofu           # open-source Terraform fork (HCL compatible)
            awscli2
            zip
            gnumake
          ];

          shellHook = ''
            echo "🏎️  F1 Race Leaderboard dev shell"
            echo "   Go:       $(go version | cut -d' ' -f3)"
            echo "   OpenTofu: $(tofu version | head -1)"
            echo "   AWS CLI:  $(aws --version 2>&1 | cut -d' ' -f1)"
            alias terraform=tofu
          '';
        };
      }
    );
}
