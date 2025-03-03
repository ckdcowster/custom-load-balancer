FROM golang:1.16-alpine

WORKDIR /app

COPY load_balancer.go /app

RUN go build -o load_balancer load_balancer.go

CMD ["./load_balancer"]