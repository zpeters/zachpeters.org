FROM docker.io/golang:1.27 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
COPY internal/ internal/
ARG VERSION=dev
RUN go build -ldflags "-X main.version=${VERSION}" .

FROM docker.io/debian:trixie-slim
WORKDIR /srv/site
COPY --from=builder /app/website .
# The sparkline fetches from github.com over HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY ./templates/ templates/
COPY ./pages/ pages/
COPY ./projects/ projects/
COPY ./assets/ assets/
COPY ./blog/ blog/
COPY ./reads/ reads/

EXPOSE 8080
ENV  GIN_MODE=release
# Run as a non-root user; the site only reads its own files
USER 10001:10001

CMD ["/srv/site/website"]
