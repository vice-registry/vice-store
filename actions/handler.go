package actions

import (
	"log"
	"net"
	"time"

	"github.com/vice-registry/vice-store/storage"
	"github.com/vice-registry/vice-util/persistence"
	"github.com/vice-registry/vice-util/storeclient"
)

func handleAction(request storeclient.StoreRequest) error {
	// get image from persistence
	image, err := persistence.GetImage(request.ImageID)
	if err != nil {
		return err
	}

	// connect to remote location
	connection, err := net.DialTimeout("tcp", request.Connection, 10*time.Second)
	if err != nil {
		log.Printf("Failed to connect to %s: %s", request.Connection, err)
		return err
	}
	log.Printf("Connect to %s, start transfer...", request.Connection)

	log.Printf("action: %s", request.Action)

	if request.Action == "retrieve" {
		// read file locally
		log.Printf("Going to retrieve file from store")		
		storage.RetrieveImage(image, connection)
	} else if request.Action == "store" {
		// store file locally
		log.Printf("Going to save file in store")		
		storage.StoreImage(image, connection)
	}

	connection.Close()
	log.Printf("Finished transfer to %s.", request.Connection)

	return nil
}
