FROM golang:1.25

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY ./src/go.mod ./src/go.sum .
RUN go mod download

COPY ./src .
RUN go build -v -o /usr/local/bin/webhook-server main.go

# Webhook listens on 8443
EXPOSE 8443

CMD ["webhook-server"]