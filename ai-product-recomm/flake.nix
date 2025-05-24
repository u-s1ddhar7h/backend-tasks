{
  description = "Python Flake - basic";

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
        python313
        python313Packages.flask
        python313Packages.pandas
        python313Packages.scikit-learn
        python313Packages.joblib
      	python313Packages.pip
      ];
      shellHook = ''
        echo "Python environment is ready!"
      '';
    };
  };
}
