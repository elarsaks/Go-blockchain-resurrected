FROM golang:1.18.0-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/wallet_server

EXPOSE 5000

CMD ["go", "run", "-buildvcs=false", "."]
