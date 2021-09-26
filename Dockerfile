FROM golang:1.17

COPY . /go/src/app

WORKDIR /go/src/app/cmd/broker

RUN go build -o broker main.go

EXPOSE 5300

CMD ["./broker"]