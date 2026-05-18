{
  description = "Open API Specs into Go Types";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    systems.url = "github:nix-systems/default";

    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };

    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    gomod2nix = {
      url = "github:nix-community/gomod2nix";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.flake-utils.inputs.systems.follows = "systems";
    };

    ux = {
      url = "github:UnstoppableMango/ux?ref=fancy-fresh";
      inputs.nixpkgs.follows = "nixpkgs";
      inputs.systems.follows = "systems";
      inputs.flake-parts.follows = "flake-parts";
      inputs.gomod2nix.follows = "gomod2nix";
      inputs.treefmt-nix.follows = "treefmt-nix";
    };
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = import inputs.systems;

      imports = with inputs; [
        treefmt-nix.flakeModule
        ux.flakeModules.default
      ];

      ux = {
        builders.openapi2go = ./nix/builder.nix;
        gen.test =
          { lib }:
          {
            builder = "openapi2go";
            config = {
              name = "petstore";
              spec = lib.fetchurl "https://petstore3.swagger.io/api/v3/openapi.json";
            };
          };
      };

      perSystem =
        {
          self',
          pkgs,
          system,
          ...
        }:
        let
          version = "0.0.1";
          openapi2go = pkgs.callPackage ./nix { inherit version; };
        in
        {
          _module.args.pkgs = import inputs.nixpkgs {
            inherit system;
            overlays = with inputs; [
              gomod2nix.overlays.default
            ];
          };

          packages = {
            inherit openapi2go;
            default = openapi2go;
          };

          apps.default = {
            type = "app";
            program = "${self'.packages.openapi2go}/bin/openapi2go";
            meta.description = "Convert OpenAPI Specs into Go Types";
          };

          devShells.default = pkgs.mkShellNoCC {
            packages = with pkgs; [
              curl
              direnv
              dprint
              go
              gomod2nix
              gopls
              golangci-lint
              ginkgo
              gnumake
              gopls
              jq
              nixfmt
            ];

            CURL = "${pkgs.curl}/bin/curl";
            DPRINT = "${pkgs.dprint}/bin/dprint";
            GO = "${pkgs.go}/bin/go";
            GOMOD2NIX = "${pkgs.gomod2nix}/bin/gomod2nix";
            GOLANGCI_LINT = "${pkgs.golangci-lint}/bin/golangci-lint";
            GINKGO = "${pkgs.ginkgo}/bin/ginkgo";
            JQ = "${pkgs.jq}/bin/jq";
          };

          treefmt.programs = {
            nixfmt.enable = true;
            gofmt.enable = true;
          };
        };
    };
}
