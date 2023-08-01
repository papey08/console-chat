FROM golang:1.20

ENV GOPATH=/

WORKDIR /go/src/console-chat
COPY . .

RUN go mod download
RUN go build -o app cmd/server/main.go

CMD ["./app"]
