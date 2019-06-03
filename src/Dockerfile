FROM golang:1.12.1

WORKDIR /go/src/github.com/jaquiel/golang-service/src
ADD . .

RUN go get
RUN go build -o main .
CMD ["go", "run", "main.go"]