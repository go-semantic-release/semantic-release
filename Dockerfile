FROM golang:alpine AS builder
ARG VERSION

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR $GOPATH/src/semantic-release
COPY . .
COPY docker-entrypoint.sh /semantic-release/

RUN GOOS=linux GOARCH=amd64 go build -ldflags="-extldflags '-static' -s -w -X main.SRVERSION=$VERSION" -o /go/bin/semantic-release ./cmd/semantic-release/

## Build clean image
FROM alpine

COPY --from=builder /etc/profile /etc/profile
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/semantic-release /usr/local/bin
COPY --from=builder /semantic-release /semantic-release

RUN chmod +x /semantic-release/docker-entrypoint.sh

ENTRYPOINT ["/semantic-release/docker-entrypoint.sh"]
CMD ["semantic-release"]
