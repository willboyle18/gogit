package commit

import (
	"fmt"
	"os"
	"bytes"
	"crypto/sha1"
	"compress/zlib"
	"path/filepath"

	"github.com/willboyle18/gogit/internal/cache"
)

func check_valid_sha1(sha1 []byte) int {
	sha1_hex := fmt.Sprintf("%x", sha1)
	directory := sha1_hex[0:2]
	file := sha1_hex[2:]
	path := filepath.Join(".gogit", "objects", directory, file)

	_, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Missing blob object\n")
		return -1
	}

	return 0
}

func store_tree_object(tree_object []byte) int {
	sha1_sum := sha1.Sum(tree_object)
	var compressed bytes.Buffer
	zw, err := zlib.NewWriterLevel(&compressed, zlib.BestCompression)
	if err != nil {
		return -1
	}

	_, err = zw.Write(tree_object)
	if err != nil {
		return -1
	}

	if err := zw.Close(); err != nil {
		return -1
	}

	compressed_bytes := compressed.Bytes()

	sha1_hex := fmt.Sprintf("%x", sha1_sum)
	directory := sha1_hex[0:2]
	file := sha1_hex[2:]

	object_dir := filepath.Join(".gogit", "objects", directory)

	err = os.MkdirAll(object_dir, 0755)
	if err != nil{
		return -1
	}

	object_path := filepath.Join(object_dir, file)

	object_fd, err := os.Create(object_path)
	if err != nil {
		return -1
	}

	_, err = object_fd.Write(compressed_bytes)
	if err != nil{
		return -1
	}


	object_fd.Close()
	return 0
}

func write_tree() {
	entries := cache.Read_Cache()
	if entries <= 0 {
		fmt.Fprintf(os.Stderr, "No file-cache to create a tree of\n")
		os.Exit(1)
	}

	entries_buffer := new(bytes.Buffer)

	for i := 0; i < cache.ActiveNR; i++ {
		cache_entry := cache.ActiveCache[i]

		if check_valid_sha1(cache_entry.Sha1[:]) < 0 {
			os.Exit(1);
		}

		entries_buffer.WriteString(fmt.Sprintf("%o %s", cache_entry.Mode, cache_entry.Name))
		entries_buffer.WriteByte(0)
		entries_buffer.Write(cache_entry.Sha1[:])
	}

	header_buffer := new(bytes.Buffer)
	header_buffer.WriteString(fmt.Sprintf("%s %d", "tree", entries_buffer.Len()))
	header_buffer.WriteByte(0)

	raw_header_buffer_bytes := header_buffer.Bytes()
	raw_entries_buffer_bytes := entries_buffer.Bytes()

	final_buffer := new(bytes.Buffer)

	final_buffer.Write(raw_header_buffer_bytes)
	final_buffer.Write(raw_entries_buffer_bytes)

	final_buffer_bytes := final_buffer.Bytes()

	store_tree_object(final_buffer_bytes)
}