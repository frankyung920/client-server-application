package server

import "github.com/streadway/amqp"

// CommandHandler defines an interface for handling commands.
type CommandHandler interface {
	HandleCommand(command string)
}

// QueueClient defines an interface for queue operations.
type QueueClient interface {
	ConsumeMessages(handleMsg func(amqp.Delivery))
}
