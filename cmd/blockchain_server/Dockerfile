FROM golang:1.18.0-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/blockchain_server

EXPOSE 5001 5002 5003

CMD ["go", "run", "-buildvcs=false", "."]
