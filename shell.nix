let

  pkgs = import (builtins.fetchTarball {
    name = "nixpkgs-alpha";
    url = https://github.com/nixos/nixpkgs/archive/34b37ad59c46fa273a03944c3f10e269f3984852.tar.gz;
    sha256 = "09hsmbpiycyffjq7k49g8qav4fwlj6nrp9nzxpk5rm1rws8h9z61";
  }) {};

  deps = [
    pkgs.capnproto
    pkgs.go
    pkgs.which
  ];

  env = pkgs.stdenv.mkDerivation {
    name = "capnpn-demo";
    buildInputs = deps;
  };

in env
