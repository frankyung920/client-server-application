package queue

import (
	"log"

	"github.com/streadway/amqp"
)

type Queue struct {
	Conn *amqp.Connection
	QueueName string
}

// Connect to RabbitMQ server
func Connect(url string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %s", err)
		return nil, err
	}
	return conn, nil
}

// PublishMessage sends a message to the specified queue
func (q *Queue) PublishMessage(body string) error {
	ch, err := q.Conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
		return err
	}
	defer ch.Close()

	
	queue, err := q.declareQueue(ch)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
		return err
	}

	err = ch.Publish(
		"",     // Default exchange
		queue.Name, // Routing key (queue name)
		false,  // RabbitMQ will silently drop the message if it cannot be routed to any queue
		false,  // Immediately deliver the message to a consumer or not, but it is deprecated
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %s", err)
		return err
	}

	return nil
}

// ConsumeMessages listens for messages on the specified queue and handles them using the provided callback function
func (q *Queue) ConsumeMessages(handleMsg func(amqp.Delivery)) {
	ch, err := q.Conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %s", err)
		return
	}
	defer ch.Close()

	queue, err := q.declareQueue(ch)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %s", err)
		return
	}

	msgs, err := ch.Consume(
		queue.Name, // Queue
		"",     // Rabbit mq will generate one for you if empty 
		true,  // true means no manual acknowledgment is required
		false,  // Allows the queue to be shared among multiple consumers 
		false,  // RabbitMQ will deliver messages to the connection that published them
		false,  // Ensure that the consumer is ready to receive messages before proceeding
		nil,    // Args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
		return
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			handleMsg(d)
		}
	}()

	log.Printf(" [*] Waiting for messages in queue %s. To exit press CTRL+C", q.QueueName)
	<-forever
}

// Make sure both consumer and publisher are using the same queue setting
func (q * Queue) declareQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		q.QueueName,
		true,  // Durable: ensures that the queue will survive a broker restart
		false, // When false, the queue will not be automatically deleted when no consumers are connected
		false, // Ensure the queue is accessible to all connections
		false, // Ensure the server will respond to the declare request
		nil,   // Other arg
	)
}