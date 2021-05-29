FROM golang:1.15

WORKDIR $GOPATH/src/github.com/incafox/golang-api

COPY go.mod .

COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o server main.go

EXPOSE 4000
EXPOSE 5432

CMD ["./server"]
