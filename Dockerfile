FROM golang:1.22.1-alpine AS builder

WORKDIR /workspace

RUN apk add --update --no-cache git && rm -rf /var/cache/apk/*
COPY go.mod go.sum /workspace/
COPY go /workspace/go
RUN go build ./go/cmd/swapi-server

FROM alpine
RUN apk add --update --no-cache ca-certificates tzdata && rm -rf /var/cache/apk/*
COPY --from=builder /workspace/swapi-server /usr/local/bin/swapi-server
CMD [ "/usr/local/bin/swapi-server", "--port=8080", "--bind=0.0.0.0", "--embed-gateway" ]
