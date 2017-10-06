package storage

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/vice-registry/vice-util/models"
)

var storageConfig = struct {
	Basepath string
}{}

// SetStorageConfig provide storage configuration
func SetStorageConfig(basepath string) {
	storageConfig.Basepath = basepath
}

// TODO extract this storage logic to separat layer, outside of vice-import!

// StoreImage stores an image in the specified location on the file system
func StoreImage(image *models.Image, reader io.Reader) error {
	filepath := storageConfig.Basepath + "/" + image.ID + ""
	// open output file
	file, err := os.Create(filepath)
	if err != nil {
		log.Printf("Error in storage: failed to create file for imageID %s: %s", image.ID, err)
		return err
	}
	// close fo on exit and check for its returned error
	defer func() {
		err := file.Close()
		if err != nil {
			log.Printf("Error in storage: failed to close file for imageID %s: %s", image.ID, err)
		}
	}()

	// make a write buffer
	writer := bufio.NewWriter(file)

	// make a buffer to keep chunks that are read
	buffer := make([]byte, 1024)
	for {
		// read a chunk
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			log.Printf("Error in storage: failed to read from reader for imageID %s: %s", image.ID, err)
			return err
		}
		if n == 0 {
			break
		}

		// write a chunk
		if _, err := writer.Write(buffer[:n]); err != nil {
			if err != nil {
				log.Printf("Error in storage: failed to write to file for imageID %s: %s", image.ID, err)
				return err
			}
		}
	}

	if err = writer.Flush(); err != nil {
		log.Printf("Error in storage: failed to flush written file for imageID %s: %s", image.ID, err)
		return err
	}

	return nil
}

// RetrieveImage returns File pointer to an image in the specified location on the file system
func RetrieveImage(image *models.Image) (*os.File, error) {
	filepath := storageConfig.Basepath + "/" + image.ID + ""
	file, err := os.Open(filepath)
	return file, err
}
