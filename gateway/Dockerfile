FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod routes.json main.go ./
RUN go build -o gateway .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/gateway .
COPY --from=builder /app/routes.json .
EXPOSE 8080
CMD ["./gateway"]
