{
  self,
  pkgs,
  ...
}:
pkgs.testers.runNixOSTest {
  name = "dankestia-nixos-module";

  nodes.machine = {
    imports = [
      self.nixosModules.dank-material-shell
    ];

    users.users.danklinux = {
      isNormalUser = true;
      extraGroups = [ "wheel" ];
    };

    programs.dank-material-shell = {
      enable = true;
      systemd.enable = true;
      plugins = {
        TestPlugin = {
          src = pkgs.emptyDirectory;
        };
      };
    };

    system.stateVersion = "25.11";
  };

  testScript = ''
    import json

    machine.wait_for_unit("multi-user.target")

    machine.succeed("command -v dankestia")
    machine.succeed("command -v quickshell")
    machine.succeed("su -- danklinux -c 'dankestia --help >/dev/null'")
    machine.succeed("test -d /etc/xdg/quickshell/dankestia-plugins")
    machine.succeed("test -f /run/current-system/sw/lib/systemd/user/dankestia.service")

    payload = json.loads(machine.succeed("su -- danklinux -c 'dankestia doctor --json'"))
    t.assertIn("summary", payload)
    t.assertIsInstance(payload.get("results"), list)
  '';
}
