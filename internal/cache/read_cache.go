package cache

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

var sha1_file_directory string

func get_name(buffer *bytes.Reader, name_length uint16) []byte {
	i := 0
	var name []byte
	for i < int(name_length) {
		next_byte, err := buffer.ReadByte()
		if err != nil {
			log.Fatal(err)
		}
		name = append(name, next_byte)
		i++
	}

	// Get through padding
	for {
		next_byte, err := buffer.ReadByte()
		if err == io.EOF {
			break
		} else if err != nil{
			log.Fatal(err)
		}
		if next_byte != 0x00 {
			break
		}
	}
	err := buffer.UnreadByte()
	if err != nil {
		log.Fatal(err)
	}
	return name
}

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
	_, err = os.Open(".gogit/index")
	if err != nil {
		if errors.Is(err, os.ErrNotExist){
			return 0;
		} else {
			log.Fatal(err)
		}
	}

	data, err := os.ReadFile(".gogit/index")
	if err != nil {
		log.Fatal(err)
	}

	var hdr CacheHeader
	buffer := bytes.NewReader(data)
	err = binary.Read(buffer, binary.BigEndian, &hdr)
	if err != nil{
		log.Fatal(err)
	}

	entries := int(hdr.Entries)

	for i := 0; i < entries; i++ {
		var entry CacheEntry

		binary.Read(buffer, binary.BigEndian, &entry.Ctime.Sec)
		binary.Read(buffer, binary.BigEndian, &entry.Ctime.Nsec)
		binary.Read(buffer, binary.BigEndian, &entry.Mtime.Sec)
		binary.Read(buffer, binary.BigEndian, &entry.Mtime.Nsec)
		binary.Read(buffer, binary.BigEndian, &entry.Dev)
		binary.Read(buffer, binary.BigEndian, &entry.Ino)
		binary.Read(buffer, binary.BigEndian, &entry.Mode)
		binary.Read(buffer, binary.BigEndian, &entry.Uid)
		binary.Read(buffer, binary.BigEndian, &entry.Gid)
		binary.Read(buffer, binary.BigEndian, &entry.Size)
		binary.Read(buffer, binary.BigEndian, &entry.Sha1)
		binary.Read(buffer, binary.BigEndian, &entry.Name_Length)

		name := get_name(buffer, entry.Name_Length)

		entry.Name = string(name)

		ActiveCache = append(ActiveCache, &entry)
		ActiveNR++
	}

	return 0
}