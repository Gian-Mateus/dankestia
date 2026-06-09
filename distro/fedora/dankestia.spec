# Feodra spec for DANKESTIA stable releases

%global debug_package %{nil}
%global version VERSION_PLACEHOLDER
%global pkg_summary Dankestia - Material 3 inspired shell for Wayland compositors

Name:           dankestia
Version:        %{version}
Release:        RELEASE_PLACEHOLDER%{?dist}
Summary:        %{pkg_summary}

License:        MIT
URL:            https://github.com/AvengeMedia/Dankestia

Source0:        dankestia-qml.tar.gz

BuildRequires:  gzip
BuildRequires:  wget
BuildRequires:  systemd-rpm-macros

Requires:       (quickshell or quickshell-git)
Requires:       accountsservice
Requires:       dankestia-cli = %{version}-%{release}
Requires:       dgop

Recommends:     cava
Recommends:     danksearch
Recommends:     matugen
Recommends:     NetworkManager
Recommends:     qt6-qtmultimedia
Suggests:       cups-pk-helper
Suggests:       qt6ct

%description
Dankestia (DANKESTIA) is a modern Wayland desktop shell built with Quickshell
and optimized for the niri and hyprland compositors. Features notifications,
app launcher, wallpaper customization, and fully customizable with plugins.

Includes auto-theming for GTK/Qt apps with matugen, 20+ customizable widgets,
process monitoring, notification center, clipboard history, dock, control center,
lock screen, and comprehensive plugin system.

%package -n dankestia-cli
Summary:        Dankestia CLI tool
License:        MIT
URL:            https://github.com/AvengeMedia/Dankestia

%description -n dankestia-cli
Command-line interface for Dankestia configuration and management.
Provides native DBus bindings, NetworkManager integration, and system utilities.

%prep
%setup -q -c -n dankestia-qml

case "%{_arch}" in
  x86_64)
    ARCH_SUFFIX="amd64"
    ;;
  aarch64)
    ARCH_SUFFIX="arm64"
    ;;
  *)
    echo "Unsupported architecture: %{_arch}"
    exit 1
    ;;
esac

# Download dankestia-cli for target architecture
wget -O %{_builddir}/dankestia-cli.gz "https://github.com/AvengeMedia/Dankestia/releases/latest/download/dankestia-distropkg-${ARCH_SUFFIX}.gz" || {
  echo "Failed to download dankestia-cli for architecture %{_arch}"
  exit 1
}
gunzip -c %{_builddir}/dankestia-cli.gz > %{_builddir}/dankestia-cli
chmod +x %{_builddir}/dankestia-cli

%build

%install
install -Dm755 %{_builddir}/dankestia-cli %{buildroot}%{_bindir}/dankestia

# Shell completions
install -d %{buildroot}%{_datadir}/bash-completion/completions
install -d %{buildroot}%{_datadir}/zsh/site-functions
install -d %{buildroot}%{_datadir}/fish/vendor_completions.d
%{_builddir}/dankestia-cli completion bash > %{buildroot}%{_datadir}/bash-completion/completions/dankestia || :
%{_builddir}/dankestia-cli completion zsh > %{buildroot}%{_datadir}/zsh/site-functions/_dankestia || :
%{_builddir}/dankestia-cli completion fish > %{buildroot}%{_datadir}/fish/vendor_completions.d/dankestia.fish || :

install -Dm644 %{_builddir}/dankestia-qml/assets/systemd/dankestia.service %{buildroot}%{_userunitdir}/dankestia.service

install -Dm644 %{_builddir}/dankestia-qml/assets/dankestia-open.desktop %{buildroot}%{_datadir}/applications/dankestia-open.desktop
install -Dm644 %{_builddir}/dankestia-qml/assets/danklogo.svg %{buildroot}%{_datadir}/icons/hicolor/scalable/apps/danklogo.svg

install -dm755 %{buildroot}%{_datadir}/quickshell/dankestia
cp -r %{_builddir}/dankestia-qml/* %{buildroot}%{_datadir}/quickshell/dankestia/

rm -rf %{buildroot}%{_datadir}/quickshell/dankestia/.git*
rm -f %{buildroot}%{_datadir}/quickshell/dankestia/.gitignore
rm -rf %{buildroot}%{_datadir}/quickshell/dankestia/.github
rm -rf %{buildroot}%{_datadir}/quickshell/dankestia/distro

echo "%{version}" > %{buildroot}%{_datadir}/quickshell/dankestia/VERSION

%posttrans
# Signal running DANKESTIA instances to reload
pkill -USR1 -x dankestia >/dev/null 2>&1 || :

%files
%license LICENSE
%doc README.md CONTRIBUTING.md
%{_datadir}/quickshell/dankestia/
%{_userunitdir}/dankestia.service
%{_datadir}/applications/dankestia-open.desktop
%{_datadir}/icons/hicolor/scalable/apps/danklogo.svg

%files -n dankestia-cli
%{_bindir}/dankestia
%{_datadir}/bash-completion/completions/dankestia
%{_datadir}/zsh/site-functions/_dankestia
%{_datadir}/fish/vendor_completions.d/dankestia.fish

%changelog
* CHANGELOG_DATE_PLACEHOLDER AvengeMedia <contact@avengemedia.com> - VERSION_PLACEHOLDER-RELEASE_PLACEHOLDER
- Stable release VERSION_PLACEHOLDER
- Built from GitHub release
