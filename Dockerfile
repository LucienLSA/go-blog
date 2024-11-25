FROM golang:latest
# FROM scratch

ENV GOPROXY https://goproxy.cn,direct
WORKDIR $GOPATH/src/go-blog
COPY . $GOPATH/src/go-blog
RUN go build cmd/ .

EXPOSE 8000
ENTRYPOINT ["./go-blog"]