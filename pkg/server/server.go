package server

import (
	"sync"

	"github.com/streadway/amqp"
)

type Server struct {
    queueClient  QueueClient
    commandHandler CommandHandler
}

func NewServer(qc QueueClient, ch CommandHandler) *Server {
    return &Server{
        queueClient:  qc,
        commandHandler: ch,
    }
}

func (s *Server) Run() {
    var wg sync.WaitGroup

    s.queueClient.ConsumeMessages(func(d amqp.Delivery) {
        wg.Add(1)
        go func() {
            defer wg.Done()
            s.commandHandler.HandleCommand(string(d.Body))
        }()
    })

    wg.Wait() // Wait for all goroutines to finish
}
