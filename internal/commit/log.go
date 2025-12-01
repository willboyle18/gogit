package commit

import (
	"fmt"
	"os"
	"path/filepath"
	"compress/zlib"
	"bytes"
	"io"
)

func Latest(){
	head, err := os.ReadFile(".gogit/HEAD")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening commit HEAD\n")
		os.Exit(1)
	}

	directory := head[:2]
	file := head[2:]
	path := filepath.Join(".gogit", "objects", string(directory), string(file))

	compressed_data, err := os.ReadFile(path)
	if err != nil{
		fmt.Fprintf(os.Stderr, "Error reading compressed data\n")
		os.Exit(1)
	}

	compressed_reader := bytes.NewReader(compressed_data)
	zlib_reader, err := zlib.NewReader(compressed_reader)
	if err != nil{
		fmt.Fprintf(os.Stderr, "Error creating zlib reader\n")
		os.Exit(1)
	}
	defer zlib_reader.Close()

	decompressed_data, err := io.ReadAll(zlib_reader)
	if err != nil{
		fmt.Fprintf(os.Stderr, "Error reading decompressed data\n")
		os.Exit(1)
	}

	fmt.Println("Latest commit: ")
	fmt.Println(string(decompressed_data))
}