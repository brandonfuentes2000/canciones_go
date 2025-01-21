FROM golang:1.23-alpine

WORKDIR /canciones

COPY . .

RUN go mod download

RUN go build -o main ./cmd

EXPOSE 3000

CMD ["./main"]
