FROM golang:1.20.0-alpine

RUN mkdir /app
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go get -u github.com/cosmtrek/air && go build -o /go/bin/air github.com/cosmtrek/air
CMD ["air", "-c", ".air.toml"]
