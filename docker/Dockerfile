FROM golang:1.24.4-alpine3.21 AS builder

RUN apk add --no-cache build-base sqlite-dev

WORKDIR /blotils

COPY ../go.mod .
COPY ../go.sum .

RUN go mod download

COPY ../ .

RUN CGO_ENABLED=1 go build -ldflags "-s -w" -o /blotils/BloTils .


FROM alpine:latest

RUN apk add --no-cache sqlite su-exec && \
    addgroup -g 1001 -S blotils && \
    adduser -u 1001 -S blotils -G blotils

WORKDIR /blotils

ARG BT_SENDGRIDAPIKEY
ARG BT_ENV

ENV BT_SENDGRIDAPIKEY=$BT_SENDGRIDAPIKEY
ENV BT_ENV=$BT_ENV

COPY --from=builder /blotils/BloTils .

COPY ../migrations /blotils/migrations/.
COPY ../static /blotils/static/.
COPY ../templates /blotils/templates/.

COPY ../docker/entrypoint.sh /blotils/.
RUN chmod +x /blotils/entrypoint.sh

EXPOSE 8000

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8000/api/v1/ping || exit 1

USER root
ENTRYPOINT ["/blotils/entrypoint.sh"]
CMD ["/blotils/BloTils"]
