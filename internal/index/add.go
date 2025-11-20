package index

import (
	"fmt"
	"log"
	"os"
	"syscall"
	"github.com/willboyle18/gogit/internal/cache"
)

func verify_path(path string) bool {
	if path[0] == '/'{ return false }
	i := 0
	length := len(path)
	current_path := ""
	for i < length {
		if path[i] == '/'{
			if current_path == ".gogit"{ return false }
			if current_path == ".."{ return false }
			current_path = ""
			i += 1
			if i == length{ return false }
		}
		current_path = current_path + string(path[i])
		i += 1
	}
	return true
}

func index_fd(path string, name_length int, cache_entry *cache.CacheEntry, fd *os.File, stats *syscall.Stat_t) int {
	// TODO: implement later
	return 0
}

func add_cache_entry(cache_entry *cache.CacheEntry) bool {
	// TODO: implement later
	return true
}

func add_file_to_cache(path string) bool {
	// Block 1: Open the file
	fd, err := os.Open(path)
	if err != nil{
		// In future write a method to remove from the cache remove_file_from_cache(path)
		log.Fatal(err)
	}
	defer fd.Close()


	// Block 2: Stat the file
	info, err := fd.Stat()
	if err != nil{
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
	if(index_fd(path, name_length, cache_entry, fd, stats) < 0){
		return false
	}

	// Block 6: Insert cache entry into the in-memory index
	return add_cache_entry(cache_entry)
}

func Add(args []string){
	fmt.Println(args)

	// Block 1: Load the existing index (we can skip for now because we are assuming the cache is empty for now)
	entries := cache.Read_Cache()
	if entries < 0 {
		fmt.Fprintf(os.Stderr, "Cache currupted")
	}

	// Block 2: Create .gogit/index.lock
	new_fd, err := os.Create(".gogit/index.lock")
	fmt.Println(new_fd)
	if err != nil{
		log.Fatal(err)
	}
	if err := os.Chmod(".gogit/index.lock", 0600); err != nil {
		log.Fatal(err)
	}

	// Block 3: Loop over all paths passed on the command line
	for i := 2; i < len(args); i++{
		// Block 4: Verify the path
		path := args[i]
		verified := verify_path(path)
		if !verified{
			fmt.Fprintf(os.Stderr, "ignoring path " + path)
			continue
		}

		// Block 5: Add the file to the index
		if add_file_to_cache(path) {

		}

	}
	// Block 6: Write the new index
	fmt.Println("finished parsing paths")

	// Block 7: Cleanup on failure
}