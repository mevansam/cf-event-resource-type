FROM alpine:latest

RUN \
    apk --no-cache add ca-certificates && update-ca-certificates

COPY cmd/check/check /opt/resource/
COPY cmd/in/in /opt/resource/
COPY cmd/out/out /opt/resource/
