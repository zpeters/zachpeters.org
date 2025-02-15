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
  podman build -t registry.zachpeters.org/zachpeters.org:latest .
  podman push registry.zachpeters.org/zachpeters.org:latest
