package main

import (
	"bufio"
	"client-server-application/internal/queue"
	"fmt"
	"log"
	"os"
)

const EXIT_COMMAND = "exit"

func main() {
	// Assuming RabbitMQ URL and queue name are provided as environment variables
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	queueName := os.Getenv("QUEUE_NAME")

	// Connect to RabbitMQ server
	conn, err := queue.Connect(rabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
	}
	defer conn.Close()

	// Initialize a new Queue instance
	q := &queue.Queue{
		Conn:      conn,
		QueueName: queueName,
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Enter commands (type '%s' to quit): \n", EXIT_COMMAND)

	// Read commands from stdin
	for scanner.Scan() {
		text := scanner.Text()

		// Exit loop if the user types 'exit'
		if text == EXIT_COMMAND {
			break
		}

		// Publish the message to the queue
		err := q.PublishMessage(text)
		if err != nil {
			log.Fatalf("Failed to publish a message: %s", err)
		} else {
			fmt.Println("Message published:", text)
		}
	}

	if scanner.Err() != nil {
		log.Fatalf("Error reading from input: %s", scanner.Err())
	}
}
