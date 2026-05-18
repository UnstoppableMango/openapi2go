{
  name,
  openapi2go,
  spec,
  writeShellApplication,
}:
writeShellApplication {
  name = "openapi2go";

  runtimeInputs = [
    openapi2go
  ];

  # TODO: --config ${config}
  text = ''
    openapi2go '${spec}' \
      --package-name '${name}' \
      --output "$out"
  '';
}
