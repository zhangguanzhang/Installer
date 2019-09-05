FROM golang:1.13 as builder

WORKDIR $GOPATH/src/Installer
COPY . $GOPATH/src/Installer
ENV GO111MODULE=on 
ARG GOPROXY=https://mirrors.aliyun.com/goproxy/
#ARG GOPROXY=https://goproxy.io
RUN CGO_ENABLED=0 go build -o /root/main main.go



FROM alpine:latest
WORKDIR /root
COPY --from=builder /root/main .
EXPOSE 8080
ENTRYPOINT ["./main"]
