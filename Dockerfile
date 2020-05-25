FROM golang:1.14 as builder

WORKDIR /

ADD . .
RUN CGO_ENABLED=0 go build -o src/ccs -a -ldflags '-s' src/main.go src/orders.go src/courier.go src/shelf.go src/conf.go &&\
    chmod +x src/ccs

FROM scratch
COPY --from=builder /src/ccs /ccs
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
CMD ["/ccs"]
