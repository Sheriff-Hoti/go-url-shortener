{
  pkgs ? import <nixpkgs> { },
}:

pkgs.mkShell {
  buildInputs = [
    pkgs.templ
    pkgs.tailwindcss_4
    pkgs.air
    pkgs.sqlite
    pkgs.sqlc
  ];
}
