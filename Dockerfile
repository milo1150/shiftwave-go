FROM golang:1.23 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o myapp .

FROM ubuntu:22.04

WORKDIR /root/

COPY --from=builder /app/myapp .

CMD ["./myapp"]



