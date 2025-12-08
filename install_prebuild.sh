#!/bin/bash


wget https://github.com/osnieltux/kde-AccentColor-logitech-g203/releases/download/release/kde_color_logitechg203 -O ~/.local/bin/kde_color_logitechg203
chmod +x ~/.local/bin/kde_color_logitechg203
mkdir -p ~/.config/systemd/user

cat << EOF > ~/.config/systemd/user/kde_color_logitechg203.service
[Unit]
Description=KDE Color for Logitech G203
After=default.target

[Service]
ExecStart=%h/.local/bin/kde_color_logitechg203
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=default.target
EOF

systemctl --user daemon-reload
systemctl --user enable kde_color_logitechg203.service
systemctl --user start kde_color_logitechg203.service
