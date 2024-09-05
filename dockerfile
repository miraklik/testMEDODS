FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod dowloand

COPY . .

RUN go build -o auth-service ./cmd/main

EXPOSE 8080

CMD ["./golang"]