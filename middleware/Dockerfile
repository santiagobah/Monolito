FROM golang:1.23.6

WORKDIR /app

COPY go.mod ./

COPY . .

RUN go mod tidy 

RUN go build -o middleware .

EXPOSE 8083

CMD ["./middleware"]