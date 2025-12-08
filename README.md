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
bash install.sh

bash uninstall.sh
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
