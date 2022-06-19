FROM golang:latest

WORKDIR go/src/app

COPY . .

RUN go build -o ./build/server main.go

CMD ["./build/server"]