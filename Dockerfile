FROM golang:1.8.3-alpine3.6

RUN apk --no-cache add curl \ 
    && echo "Pulling watchdog binary from Github." \
    && curl -sSL https://github.com/openfaas/faas/releases/download/0.6.5/fwatchdog > /usr/bin/fwatchdog \
    && chmod +x /usr/bin/fwatchdog \
    && apk del curl --no-cache

WORKDIR /go/src/handler
COPY . .

# Run a gofmt and exclude all vendored code.
RUN test -z "$(gofmt -l $(find . -type f -name '*.go' -not -path "./function/vendor/*"))" || { echo "Run \"gofmt -s -w\" on your Golang code"; exit 1; }

RUN CGO_ENABLED=0 GOOS=linux \
    go build --ldflags "-s -w" -a -installsuffix cgo -o handler . && \
    go test $(go list ./... | grep -v /vendor/) -cover

FROM alpine:3.6
RUN apk --no-cache add ca-certificates

# Add non root user
RUN addgroup -S app && adduser -S -g app app
RUN mkdir -p /home/app
RUN chown app /home/app

WORKDIR /home/app

COPY --from=0 /go/src/handler/handler    .
COPY --from=0 /usr/bin/fwatchdog         .
COPY --from=0 /go/src/handler/function/monorest.yml          .
COPY --from=0 /go/src/handler/function/english.tet          .
COPY --from=0 /go/src/handler/function/english.tri          .
COPY --from=0 /go/src/handler/function/wp.3          .
COPY --from=0 /go/src/handler/function/wp.3.spaces          .
COPY --from=0 /go/src/handler/function/wp.4          .
COPY --from=0 /go/src/handler/function/wp.4.spaces          .
COPY --from=0 /go/src/handler/function/wp.5.spaces          .
COPY --from=0 /go/src/handler/function/wp.6.spaces          .
COPY --from=0 /go/src/handler/function/wp7          .
RUN chown app /home/app/monorest.yml
RUN chown app /home/app/english.tet
RUN chown app /home/app/english.tri
RUN chown app /home/app/wp.3
RUN chown app /home/app/wp.3.spaces
RUN chown app /home/app/wp.4
RUN chown app /home/app/wp.4.spaces
RUN chown app /home/app/wp.5.spaces
RUN chown app /home/app/wp.6.spaces
RUN chown app /home/app/wp7

USER app

ENV fprocess="./handler"

CMD ["./fwatchdog"]
