package commit

import (
	"fmt"
	"os"
	"path/filepath"
	"compress/zlib"
	"bytes"
	"io"
)

func get_head_path() string {
	head, err := os.ReadFile(".gogit/HEAD")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening commit HEAD\n")
		os.Exit(1)
	}

	directory := head[:2]
	file := head[2:]
	path := filepath.Join(".gogit", "objects", string(directory), string(file))
	return path
}

func Latest(){
	path := get_head_path()

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

	final_data := string(decompressed_data)

	start := 0
	for final_data[start] != '\x00'{
		start++
	}

	start++

	fmt.Println("Latest commit: ")
	fmt.Println(final_data[start:])
}

func Log(){
	path := get_head_path()
	for {
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

		final_data := string(decompressed_data)

		start := 0
		for final_data[start] != '\x00'{
			start++
		}

		parent_or_initial := 0
		for final_data[parent_or_initial] != '\n'{
			parent_or_initial++
		}

		parent_or_initial++
		start++
		if final_data[parent_or_initial] == 'p'{
			parent_or_initial += 7
			path = ""
			for final_data[parent_or_initial] != '\n'{
				path = path + string(final_data[parent_or_initial])
				parent_or_initial++
			}
			directory := path[:2]
			file := path[2:]
			path = filepath.Join(".gogit", "objects", directory, file)
			fmt.Println(string(final_data[start:]))
		} else if final_data[parent_or_initial] == 'a'{
			fmt.Println(final_data[start:])
			break
		}
	}
}