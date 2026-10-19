image := "registry.home.zach/zachpeters.org"
version := "0.2.0"
# This machine is arm64; the cluster is x86_64, so build on a remote x86
# podman. Create the connection once:
#   podman system connection add cluster ssh://user@host/run/user/1000/podman/podman.sock
builder := "cluster"

build:
  # Build your Go application
  go build -o bin/website .

run: build
  # Run your application locally
  ./bin/website

proxy: build
  tmux kill-session -t dev 2>/dev/null || true
  tmux new-session -d -s dev -n main 'mitmproxy --mode reverse:http://localhost:8080 --listen-port 8443'
  tmux split-window -t dev -v 'fd . --exclude bin/ --exclude tmp/ | entr -r go run .'
  tmux attach -t dev

stop:
  tmux kill-session -t dev 

dev:
  fd . --exclude bin/ --exclude tmp/ | entr -r go run .

lint:
  # Lint your Go codebase
  fd -e go -E vendor -X gofmt -l
  go vet ./...
  golangci-lint run ./...

vuln:
  # Check for known vulnerabilities
  govulncheck ./...

# build the image on the remote x86 podman
image-build:
  podman --connection {{builder}} build --platform linux/amd64 -t {{image}}:{{version}} .

# fallback: build here via qemu emulation
image-build-local:
  podman build --platform linux/amd64 -t {{image}}:{{version}} .

# push the image to the home registry (from the remote builder)
image-push:
  podman --connection {{builder}} push --tls-verify=false {{image}}:{{version}}

deploy: image-build image-push
