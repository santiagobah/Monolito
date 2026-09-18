FROM golang:1.25.0

WORKDIR /app

COPY go.mod ./
COPY . .

RUN go mod tidy
RUN go build -o loadbalancer .

EXPOSE 9100

CMD ["./loadbalancer"]
