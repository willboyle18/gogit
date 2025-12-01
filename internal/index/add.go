package index

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"syscall"

	"github.com/willboyle18/gogit/internal/cache"
)


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
	pos := cache_name_pos(cache_entry.Name)

	if pos < 0 {
		cache.ActiveCache[-pos-1] = cache_entry
		return true
	}

	cache.ActiveNR++
	cache.ActiveCache = append(cache.ActiveCache, nil)
	copy(cache.ActiveCache[pos+1:], cache.ActiveCache[pos:])
	cache.ActiveCache[pos] = cache_entry
	return true
}


func write_cache(new_fd *os.File){
	var cache_header cache.CacheHeader

	cache_header.Signature = cache.CACHE_SIGNATURE
	cache_header.Version = 1
	cache_header.Entries = uint32(cache.ActiveNR)

	// Initialize the SHA1 hasher
	hasher := sha1.New()
	buffer := new(bytes.Buffer)

	// Write the first 12 header bytes to the hasher
	err := binary.Write(buffer, binary.BigEndian, cache_header.Signature)
	if err != nil {
		log.Fatal(err)
	}
	err = binary.Write(buffer, binary.BigEndian, cache_header.Version)
	if err != nil {
		log.Fatal(err)
	}
	err = binary.Write(buffer, binary.BigEndian, cache_header.Entries)
	if err != nil {
		log.Fatal(err)
	}
	_, err = hasher.Write(buffer.Bytes())
	if err != nil{
		log.Fatal(err)
	}

	// Add entires to SHA1 hasher
	for i := 0; i < int(cache.ActiveNR); i++ {
		raw_entry_bytes := turn_cache_entry_into_raw_bytes(cache.ActiveCache[i])
		_, err = hasher.Write(raw_entry_bytes)
		if err != nil{
			log.Fatal(err)
		}
	}

	// Compute SHA1
	final := hasher.Sum(nil)
	copy(cache_header.Sha1[:], final)

	// Write SHA1 to the buffer
	err = binary.Write(buffer, binary.BigEndian, cache_header.Sha1[:])
	if err != nil {
		log.Fatal(err)
	}

	// Write header to the index file
	_, err = new_fd.Write(buffer.Bytes())
	if err != nil{
		log.Fatal(err)
	}

	// Write file entry bytes to the index
	for i := 0; i < int(cache.ActiveNR); i++{
		raw_entry_bytes := turn_cache_entry_into_raw_bytes(cache.ActiveCache[i])
		_, err = new_fd.Write(raw_entry_bytes)
		if err != nil{
			log.Fatal(err)
		}
	}
}

func remove_file_from_cache(path string){
	pos := cache_name_pos(path)

	// Deleted file needs to be removed from the active cache
	if pos < 0 {
		cache.ActiveCache = append(cache.ActiveCache[:-pos-1], cache.ActiveCache[-pos:]...)
		fmt.Println("removed", path)
		cache.ActiveNR--
		return
	}
	fmt.Println("file does not exist")
}

func add_file_to_cache(path string) bool {
	// Block 1: Open the file
	fd, err := os.Open(path)
	if os.IsNotExist(err){
		remove_file_from_cache(path)
		return true
	} else if err != nil {
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
	cache_entry.Name_Length = uint16(len(path))

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
		return
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
	defer new_fd.Close()

	// Block 3: Loop over all paths passed on the command line
	for i := 2; i < len(args); i++ {
		// Block 4: Verify the path
		path := args[i]
		verified := verify_path(path) // located in add_helpers.go
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
	write_cache(new_fd)
	err = os.Rename(".gogit/index.lock", ".gogit/index")
	if err != nil{
		log.Fatal(err)
	}
}
