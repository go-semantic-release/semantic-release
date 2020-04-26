FROM golang:alpine AS builder
ARG VERSION

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR $GOPATH/src/semantic-release
COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a --installsuffix cgo -ldflags="-extldflags \'-static\' -s -w -X main.SRVERSION=$VERSION" -o /go/bin/semantic-release ./cmd/semantic-release/


FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/semantic-release /

ENTRYPOINT ["/semantic-release"]
