#!/bin/bash
set -eo pipefail

podman build --platform linux/amd64 -t podman:5000/zachpeters.org:latest .
podman push podman:5000/zachpeters.org:latest
ssh podman "podman auto-update"
ssh podman "systemctl stop zachpeters"
ssh podman "systemctl start zachpeters"
ssh podman "journalctl -u zachpeters -f"
