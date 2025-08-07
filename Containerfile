FROM docker.io/library/golang:1.24.6-alpine as builder

WORKDIR /build

COPY ./ ./

RUN go mod download

RUN go build -o desec-dyndns-client

FROM scratch

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /build/desec-dyndns-client ./desec-dyndns-client

CMD [ "/desec-dyndns-client" ]