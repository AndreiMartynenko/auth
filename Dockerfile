FROM golang:1.20.3-alpine AS builder

COPY . /github.com/AndreiMartynenko/auth/
WORKDIR /github.com/AndreiMartynenko/auth/

RUN go mod download
RUN go build -o ./bin/crud_server cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/AndreiMartynenko/auth/source/bin/crud_server .

CMD ["./crud_server"]