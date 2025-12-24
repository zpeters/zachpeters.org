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
  podman build --isolation=chroot -t registry.home.zach/zachpeters.org:0.1.1 .
  podman push --tls-verify=false registry.home.zach/zachpeters.org:0.1.1
