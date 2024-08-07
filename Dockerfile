## 1 step
#FROM golang:1.20.3-alpine AS builder
#
## . means current repo
#COPY . /github.com/AndreiMartynenko/auth/source/
#WORKDIR /github.com/AndreiMartynenko/auth/source/
#
#RUN go mod download
#RUN go build -o ./bin/crud_server grpc/cmd/grpc_server/main.go
#
## 2 step
#FROM alpine:latest
#
#WORKDIR /root/
#
## . here is copy everything in the root
#COPY --from=builder /github.com/AndreiMartynenko/auth/source/bin/crud_server .
#
## run our server
#CMD ["./crud_server"]



FROM golang:1.22.1-alpine3.19 AS builder
ARG ENV

RUN apk update && apk upgrade --available && \
    apk add make && \
    adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "10001" \
    "auth"

WORKDIR /opt/app/
COPY . .

RUN go mod download && go mod verify
RUN make build-app ENV=${ENV}
RUN mkdir logs/

FROM scratch
ARG CONFIG

WORKDIR /opt/app/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder --chown=auth:auth /opt/app/bin/main .
COPY --from=builder --chown=auth:auth /opt/app/${CONFIG} ./config
COPY --from=builder --chown=auth:auth /opt/app/tls/ ./tls/
COPY --from=builder --chown=auth:auth /opt/app/logs/ ./logs/

USER auth:auth

CMD ["./main", "-config=./config"]
