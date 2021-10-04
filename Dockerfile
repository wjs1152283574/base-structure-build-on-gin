FROM golang:1.15.6

ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    GIN_MODE=release \
    PORT=80
    
WORKDIR /app

COPY . .

RUN go build .

EXPOSE 80

ENTRYPOINT ["./goweb"]
