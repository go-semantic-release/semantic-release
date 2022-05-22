FROM alpine

ADD ./docker/entrypoint.sh /usr/local/bin/docker-entrypoint
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
COPY "./dist/semantic-release_linux_amd64/semantic-release" /usr/local/bin/semantic-release

ENTRYPOINT ["docker-entrypoint"]
