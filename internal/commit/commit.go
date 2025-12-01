package commit

import (
	"fmt"
	"os"
	"os/user"
	"bytes"
	"time"
	"crypto/sha1"
	"compress/zlib"
	"path/filepath"
)

func get_timestamp() string {
	now := time.Now()

	seconds := now.Unix()

	_, offset_seconds := now.Zone()

	sign := "+"
	if offset_seconds < 0{
		sign = "-"
		offset_seconds = -offset_seconds
	}

	offset_hours := offset_seconds / 3600
	offset_minutes := (offset_seconds % 3600) / 60

	return fmt.Sprintf("%d %s%02d%02d", seconds, sign, offset_hours, offset_minutes)
}

func Check_Message_Format(message string) string {
	if len(message) == 0{
		return ""
	}

	fmt.Println(message)

	return message
}

func create_commit_object(message string, sha1_hex string) []byte {
	buffer := new(bytes.Buffer)
	buffer.WriteString(fmt.Sprintf("tree %s\n", sha1_hex))

	fd, err := os.OpenFile(".gogit/HEAD", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error accessing commit HEAD\n")
		os.Exit(1)
	}
	defer fd.Close()

	stats, err := os.Stat(".gogit/HEAD")
	if err != nil {
		fmt.Fprint(os.Stderr, "Error accessing HEAD stats\n")
		os.Exit(1)
	}

	if stats.Size() > 0{
		buffer.WriteString(fmt.Sprintf("parent "))
		parent_hash, err := os.ReadFile(".gogit/HEAD")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading commit HEAD\n")
			os.Exit(1)
		}
		buffer.Write(parent_hash)
		buffer.WriteString("\n")
	}

	user, err := user.Current()
	if err != nil{
		fmt.Fprintf(os.Stderr, "Error getting author info\n")
		os.Exit(1)
	}

	name := user.Name
	username := user.Username
	hostname, err := os.Hostname()
	if err != nil{
		fmt.Fprintf(os.Stderr, "Error getting hostname\n")
		os.Exit(1)
	}

	email := fmt.Sprintf("%s@%s", username, hostname)

	timestamp := get_timestamp()

	buffer.WriteString(fmt.Sprintf("author %s <%s> %s\n", name, email, timestamp))
	buffer.WriteString(fmt.Sprintf("committer %s <%s> %s\n\n", name, email, timestamp))
	buffer.WriteString(fmt.Sprintf("%s\n", message))

	header_buffer := new(bytes.Buffer)
	header_buffer.WriteString(fmt.Sprintf("%s %d", "commit", buffer.Len()))
	header_buffer.WriteByte(0)

	final_buffer := new(bytes.Buffer)
	final_buffer.Write(header_buffer.Bytes())
	final_buffer.Write(buffer.Bytes())

	return final_buffer.Bytes()
}

func save_commit_object(commit_object []byte){
	sha1_sum := sha1.Sum(commit_object)
	sha1_hex := fmt.Sprintf("%x", sha1_sum)

	fd, err := os.OpenFile(".gogit/HEAD", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil{
		fmt.Fprintf(os.Stderr, "Error accessing commit HEAD\n")
		os.Exit(1)
	}
	defer fd.Close()

	_, err = fd.Write([]byte(sha1_hex))
	if err != nil{
		fmt.Fprintf(os.Stderr, "Error writing SHA-1 hash to HEAD\n")
		os.Exit(1)
	}

	var compressed bytes.Buffer
	zw, err := zlib.NewWriterLevel(&compressed, zlib.BestCompression)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating writer\n")
		os.Exit(1)
	}

	_, err = zw.Write(commit_object)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to writer\n")
		os.Exit(1)
	}

	if err := zw.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "Error closing writer\n")
		os.Exit(1)
	}

	compressed_bytes := compressed.Bytes()

	directory := sha1_hex[0:2]
	file := sha1_hex[2:]

	object_dir := filepath.Join(".gogit", "objects", directory)

	object_path := filepath.Join(object_dir, file)

	object_fd, err := os.Create(object_path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating object file\n")
		os.Exit(1)
	}

	_, err = object_fd.Write(compressed_bytes)
	if err != nil{
		fmt.Fprintf(os.Stderr, "Error writing to object file\n")
		os.Exit(1)
	}
}

func Commit(message string) {
	sha1_hex := write_tree()

	commit_object := create_commit_object(message, sha1_hex)

	save_commit_object(commit_object)
}