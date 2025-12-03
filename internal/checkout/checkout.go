package checkout

import(
	"fmt"
	"os"
	"log"
	"bytes"
	"io"
	"compress/zlib"
	"path/filepath"
)

func verify_hash(commit_hash string) bool {
	current_hash, err := os.ReadFile(".gogit/HEAD")
	if err != nil{
		return false
	}

	for {
		directory := current_hash[:2]
		file := current_hash[2:]
		path := filepath.Join(".gogit", "objects", string(directory), string(file))

		compressed_data, err := os.ReadFile(path)
		if err != nil{
			fmt.Fprintf(os.Stderr, "Error reading compressed data\n")
			return false
		}

		compressed_reader := bytes.NewReader(compressed_data)
		zlib_reader, err := zlib.NewReader(compressed_reader)
		if err != nil{
			fmt.Fprintf(os.Stderr, "Error creating zlib reader\n")
			return false
		}
		defer zlib_reader.Close()

		decompressed_data, err := io.ReadAll(zlib_reader)
		if err != nil{
			fmt.Fprintf(os.Stderr, "Error reading decompressed data\n")
			return false
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
		start += 6

		// Extract tree
		tree := ""
		for final_data[start] != '\n' {
			tree = tree + string(final_data[start])
			start++
		}

		fmt.Println("tree:", tree)
		fmt.Println("commit_hash:", commit_hash)
		if tree == commit_hash{
			return true
		}

		if final_data[parent_or_initial] == 'a'{
			break
		}

		parent_or_initial += 7
		parent := new(bytes.Buffer)
		for final_data[parent_or_initial] != '\n'{
			parent.WriteByte(final_data[parent_or_initial])
			parent_or_initial++
		}

		current_hash = parent.Bytes()
	}

	return false
}

func Checkout(commit_hash string){

	verified := verify_hash(commit_hash)
	fmt.Println("verified:", verified)
	return
	
	
	files, err := os.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file.Name())
		if file.Name() == ".gogit"{
			continue
		}
		os.RemoveAll(file.Name())
	}


}