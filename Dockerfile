FROM golang:1.16 AS builder
WORKDIR /build
COPY . .
ENV CGO_ENABLED=0
RUN go build -ldflags '-w' -o server ./server.go

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/server /bin/cmd
USER 1000
ENTRYPOINT ["/bin/cmd"]