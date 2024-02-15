FROM docker.io/golang:1.22 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .

FROM docker.io/debian:bookworm-slim
WORKDIR /root/
COPY --from=builder /app/main .
COPY ./templates/ templates/
COPY ./assets/ assets/

EXPOSE 8080
ENV  GIN_MODE=release

CMD ["/app/main"]
