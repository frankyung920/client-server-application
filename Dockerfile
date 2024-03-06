FROM arm64v8/golang:1.21

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o main ./cmd/server

CMD ["./main"]