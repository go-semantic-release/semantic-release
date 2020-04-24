# Obtain latest ca-certificates
FROM alpine:latest as certs
RUN apk --update add ca-certificates

FROM scratch
ADD out/semantic-release /usr/local/bin/release
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
CMD ["release"]