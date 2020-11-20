FROM alpine
ARG VERSION

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
COPY "./bin/semantic-release_v${VERSION}_linux_amd64" /usr/local/bin/semantic-release

ENTRYPOINT ["semantic-release"]
