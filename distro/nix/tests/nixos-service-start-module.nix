{
  self,
  pkgs,
  ...
}:
let
  fakeDms = pkgs.writeShellScriptBin "dankestia" ''
    printf '%s\n' "$@" > /tmp/dankestia-service-args
    exec ${pkgs.coreutils}/bin/sleep 300
  '';
in
pkgs.testers.runNixOSTest {
  name = "dankestia-nixos-service-start-module";

  nodes.machine = {
    imports = [
      self.nixosModules.dank-material-shell
    ];

    users.users.danklinux = {
      isNormalUser = true;
      linger = true;
      extraGroups = [ "wheel" ];
    };

    programs.dank-material-shell = {
      enable = true;
      package = fakeDms;
      systemd = {
        enable = true;
        target = "default.target";
      };
    };

    system.stateVersion = "25.11";
  };

  testScript = ''
    machine.wait_for_unit("multi-user.target")
    machine.wait_for_unit("user@1000.service")

    machine.succeed("systemctl --machine=danklinux@ --user start dankestia.service")
    machine.wait_until_succeeds("systemctl --machine=danklinux@ --user is-active dankestia.service")
    machine.wait_until_succeeds("test -f /tmp/dankestia-service-args")
    machine.succeed("grep -Fx run /tmp/dankestia-service-args")
    machine.succeed("grep -Fx -- --session /tmp/dankestia-service-args")
  '';
}
