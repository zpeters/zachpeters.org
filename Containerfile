FROM docker.io/golang:1.22 as builder
WORKDIR /app
COPY go.mod go.sum main.go ./
RUN go mod download
RUN go build  .

FROM docker.io/debian:bookworm-slim
WORKDIR /root/
COPY --from=builder /app/website .
COPY ./templates/ templates/
COPY ./assets/ assets/
COPY ./CHANGELOG.md CHANGELOG.md

EXPOSE 8080
ENV  GIN_MODE=release

CMD ["/root/website"]
