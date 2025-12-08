# kde AccentColor logitech g203

Set Logitech G203 mouse color to KDE AccentColor

## Dependencies

- [kreadconfig6](https://community.kde.org/Frameworks)
- [libratbag](https://github.com/libratbag/libratbag)
- [golang](https://go.dev/) (just to compile)
- [dbus](https://www.freedesktop.org/wiki/Software/dbus/)

## Install dependencies (Manjaro):

```bash
pamac install libratbag kconfig go dbus
```

## Install/Uninstall systemd service for current user.

```bash
# download
wget https://github.com/osnieltux/kde-AccentColor-logitech-g203/archive/refs/heads/main.zip
unzip main.zip 
cd kde-AccentColor-logitech-g203-main

# install or uninstall (systemd userspace)
bash install.sh

bash uninstall.sh
```

### Install without compiling with Go, using the binary release.
```bash
bash install_prebuild.sh
```

## Check

```bash
systemctl --user status kde_color_logitechg203.service
```

### TODO

- Implement config.toml
- Add support for list of devices
- Indicate LED color per device instead of 0 (current default)
- use "github.com/gotmc/libusb/v2" instead "libratbag" (i need a lot of time)
