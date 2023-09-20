FROM golang:1.21.0

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /go-docker-bin

EXPOSE 9099 8081 8082

CMD ["/go-docker-bin"]

