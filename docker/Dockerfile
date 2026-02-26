
# Build stage
FROM golang:1.24.4-alpine3.22 AS builder

RUN apk add --no-cache build-base sqlite-dev

WORKDIR /build

# Cache dependencies
COPY ../go.mod ../go.sum ./
RUN go mod download


COPY ../ .

# Build arguments
ARG VERSION=unknown
ARG COMMIT_SHA=unknown
ARG BUILD_TIME=unknown

# Build with embedded version info
RUN CGO_ENABLED=1 go build \
    -ldflags "-s -w \
    -X BloTils/src/models.Version=${VERSION} \
    -X BloTils/src/models.CommitSHA=${COMMIT_SHA} \
    -X BloTils/src/models.BuildTime=${BUILD_TIME}" \
    -o /build/BloTils .

# Runtime stage
FROM alpine:3.22

RUN apk add --no-cache sqlite su-exec && \
    addgroup -g 1001 -S blotils && \
    adduser -u 1001 -S blotils -G blotils && \
    mkdir -p /blotils && \
    chown blotils:blotils /blotils

WORKDIR /blotils

# Copy application files
COPY --from=builder /build/BloTils .
COPY --chown=blotils:blotils ../migrations/ ./migrations/
COPY --chown=blotils:blotils ../static/ ./static/
COPY --chown=blotils:blotils ../templates/ ./templates/
COPY --chown=blotils:blotils ../docker/entrypoint.sh .
RUN chmod +x entrypoint.sh

EXPOSE 8000

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8000/api/v1/ping || exit 1

USER root
ENTRYPOINT ["/blotils/entrypoint.sh"]
CMD ["/blotils/BloTils"]
