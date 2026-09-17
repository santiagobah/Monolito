FROM golang:1.25

WORKDIR /app

COPY go.mod ./

COPY . .

RUN go mod tidy 

RUN go build -o middleware .

EXPOSE 8100

CMD ["./middleware"]