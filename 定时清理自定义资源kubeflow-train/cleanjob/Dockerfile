FROM golang:1.19.8  as builder

WORKDIR /go/src/cleanjob
COPY main.go main.go
COPY client/ client/
COPY model/ model/
COPY server server/


# Build
RUN  go env -w GOPROXY=https://goproxy.cn,direct && \
     go env -w GO111MODULE=on && \
     cd /go/src/cleanjob/ && \
     go mod init cleanjob && \
     go mod tidy && \
     CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=mod -o builder  /go/src/cleanjob/main.go

FROM 192.168.2.106:1180/fangzhou/alpine:v1.0.0

WORKDIR /root/

#RUN apk update && apk upgrade && \
#    apk add --no-cache bash git openssh

COPY --from=builder /go/src/cleanjob/builder .
CMD ["./builder"]
