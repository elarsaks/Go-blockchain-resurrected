FROM golang:1.18.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -buildvcs=false -o /bin/blockchain-server ./cmd/blockchain_server

FROM alpine:3.19

COPY --from=builder /bin/blockchain-server /bin/blockchain-server

EXPOSE 5001 5002 5003

CMD ["/bin/blockchain-server"]
