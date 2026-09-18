FROM golang:1.25.0

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /enlaces

EXPOSE 8080

CMD ["/enlaces"]