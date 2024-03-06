# Client-Server Application with RabbitMQ

This repository contains the implementation of a Client-Server application using RabbitMQ as the message broker. The server reads commands from an external RabbitMQ queue and executes them, while clients can send commands to the server via the same RabbitMQ queue.

## Features

- **Server**: Processes commands received from the RabbitMQ queue and executes them.
- **Client**: Sends commands to the RabbitMQ queue.
- **Ordered Map**: The server uses an ordered map data structure for storing key-value pairs.
- **Scalability**: The server and clients are designed to scale quickly and serve multiple clients simultaneously.

## Folder Structure

The project is structured as follows:

client-server-application/

```
│
├── cmd/
│ ├── client
│ │ └── main.go # Entry point for client
│ │
│ └── server
│   └── main.go # Entry point for server
│
├── internal/
│ ├── orderedmap/ # Implementation of the ordered map
│ └── queue/ # Implementation of RabbitMQ
│
├── pkg/
│ └── server/ # Implementation of the server and related components
│
├── README.md
├── docker-compose.yml # For spinning up rabbit mq + server
├── Dockerfile # Simple Docker file for dockerizing the server
├── Makefile # Some useful commands
├── go.mod
└── go.sum
```

## Assumptions

- **RabbitMQ Setup**: Assumes that RabbitMQ server is already installed and running on the local machine or at a specified URL.
- **Data Format**: Below are some date format assumptions
  1. Data exchanged between clients and the server are in the form of strings
  2. No character ’ in between key or value as they will be removed when the application extracts the command to be arguments
- **Docker**: Assumes Docker is installed in the PC
- **Environment Variables**
  1. Environment Variables RABBITMQ_URL is set for RabbitMQ connection info
  2. Environment Variables QUEUE_NAME is set for RabbitMQ queue name

## Dependencies

The project uses the following dependencies:

github.com/streadway/amqp: Go client for RabbitMQ.

## To run the app

1. By using docker-compose
   For server side, it support using docker-compose to spin up the application with rabbit-mq (see docker-compose.yml), we can still dockerizing the client but it may not be helpful as the client is for us to input command.
   We can run below command to start the server app and rabbit mq

```
docker-compose up -d
```

2. For client or if you want to run server locally
   Please run below command, where `cmd/client/main.go` is for client and `cmd/server/main.go` is for server

```
# For client
QUEUE_NAME=task_queue RABBITMQ_URL=amqp://guest:guest@localhost:5672 go run cmd/client/main.go

# For server
QUEUE_NAME=task_queue RABBITMQ_URL=amqp://guest:guest@localhost:5672 go run cmd/client/server.go
```

We can make sure of the Makefile to start a rabbit mq instance

```
make run-mq
```

3. To ensure the app can be compiled in ARM64, we can do

```
GOARCH=arm64 go build ./...
```

## Testing

We can make use of the Makefile to run test or get test coverage

```
# Run test
make test

# Run test coverage
make coverage
```
