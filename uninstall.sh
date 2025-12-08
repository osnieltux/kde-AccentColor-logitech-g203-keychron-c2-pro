#!/bin/bash

systemctl --user stop kde_color_logitechg203.service
rm ~/.config/systemd/user/kde_color_logitechg203.service
rm ~/.local/bin/kde_color_logitechg203

systemctl --user daemon-reload