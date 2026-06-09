{
  description = "Dank Material Shell";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    flake-compat = {
      url = "github:NixOS/flake-compat";
      flake = false;
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      ...
    }:
    let
      goModVersion =
        let
          content = builtins.readFile ./core/go.mod;
          lines = builtins.filter builtins.isString (builtins.split "\n" content);
          goLines = builtins.filter (l: builtins.match "go [0-9]+\\..*" l != null) lines;
          matched =
            if goLines != [ ] then builtins.match "go ([0-9]+)\\.([0-9]+).*" (builtins.head goLines) else null;
        in
        if matched != null then
          {
            major = builtins.elemAt matched 0;
            minor = builtins.elemAt matched 1;
          }
        else
          {
            major = "1";
            minor = "25";
          };
      goForPkgs = pkgs: pkgs.${"go_${goModVersion.major}_${goModVersion.minor}"};
      forEachSystem =
        fn:
        nixpkgs.lib.genAttrs [ "aarch64-darwin" "aarch64-linux" "x86_64-darwin" "x86_64-linux" ] (
          system: fn system nixpkgs.legacyPackages.${system}
        );
      forEachLinuxSystem =
        fn:
        nixpkgs.lib.genAttrs [ "aarch64-linux" "x86_64-linux" ] (
          system: fn system nixpkgs.legacyPackages.${system}
        );

      mkModuleWithDmsPkgs =
        modulePath:
        args@{ pkgs, ... }:
        {
          imports = [
            (import modulePath (args // { dankestiaPkgs = buildDmsPkgs pkgs; }))
          ];
        };

      mkQmlImportPath =
        pkgs: qmlPkgs:
        pkgs.lib.concatStringsSep ":" (map (o: "${o}/${pkgs.qt6.qtbase.qtQmlPrefix}") qmlPkgs);

      mkQtPluginPath =
        pkgs: qtPkgs:
        pkgs.lib.concatStringsSep ":" (map (o: "${o}/${pkgs.qt6.qtbase.qtPluginPrefix}") qtPkgs);

      qmlPkgs =
        pkgs: with pkgs.kdePackages; [
          kirigami.unwrapped
          sonnet
          qtmultimedia
          qtimageformats
          kimageformats
        ];

      # Allows downstream modules to provide their own 'pkgs' (with overlays)
      # instead of being forced to use the flake's locked nixpkgs.
      mkDmsShell =
        pkgs:
        let
          mkDate =
            longDate:
            pkgs.lib.concatStringsSep "-" [
              (builtins.substring 0 4 longDate)
              (builtins.substring 4 2 longDate)
              (builtins.substring 6 2 longDate)
            ];
          version =
            let
              rawVersion = pkgs.lib.removePrefix "v" (pkgs.lib.trim (builtins.readFile ./quickshell/VERSION));
              cleanVersion = builtins.replaceStrings [ " " ] [ "" ] rawVersion;
              dateSuffix = "+date=" + mkDate (self.lastModifiedDate or "19700101");
              revSuffix = "_" + (self.shortRev or "dirty");
            in
            "${cleanVersion}${dateSuffix}${revSuffix}";
        in
        pkgs.lib.makeOverridable (
          {
            extraQtPackages ? [ ],
          }:
          (pkgs.buildGoModule.override { go = goForPkgs pkgs; }) (
            let
              rootSrc = ./.;
              qtPackages = (qmlPkgs pkgs) ++ extraQtPackages;
            in
            {
              inherit version;
              pname = "dankestia-shell";
              src = ./core;
              vendorHash = "sha256-nvxFHQhOfBGl3h51fgYDb39K0NCj+H8mAEyKr1qOwJQ=";

              subPackages = [ "cmd/dankestia" ];

              ldflags = [
                "-s"
                "-w"
                "-X 'main.Version=${version}'"
              ];

              nativeBuildInputs = with pkgs; [
                installShellFiles
                makeWrapper
              ];

              postInstall = ''
                mkdir -p $out/share/quickshell/dankestia
                cp -r ${rootSrc}/quickshell/. $out/share/quickshell/dankestia/

                chmod u+w $out/share/quickshell/dankestia/VERSION
                echo "${version}" > $out/share/quickshell/dankestia/VERSION

                # Install desktop file and icon
                install -D ${rootSrc}/assets/dankestia-open.desktop \
                  $out/share/applications/dankestia-open.desktop
                install -D ${rootSrc}/core/assets/danklogo.svg \
                  $out/share/hicolor/scalable/apps/danklogo.svg

                wrapProgram $out/bin/dankestia \
                  --add-flags "-c $out/share/quickshell/dankestia" \
                  --prefix "NIXPKGS_QT6_QML_IMPORT_PATH" ":" "${mkQmlImportPath pkgs qtPackages}" \
                  --prefix "QT_PLUGIN_PATH" ":" "${mkQtPluginPath pkgs qtPackages}"

                install -Dm644 ${rootSrc}/assets/systemd/dankestia.service \
                  $out/lib/systemd/user/dankestia.service

                substituteInPlace $out/lib/systemd/user/dankestia.service \
                  --replace-fail /usr/bin/dankestia $out/bin/dankestia \
                  --replace-fail /usr/bin/pkill ${pkgs.procps}/bin/pkill

                substituteInPlace $out/share/quickshell/dankestia/Modules/Greetd/assets/dankestia-greeter \
                  --replace-fail /bin/bash ${pkgs.bashInteractive}/bin/bash

                substituteInPlace $out/share/quickshell/dankestia/assets/pam/fprint \
                  --replace-fail pam_fprintd.so ${pkgs.fprintd}/lib/security/pam_fprintd.so

                substituteInPlace $out/share/quickshell/dankestia/assets/pam/u2f \
                  --replace-fail pam_u2f.so ${pkgs.pam_u2f}/lib/security/pam_u2f.so

                installShellCompletion --cmd dankestia \
                  --bash <($out/bin/dankestia completion bash) \
                  --fish <($out/bin/dankestia completion fish) \
                  --zsh <($out/bin/dankestia completion zsh)
              '';

              meta = {
                description = "Desktop shell for wayland compositors built with Quickshell & GO";
                homepage = "https://danklinux.com";
                changelog = "https://github.com/AvengeMedia/Dankestia/releases/tag/v${version}";
                license = pkgs.lib.licenses.mit;
                mainProgram = "dankestia";
                platforms = pkgs.lib.platforms.linux;
              };
            }
          )
        ) { };

      buildDmsPkgs = pkgs: {
        dankestia-shell = mkDmsShell pkgs;
      };
    in
    {
      packages = forEachSystem (
        system: pkgs: {
          dankestia-shell = mkDmsShell pkgs;
          default = self.packages.${system}.dankestia-shell;
          quickshell = builtins.warn "dank-material-shell: the package Quickshell is not included in the DANKESTIA flake anymore. We recommend you to use the one from nixos-unstable branch of Nixpkgs or the upstream flake." pkgs.quickshell;
        }
      );

      lib = { inherit mkDmsShell buildDmsPkgs; };

      homeModules.dank-material-shell = mkModuleWithDmsPkgs ./distro/nix/home.nix;

      homeModules.default = self.homeModules.dank-material-shell;

      homeModules.niri = import ./distro/nix/niri.nix;

      homeModules.dankMaterialShell.default = builtins.warn "dank-material-shell: flake output `homeModules.dankMaterialShell.default` has been renamed to `homeModules.dank-material-shell`" self.homeModules.dank-material-shell;

      homeModules.dankMaterialShell.niri = builtins.warn "dank-material-shell: flake output `homeModules.dankMaterialShell.niri` has been renamed to `homeModules.niri`" self.homeModules.niri;

      nixosModules.dank-material-shell = mkModuleWithDmsPkgs ./distro/nix/nixos.nix;

      nixosModules.default = self.nixosModules.dank-material-shell;

      nixosModules.greeter = mkModuleWithDmsPkgs ./distro/nix/greeter.nix;

      nixosModules.dankMaterialShell = builtins.warn "dank-material-shell: flake output `nixosModules.dankMaterialShell` has been renamed to `nixosModules.dank-material-shell`" self.nixosModules.dank-material-shell;

      devShells = forEachSystem (
        system: pkgs:
        let
          devQmlPkgs = with pkgs;
          [
            quickshell
            kdePackages.qtdeclarative
          ]
          ++ (qmlPkgs pkgs);
        in
        {
          default = pkgs.mkShell {
            buildInputs =
              with pkgs;
              [
                (goForPkgs pkgs)
                go-mockery_2
                gopls
                delve
                go-tools
                gnumake

                prek
                uv # for prek
                shellcheck

                # Nix development tools
                nixd
                nil
              ]
              ++ devQmlPkgs;

            shellHook = ''
              touch quickshell/.qmlls.ini 2>/dev/null
              if [ ! -f .git/hooks/pre-commit ]; then prek install; fi
            '';

            QML2_IMPORT_PATH = mkQmlImportPath pkgs devQmlPkgs;
            QT_PLUGIN_PATH = mkQtPluginPath pkgs devQmlPkgs;
          };
        }
      );

      nixosTests = forEachLinuxSystem (
        system: pkgs:
        import ./distro/nix/tests {
          inherit
            self
            pkgs
            ;
          lib = pkgs.lib;
        }
      );
    };
}
