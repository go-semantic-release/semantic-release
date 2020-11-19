FROM alpine AS certs
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

FROM scratch
ARG VERSION

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY "./bin/semantic-release_v${VERSION}_linux_amd64" /semantic-release

ENTRYPOINT ["/semantic-release"]
