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


build-image:
  # podman build --platform linux/arm64 -t registry.home.zach:5000/zachpeters.org:latest .
  podman build -t registry.home.zach:5000/zachpeters.org:latest .
  podman push registry.home.zach:5000/zachpeters.org:latest
