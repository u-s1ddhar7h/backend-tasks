{
  description = "E-commerce Backend Development Environment (Go)";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = { self, nixpkgs }: 
  let
    system = "x86_64-linux";
    pkgs = nixpkgs.legacyPackages.${system};
  in {
    devShells.${system}.default = pkgs.mkShell {
      buildInputs = with pkgs; [
      	go
      	gopls
      ];
      shellHook = ''echo "Entering Go environment ..."'';
    };
  };
}
