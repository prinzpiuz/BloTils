FROM golang:1.24.4-alpine3.21 AS builder

RUN apk add --no-cache build-base sqlite-dev

WORKDIR /blotils

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 go build -ldflags "-s -w" -o /blotils/BloTils .


FROM alpine:latest

WORKDIR /blotils

ARG BT_SENDGRIDAPIKEY
ARG BT_ENV

ENV BT_SENDGRIDAPIKEY=$BT_SENDGRIDAPIKEY
ENV BT_ENV=$BT_ENV

COPY --from=builder /blotils/BloTils .
COPY migrations /blotils/migrations/.
COPY static /blotils/static/.
COPY templates /blotils/templates/.

EXPOSE 8000

CMD ["/blotils/BloTils"]
