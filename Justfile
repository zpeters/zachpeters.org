build:
  # Build your Go application
  go build -o bin/website .

run: build
  # Run your application locally
  ./bin/website

dev:
  fd . --exclude bin/ --exclude tmp/ | entr -r go run .

lint:
  # Lint your Go codebase
  gofmt -l .
  go vet ./...
  golangci-lint run ./...


deploy:
  # Build the container with podman
  podman build --platform linux/amd64 -t podman:5000/zachpeters.org:latest .
  # Push the container image
  podman push podman:5000/zachpeters.org:latest
  # Remote commands to update the service
  ssh root@podman "podman auto-update"
  ssh root@podman "systemctl stop zachpeters"
  ssh root@podman "systemctl start zachpeters"
  ssh root@podman "if systemctl is-active --quiet zachpeters; then echo -e '\033[32mOK: Service started successfully\033[0m'; else echo -e '\033[31mFAIL: Service failed to start\033[0m'; fi"
