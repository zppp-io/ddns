FROM golang AS builder

ENV GOPROXY https://goproxy.io,direct
ENV GO111MODULE on
WORKDIR $GOPATH/src/app

COPY .  .

RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -tags=jsoniter   -o /ddns

#RUNTIME
FROM alpine
MAINTAINER Echo echo@zppp.io
WORKDIR /home

COPY --from=builder /ddns /bin/ddns

CMD ["ddns"]