FROM golang:1.20-alpine
WORKDIR /app

COPY go.mod ./
#COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /desec-dyndns-client

CMD [ "/desec-dyndns-client" ]