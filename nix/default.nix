{
  buildGoApplication,
  lib,
  ginkgo,
  git,
  version,
}:
buildGoApplication {
  pname = "";
  inherit version;

  src = lib.cleanSource ../.;
  modules = ../gomod2nix.toml;

  nativeCheckInputs = [
    ginkgo
    git
  ];

  checkPhase = ''
    ginkgo run -r --label-filter='!E2E'
  '';
}
