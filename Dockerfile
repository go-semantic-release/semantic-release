FROM alpine
ARG VERSION

ADD ./docker/entrypoint.sh /usr/local/bin/docker-entrypoint
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
COPY "./bin/semantic-release_v${VERSION}_linux_amd64" /usr/local/bin/semantic-release

ENTRYPOINT ["docker-entrypoint"]
