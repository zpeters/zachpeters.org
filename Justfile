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
  podman build -t registry.home.zach/zachpeters.org:latest .
  podman push registry.home.zach/zachpeters.org:latest
