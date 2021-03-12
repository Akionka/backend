FROM golang:1.16.0-alpine3.13

WORKDIR /build
COPY . .
RUN go build -o /app/app cmd/main.go

WORKDIR /app

ENTRYPOINT ["./app"]
EXPOSE 8080