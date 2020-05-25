FROM golang:1.14 as builder

WORKDIR /app
ENV GOPROXY=http://goproxy.ymt360.com
ENV GOSUMDB=off

ADD . .
RUN CGO_ENABLED=0 go build -o src/ccs -a -ldflags '-s' src/main.go &&\
    chmod +x src/ccs

FROM scratch
COPY --from=builder /usr/local/go/lib/time/zoneinfo.zip /usr/local/go/lib/time/zoneinfo.zip
CMD ["/app/src/ccs"]
