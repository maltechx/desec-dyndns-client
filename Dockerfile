FROM library/golang:1.24.0-alpine
WORKDIR /app

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /desec-dyndns-client

CMD [ "/desec-dyndns-client" ]