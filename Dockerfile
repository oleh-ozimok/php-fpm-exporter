FROM golang:1.11.6-alpine3.9 as builder

LABEL maintainer="Oleg Ozimok oleg.ozimok@corp.kismia.com"

COPY . /go/src/github.com/oleh-ozimok/php-fpm-exporter

WORKDIR /go/src/github.com/oleh-ozimok/php-fpm-exporter

RUN go build -tags=jsoniter -o /php-fpm-exporter ./cmd/php-fpm-exporter

FROM alpine:3.9

COPY --from=builder /php-fpm-exporter /usr/bin/php-fpm-exporter

EXPOSE 8080

STOPSIGNAL SIGTERM

ENTRYPOINT ["/usr/bin/php-fpm-exporter"]