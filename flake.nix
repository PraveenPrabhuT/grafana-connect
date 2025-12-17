{
  description = "Grafana Connect - Context-aware Grafana dashboard launcher";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    gomod2nix.url = "github:nix-community/gomod2nix";
    gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
    gomod2nix.inputs.flake-utils.follows = "flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, gomod2nix }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ gomod2nix.overlays.default ];
        };
      in
      {
        packages = rec {
          grafana-connect = pkgs.buildGoApplication {
            pname = "grafana-connect";
            version = self.shortRev or "dirty";

            src = ./.;

            # You must generate this file! See instructions below.
            modules = ./gomod2nix.toml;

            # Inject version info into the cmd package
            # Note: We use lowercase 'version' variable as defined in your cmd/version.go
            ldflags = [
              "-s"
              "-w"
              "-X github.com/PraveenPrabhuT/grafana-connect/cmd.version=${self.shortRev or "dirty"}"
              "-X github.com/PraveenPrabhuT/grafana-connect/cmd.commit=${self.shortRev or "dirty"}"
            ];

            meta = with pkgs.lib; {
              description = "Context-aware Grafana dashboard launcher";
              homepage = "https://github.com/PraveenPrabhuT/grafana-connect";
              license = licenses.mit;
              maintainers = [ ];
              platforms = platforms.unix;
              mainProgram = "grafana-connect";
            };
          };

          # Set the default package
          default = grafana-connect;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gotools
            gopls
            golangci-lint
            gomod2nix.packages.${system}.default
            # Useful for testing clipboard on Linux inside the shell
            xclip 
          ];
        };

        apps.default = {
          type = "app";
          program = "${self.packages.${system}.default}/bin/grafana-connect";
        };
      }
    );
}
