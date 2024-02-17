#!/bin/bash
set -eo pipefail

podman build --platform linux/amd64 -t podman:5000/zachpeters.org-dev:latest .
podman push podman:5000/zachpeters.org-dev:latest
ssh podman "podman auto-update"
ssh podman "systemctl stop zachpeters-dev"
ssh podman "systemctl start zachpeters-dev"
ssh podman "journalctl -u zachpeters-dev -f"
