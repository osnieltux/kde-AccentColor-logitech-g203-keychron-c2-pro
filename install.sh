#!/bin/bash

go build -ldflags="-s -w" -o kde_color_logitechg203 main.go

cp kde_color_logitechg203 ~/.local/bin/kde_color_logitechg203 -f -n

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