package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	flags "github.com/jessevdk/go-flags"
	"github.com/vice-registry/vice-store/actions"
	"github.com/vice-registry/vice-store/storage"
	"github.com/vice-registry/vice-util/communication"
	"github.com/vice-registry/vice-util/persistence"
)

func main() {

	// default flags
	var opts struct {
	}

	// CouchbaseFlags cli Configuration options for couchbase connection
	var couchbaseFlags = struct {
		Location string `short:"" long:"couchbase-location" description:"Location of the Couchbase cluster to connect to (e.g. localhost)"`
		Username string `short:"" long:"couchbase-user" description:"Username to log in to Couchbase cluster"`
		Password string `short:"" long:"couchbase-pass" description:"Password to log in to Couchbase cluster"`
	}{}

	// RabbitmqFlags cli Configuration options for rabbitmq connection
	var rabbitmqFlags = struct {
		Location string `short:"" long:"rabbitmq-location" description:"Location of the RabbitMQ to connect to (e.g. localhost)"`
		Username string `short:"" long:"rabbitmq-user" description:"Username to log in to RabbitMQ"`
		Password string `short:"" long:"rabbitmq-pass" description:"Password to log in to RabbitMQ"`
	}{}

	// StorageFlags cli Configuration options for rabbitmq connection
	var storageFlags = struct {
		Basepath string `short:"" long:"storage-basepath" description:"Basepath to store the imported images"`
	}{}

	// initialize parser for flags
	parser := flags.NewParser(&opts, flags.Default)
	parser.ShortDescription = "ViCE Image Registry Store"
	parser.LongDescription = "Store component of the ViCE Image Registry"
	parser.AddGroup("Couchbase Connection", "Configuration options for couchbase connection", &couchbaseFlags)
	parser.AddGroup("RabbitMQ Connection", "Configuration options for RabbitMQ connection", &rabbitmqFlags)
	parser.AddGroup("Storage Connection", "Configuration options for Image Storage", &storageFlags)
	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	log.Printf("Starting vice-store service ...")

	// catch SIGINT signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		shutdown()
	}()

	// initialize couchbase
	persistence.SetCouchbaseCredentials(couchbaseFlags.Location, couchbaseFlags.Username, couchbaseFlags.Password)
	persistence.InitViceCouchbase()

	// initialize rabbitmq
	err := communication.SetRabbitmqCredentials(rabbitmqFlags.Location, rabbitmqFlags.Username, rabbitmqFlags.Password)
	if err != nil {
		log.Printf("Unable to connect to RabbitMQ: %s", err)
		shutdown()
	}

	// initialize storage
	storage.SetStorageConfig(storageFlags.Basepath)

	log.Print("Wait for incoming actions...")
	err = actions.WaitForActions()
	if err != nil {
		log.Printf("Cannot WaitForActions: %s", err)
		shutdown()
	}

}

func shutdown() {
	// clean up before termination
	persistence.CloseConnection()
	communication.CloseConnection()
	os.Exit(1)
}
