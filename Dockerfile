FROM golang:1.19.5

ENV GOPROXY=https://goproxy.cn


RUN mkdir /app

WORKDIR /app

ADD . /app

RUN go build -o main main.go

EXPOSE 8080

CMD /app/main

