package actions

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/vice-registry/vice-util/communication"
	"github.com/vice-registry/vice-util/storeclient"
)

// WaitForActions listens on RabbitMQ channel and accepts one message at a time
func WaitForActions() error {

	msgs, err := communication.NewConsumer("store")
	if err != nil {
		log.Printf("Error while registering new consumer: %s", err)
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			request := storeclient.StoreRequest{}
			err := json.Unmarshal(d.Body, &request)
			if err != nil {
				fmt.Println("error:", err)
			}
			err = handleAction(request)
			if err != nil {
				log.Printf("Failed to handle action for request %v: %s", request, err)
			}
			d.Ack(false)
		}
	}()

	// wait forever until interrupted
	<-forever

	return nil
}
