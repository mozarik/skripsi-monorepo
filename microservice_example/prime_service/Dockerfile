FROM --platform=linux/amd64 golang:1.21-alpine AS builder
WORKDIR /app
COPY main.go .
RUN GOOS=linux GOARCH=amd64 go build -o prime main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/prime .
EXPOSE 8080
CMD ["./prime"]
