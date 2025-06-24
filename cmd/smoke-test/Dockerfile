ARG GO_VERSION=1.24
FROM golang:${GO_VERSION}-bookworm AS builder

WORKDIR /usr/src/app
COPY go.mod ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /app .


FROM debian:bookworm

COPY --from=builder /app /usr/local/bin/
CMD ["app"]
