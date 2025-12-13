# kde AccentColor for logitech g203 and keychron c2 pro

Set Logitech G203 mouse and keychron c2 pro color to KDE AccentColor

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
cd de-AccentColor-logitech-g203-keychron-c2-pro-main 

# install or uninstall (systemd userspace)
bash install.sh

bash uninstall.sh
```

### Udev permisions. Example of Keychron Keychron C2 Pro
```bash
lsusb
Bus 001 Device 003: ID 3434:0526 Keychron Keychron C2 Pro

sudo nano /etc/udev/rules.d/99-keychron.rules
SUBSYSTEM=="hidraw", ATTRS{idVendor}=="3434", ATTRS{idProduct}=="0526", MODE="0666"

sudo udevadm control --reload-rules
sudo udevadm trigger
ls -l /dev/hidraw*
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

### how to capture codes
```bash
sudo hexdump -Cv /dev/hidraw4 # use your device
```
