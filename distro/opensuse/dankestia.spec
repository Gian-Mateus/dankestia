# Spec for DANKESTIA for OpenSUSE/OBS

%global debug_package %{nil}

Name:           dankestia
Version:        1.2.3
Release:        1%{?dist}
Summary:        Dankestia - Material 3 inspired shell for Wayland compositors

License:        MIT
URL:            https://github.com/AvengeMedia/Dankestia
Source0:        dankestia-source.tar.gz
Source1:        dankestia-distropkg-amd64.gz
Source2:        dankestia-distropkg-arm64.gz

BuildRequires:  gzip
BuildRequires:  systemd-rpm-macros

# Core requirements
Requires:       (quickshell or quickshell-git)
Requires:       accountsservice
Requires:       dgop

# Core utilities (Highly recommended for DANKESTIA functionality)
Recommends:     cava
Recommends:     danksearch
Recommends:     matugen
Recommends:     NetworkManager
Recommends:     qt6-qtmultimedia
Suggests:       cups-pk-helper
Suggests:       qt6ct

%description
Dankestia (DANKESTIA) is a modern Wayland desktop shell built with Quickshell
and optimized for niri, Hyprland, Sway, and other wlroots compositors. Features
notifications, app launcher, wallpaper customization, and plugin system.

Includes auto-theming for GTK/Qt apps with matugen, 20+ customizable widgets,
process monitoring, notification center, clipboard history, dock, control center,
lock screen, and comprehensive plugin system.

%prep
%setup -q -n Dankestia-%{version}

%ifarch x86_64
gunzip -c %{SOURCE1} > dankestia
%endif
%ifarch aarch64
gunzip -c %{SOURCE2} > dankestia
%endif
chmod +x dankestia

%build

%install
install -Dm755 dankestia %{buildroot}%{_bindir}/dankestia

install -d %{buildroot}%{_datadir}/bash-completion/completions
install -d %{buildroot}%{_datadir}/zsh/site-functions
install -d %{buildroot}%{_datadir}/fish/vendor_completions.d
./dankestia completion bash > %{buildroot}%{_datadir}/bash-completion/completions/dankestia || :
./dankestia completion zsh > %{buildroot}%{_datadir}/zsh/site-functions/_dankestia || :
./dankestia completion fish > %{buildroot}%{_datadir}/fish/vendor_completions.d/dankestia.fish || :

install -Dm644 assets/systemd/dankestia.service %{buildroot}%{_userunitdir}/dankestia.service

install -Dm644 assets/dankestia-open.desktop %{buildroot}%{_datadir}/applications/dankestia-open.desktop
install -Dm644 assets/danklogo.svg %{buildroot}%{_datadir}/icons/hicolor/scalable/apps/danklogo.svg

install -dm755 %{buildroot}%{_datadir}/quickshell/dankestia
cp -r quickshell/* %{buildroot}%{_datadir}/quickshell/dankestia/

rm -rf %{buildroot}%{_datadir}/quickshell/dankestia/.git*
rm -f %{buildroot}%{_datadir}/quickshell/dankestia/.gitignore
rm -rf %{buildroot}%{_datadir}/quickshell/dankestia/.github
rm -rf %{buildroot}%{_datadir}/quickshell/dankestia/distro
rm -rf %{buildroot}%{_datadir}/quickshell/dankestia/core

echo "%{version}" > %{buildroot}%{_datadir}/quickshell/dankestia/VERSION

%posttrans
# Signal running DANKESTIA instances to reload
pkill -USR1 -x dankestia >/dev/null 2>&1 || :

%files
%license LICENSE
%doc CONTRIBUTING.md
%doc quickshell/README.md
%{_bindir}/dankestia
%dir %{_datadir}/fish
%dir %{_datadir}/fish/vendor_completions.d
%{_datadir}/fish/vendor_completions.d/dankestia.fish
%dir %{_datadir}/zsh
%dir %{_datadir}/zsh/site-functions
%{_datadir}/zsh/site-functions/_dankestia
%{_datadir}/bash-completion/completions/dankestia
%dir %{_datadir}/quickshell
%{_datadir}/quickshell/dankestia/
%{_userunitdir}/dankestia.service
%{_datadir}/applications/dankestia-open.desktop
%dir %{_datadir}/icons/hicolor
%dir %{_datadir}/icons/hicolor/scalable
%dir %{_datadir}/icons/hicolor/scalable/apps
%{_datadir}/icons/hicolor/scalable/apps/danklogo.svg

%changelog
* Mon Dec 16 2025 AvengeMedia <maintainer@avengemedia.com> - 1.0.3-1
- Update to stable v1.0.3 release

* Fri Dec 12 2025 AvengeMedia <maintainer@avengemedia.com> - 1.0.2-1
- Update to stable v1.0.2 release
- Bug fixes and improvements

* Fri Nov 22 2025 AvengeMedia <maintainer@avengemedia.com> - 0.6.2-1
- Stable release build with pre-built binaries
- Multi-arch support (x86_64, aarch64)
