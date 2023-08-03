FROM golang:1.19.8  as builder

WORKDIR /go/src/webhook
COPY main.go main.go
COPY routers/ routers/
COPY model/ model/
COPY server server/


# Build
RUN  go env -w GOPROXY=https://goproxy.cn,direct && \
     go env -w GO111MODULE=on && \
     cd /go/src/webhook/ && \
     go mod init webhook && \
     go mod tidy && \
     CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=mod -o webhook  /go/src/webhook/main.go

FROM 192.168.2.106:1180/fangzhou/alpine:v1.0.0

WORKDIR /root/

#RUN apk update && apk upgrade && \
#    apk add --no-cache bash git openssh

COPY --from=builder /go/src/webhook/webhook .
CMD ["./webhook"]
