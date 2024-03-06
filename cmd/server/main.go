package main

import (
	"client-server-application/internal/orderedmap"
	"client-server-application/internal/queue"
	"client-server-application/pkg/server"
	"log"
	"os"
	"sync"
)

func main() {
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	queueName := os.Getenv("QUEUE_NAME")

	conn, err := queue.Connect(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	q := &queue.Queue{
		Conn:      conn,
		QueueName: queueName,
	}

	om := orderedmap.NewOrderedMap()
	mu := &sync.Mutex{}
	ch := server.NewOrderedMapCommandHandler(om, mu)

	srv := server.NewServer(q, ch)
	srv.Run()
}
