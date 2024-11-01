FROM golang:1.23.2

WORKDIR /app

COPY . .

RUN apt-get update

RUN go mod download
RUN go build -o link-base ./cmd/main.go

EXPOSE 8080

CMD ["./link-base"]