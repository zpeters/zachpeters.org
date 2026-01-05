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
  gofmt -l .
  go vet ./...
  golangci-lint run ./...


deploy:
  podman build --isolation=chroot -t registry.home.zach/zachpeters.org:0.1.2 .
  podman push --tls-verify=false registry.home.zach/zachpeters.org:0.1.2
