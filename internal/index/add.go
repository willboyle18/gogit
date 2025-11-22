package index

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"syscall"
	"path/filepath"

	"github.com/willboyle18/gogit/internal/cache"
)

func verify_path(path string) bool {
	if path[0] == '/' {
		return false
	}
	i := 0
	length := len(path)
	current_path := ""
	for i < length {
		if path[i] == '/' {
			if current_path == ".gogit" {
				return false
			}
			if current_path == ".." {
				return false
			}
			current_path = ""
			i += 1
			if i == length {
				return false
			}
		}
		current_path = current_path + string(path[i])
		i += 1
	}
	return true
}

func index_fd(path string, name_length int, cache_entry *cache.CacheEntry, fd *os.File, stats *syscall.Stat_t) int {
	data, err := io.ReadAll(fd) // Reads all the data into a buffer
	if err != nil {
		return -1
	}
	fd.Close()

	header := fmt.Sprintf("blob %d\x00", len(data)) // Build the Git blob header: "blob <size>\0" (required for hashing)
	header_bytes := []byte(header) // Convert header string into raw bytes so we can concatenate it with the file contents

	object := append(header_bytes, data...) // Build the complete uncompressed Git blob object: header + file contents

	sha1_sum := sha1.Sum(object) // Compute the sha1 of the uncompressed object
	copy(cache_entry.Sha1[:], sha1_sum[:]) // Store the blob’s SHA-1 hash in the cache entry (used by the index)

	// Compression step
	var compressed bytes.Buffer // Destination buffer (Dynamically growing buffer that will hold the compressed blob object)
	zw, err := zlib.NewWriterLevel(&compressed, zlib.BestCompression) // zlib compresser that will write its compressed output into the compressed buffer
	if err != nil {
		return -1
	}

	_, err = zw.Write(object) // Compress the uncompressed blob and writes the compressed blob to 'compressed'
	if err != nil {
		return -1
	}

	// Finalize compression
	if err := zw.Close(); err != nil {
		return -1
	}

	compressed_bytes := compressed.Bytes()

	shaHex := fmt.Sprintf("%x", sha1_sum)
	dir := shaHex[:2]
	file := shaHex[2:] 
	
	object_dir := filepath.Join(".gogit", "objects", dir)

	fmt.Println(object_dir)

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

func add_cache_entry(cache_entry *cache.CacheEntry) bool {
	// TODO: implement later
	cache.ActiveNR += 1
	cache.ActiveCache = append(cache.ActiveCache, cache_entry)
	return true
}

func write_cache(new_fd *os.File, active_cache []*cache.CacheEntry, entries int){

}

func add_file_to_cache(path string) bool {
	// Block 1: Open the file
	fd, err := os.Open(path)
	if err != nil {
		// In future write a method to remove from the cache remove_file_from_cache(path)
		log.Fatal(err)
	}
	defer fd.Close()

	// Block 2: Stat the file
	info, err := fd.Stat()
	if err != nil {
		log.Fatal(err)
	}

	// Block 3: Allocate a cache_entry struct
	name_length := len(path)
	// size := cache.Cache_Entry_Size(name_length)
	cache_entry := &cache.CacheEntry{}
	cache_entry.Name = path

	// Block 4: Fill metadata
	stats := info.Sys().(*syscall.Stat_t)
	cache_entry.Ctime.Sec = uint32(stats.Ctim.Sec)
	cache_entry.Ctime.Nsec = uint32(stats.Ctim.Nsec)
	cache_entry.Mtime.Sec = uint32(stats.Mtim.Sec)
	cache_entry.Mtime.Nsec = uint32(stats.Mtim.Nsec)
	cache_entry.Dev = uint32(stats.Dev)
	cache_entry.Ino = uint32(stats.Ino)
	cache_entry.Mode = uint32(stats.Mode)
	cache_entry.Uid = uint32(stats.Uid)
	cache_entry.Gid = uint32(stats.Gid)
	cache_entry.Size = uint32(stats.Size)

	// Block 5: Process file contents, compute SHA-1, write blob object
	if index_fd(path, name_length, cache_entry, fd, stats) < 0 {
		return false
	}

	// Block 6: Insert cache entry into the in-memory index
	return add_cache_entry(cache_entry)
}

func Add(args []string) {
	fmt.Println(args)

	// Block 1: Load the existing index (we can skip for now because we are assuming the cache is empty for now)
	entries := cache.Read_Cache()
	if entries < 0 {
		fmt.Fprintf(os.Stderr, "Cache currupted")
	}

	// Block 2: Create .gogit/index.lock
	new_fd, err := os.Create(".gogit/index.lock")
	fmt.Println(new_fd)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.Chmod(".gogit/index.lock", 0600); err != nil {
		log.Fatal(err)
	}

	// Block 3: Loop over all paths passed on the command line
	for i := 2; i < len(args); i++ {
		// Block 4: Verify the path
		path := args[i]
		verified := verify_path(path)
		if !verified {
			fmt.Fprintf(os.Stderr, "ignoring path " + path +"\n")
			continue
		}

		// Block 5: Add the file to the index
		if !add_file_to_cache(path) {
			fmt.Fprintf(os.Stderr, "unable to add " + path + " to database\n")
			os.Remove(".gogit/index.lock")
			return
		}

	}
	// Block 6: Write the new index
	write_cache(new_fd, cache.ActiveCache, cache.ActiveNR)
	err = os.Rename(".gogit/index.lock", ".gogit/index")
	if err != nil{
		log.Fatal(err)
	}
}
