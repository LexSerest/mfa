{
  description = "Secure 2FA CLI manager";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-25.11";
    utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, utils }:
    utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
        version = "1.0.0";
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "mfa";
          inherit version;

          src = ./.;
          vendorHash = "sha256-vWauVUXY+th7KR9cKMscyYiiLVXa2dBXkKFcTIN8ak4=";
          ldflags = [ "-s" "-w" ];
          preBuild = ''
            make completions
          '';

          postInstall = ''
            installShellCompletion --zsh completions/mfa.zsh
            installShellCompletion --fish completions/mfa.fish
            installShellCompletion --bash completions/mfa.bash
          '';

          nativeBuildInputs = [ pkgs.installShellFiles ];

          meta = with pkgs.lib; {
            description = "Secure 2FA CLI manager";
            homepage = "https://github.com/LexSerest/mfa";
            license = licenses.mit;
            maintainers = [ "LexSerest <lexserest@gmail.com>" ];
            mainProgram = "mfa";
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gnumake
            goreleaser
          ];
        };
      });
}
