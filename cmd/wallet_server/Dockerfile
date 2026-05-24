FROM golang:1.18.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -buildvcs=false -o /bin/wallet-server ./cmd/wallet_server

FROM alpine:3.19

COPY --from=builder /bin/wallet-server /bin/wallet-server

EXPOSE 5000

CMD ["/bin/wallet-server"]
