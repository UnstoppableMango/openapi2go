{
  buildGoApplication,
  fetchurl,
  lib,
  ginkgo,
  git,
  version,
}:
let
  petstoreSpec = fetchurl {
    url = "https://petstore3.swagger.io/api/v3/openapi.json";
    hash = "sha256-AEQcBa3WDyjaVetFY9P7a72jZLqOt7OB4uLJhhMAXII=";
  };
in
buildGoApplication {
  pname = "openapi2go";
  inherit version;

  src = lib.cleanSource ../.;
  modules = ./gomod2nix.toml;

  nativeCheckInputs = [
    ginkgo
    git
  ];

  checkPhase = ''
    mkdir bin && cp ${petstoreSpec} bin/petstore.json
    ginkgo run -r
  '';
}
