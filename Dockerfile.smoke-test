ARG GO_VERSION=1.24
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app
COPY go.mod ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /smoke-test ./cmd/smoke-test


FROM debian:bookworm

COPY --from=builder /smoke-test /usr/local/bin/
CMD ["smoke-test"]
