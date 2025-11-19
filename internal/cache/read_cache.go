package cache

import (
	"errors"
	"fmt"
	"log"
	"os"
)

var sha1_file_directory string

func Read_Cache() int {

	// Block 1: Prevent loading multiple index files
	if ActiveCache != nil {
		fmt.Fprintf(os.Stderr, "more then one cachefile")
		os.Exit(1)
	}

	// Block 2: Determine the object directory
	sha1_file_directory = ".gogit/objects"
	_, err := os.Stat(sha1_file_directory)
	if err != nil {
		log.Fatal(err)
	}

	// Block 3: Check that object directory is accessible
	f, err := os.CreateTemp(sha1_file_directory, "example")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(f.Name()) // clean up

	if _, err := f.Write([]byte("test")); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

	// Block 4: Open the index file
	fd, err := os.Open(".gogit/index")
	if err != nil {
		if errors.Is(err, os.ErrNotExist){
			return 0;
		} else {
			log.Fatal(err)
		}
	}
}