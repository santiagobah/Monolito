FROM golang:1.21

WORKDIR /app

COPY go.mod ./
COPY . .

# esto descarga las dependencias (nomas el driver de postgres) y compila
RUN go mod tidy
RUN go build -o server .

EXPOSE 8080

CMD ["./server"]
