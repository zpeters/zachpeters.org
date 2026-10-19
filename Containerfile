FROM docker.io/golang:1.27 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go build  .

FROM docker.io/debian:trixie-slim
WORKDIR /root/
COPY --from=builder /app/website .
# The sparkline fetches from github.com over HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY ./templates/ templates/
COPY ./assets/ assets/
COPY ./blog/ blog/
COPY ./reads/ reads/

EXPOSE 8080
ENV  GIN_MODE=release

CMD ["/root/website"]
