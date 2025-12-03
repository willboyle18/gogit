package checkout

import(
	"fmt"
	"os"
	"strconv"
	"bytes"
	"log"
	"io"
	"compress/zlib"
	"path/filepath"
)

func decompress_file(hash []byte) (string, bool) {
	directory := hash[:2]
	file := hash[2:]
	path := filepath.Join(".gogit", "objects", string(directory), string(file))

	compressed_data, err := os.ReadFile(path)
	if err != nil{
		fmt.Fprintf(os.Stderr, "Error reading compressed data\n")
		return "", false
	}

	compressed_reader := bytes.NewReader(compressed_data)
	zlib_reader, err := zlib.NewReader(compressed_reader)
	if err != nil{
		fmt.Fprintf(os.Stderr, "Error creating zlib reader\n")
		return "", false
	}
	defer zlib_reader.Close()

	decompressed_data, err := io.ReadAll(zlib_reader)
	if err != nil{
		fmt.Fprintf(os.Stderr, "Error reading decompressed data\n")
		return "", false
	}

	final_data := string(decompressed_data)
	return final_data, true
}

func verify_hash(commit_hash string) bool {
	current_hash, err := os.ReadFile(".gogit/HEAD")
	if err != nil{
		return false
	}

	for {
		final_data, confirmed := decompress_file(current_hash)
		if !confirmed {
			return false
		}

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

		tree := ""
		for final_data[start] != '\n' {
			tree = tree + string(final_data[start])
			start++
		}

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
	// Verify hash exists before deleting everything
	verified := verify_hash(commit_hash)
	if !verified {
		fmt.Fprintf(os.Stderr, "Hash not verified\n")
		os.Exit(1)
	}
	
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

	buffer := new(bytes.Buffer)
	buffer.WriteString(commit_hash)

	data, confirmed := decompress_file(buffer.Bytes())
	if !confirmed{
		fmt.Fprintf(os.Stderr, "Error decompressing data\n")
		os.Exit(1)
	}

	fmt.Println(data)
	index := 0
	for data[index] != '\x00'{
		index++
	}
	index++

	// Extract all files
	for {
		mode_string := ""
		for data[index] != ' '{
			mode_string = mode_string + string(data[index])
			index++
		}
		mode, err := strconv.ParseInt(mode_string, 8, 0) 
		if err != nil {
			fmt.Println("Error converting string to octal:", err)
			os.Exit(1)
		}
		fmt.Println("mode", mode)

		index++
		file := ""
		for data[index] != '\x00'{
			file = file + string(data[index])
			index++
		}
		index++
		fmt.Println("file", file)

		sha_buffer := new(bytes.Buffer)
		i := 0
		for i < 20{
			sha_buffer.WriteByte(data[index])
			index++
			i++
		}

		raw_sha := sha_buffer.Bytes()
		sha_hex := fmt.Sprintf("%x", raw_sha)
		fmt.Println("sha1", sha_hex)

		// Create the directory and file
		directory := filepath.Dir(file)
		file_name := filepath.Base(file)
		err = os.MkdirAll(directory, 0750)
		if err != nil{
			fmt.Fprintf(os.Stderr, "Error creating directory\n")
			os.Exit(1)
		}

		file_path := filepath.Join(directory, file_name)
		fd, err := os.Create(file_path)
		if err != nil{
			fmt.Fprintf(os.Stderr, "Error creating file\n")
			os.Exit(1)
		}

		// Decompress the file contents
		sha_buffer = new(bytes.Buffer)
		sha_buffer.WriteString(sha_hex)
		raw_sha = sha_buffer.Bytes()
		decompressed_blob, confirmed := decompress_file(raw_sha)
		if !confirmed {
			os.Exit(1)
		}

		// Write the decompressed contents to the file
		start := 0
		for decompressed_blob[start] != '\x00'{
			start++
		}
		start++


		_, err = fd.WriteString(decompressed_blob[start:])
		if err != nil{
			fmt.Fprintf(os.Stderr, "Failed to write to file\n")
			os.Exit(1)
		}

		if index >= len(data){
			break
		}
	}
	err = os.Remove(".gogit/index")
	if err != nil{
		fmt.Fprintf(os.Stderr, "Failed to delete index\n")
	}
}