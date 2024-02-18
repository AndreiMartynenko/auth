# 1st step
FROM golang:1.20.3-alpine AS builder

# . means current repo
COPY . /github.com/AndreiMartynenko/auth/grpc/source/
WORKDIR /github.com/AndreiMartynenko/auth/grpc/source/

RUN go mod download
RUN go build -o ./bin/crud_server cmd/grpc_server/main.go

# 2 step
FROM alpine:latest

WORKDIR /root/

# . here is copy everything in the root
COPY --from=builder /github.com/AndreiMartynenko/auth/source/bin/crud_server .

# run our server
CMD ["./crud_server"]