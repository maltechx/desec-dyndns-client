FROM docker.io/library/golang:1.25.5-alpine as builder

WORKDIR /build

COPY ./ ./

RUN go mod download

RUN go build ./cmd/desecdyndns

FROM scratch

WORKDIR /app

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder /build/desecdyndns /desecdyndns

CMD [ "/desecdyndns" ]