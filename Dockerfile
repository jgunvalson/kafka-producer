FROM golang:1.18.1 as builder
ENV GOMODULE=on \
    GOOS=linux \
    GOARCH=amd64 \
    CGO_ENABLED=0
WORKDIR /workspace
COPY ./ .
RUN go mod download
RUN go build -o producer .

FROM golang:1.18.1
WORKDIR /build
COPY --from=builder /workspace/producer /usr/bin/producer
CMD ["/usr/bin/producer"]