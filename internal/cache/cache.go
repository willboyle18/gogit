package cache

import (
	// "crypto/sha1"
	// "compress/zlib"
)

const CACHE_SIGNATURE uint32 = 0x44495243 /* DIRC */
type CacheHeader struct {
	Signature uint32
	Version uint32
	Entries uint32
	Sha1 [20]byte	
}

type CacheTime struct {
	Sec uint32
	Nsec uint32
}

type CacheEntry struct {
	Ctime CacheTime
	Mtime CacheTime
	Dev uint32
	Ino uint32
	Mode uint32
	Uid uint32
	Gid uint32
	Size uint32
	Sha1 [20]byte
	Name string
}

func Cache_Entry_Size(filename_length int) int {
	fixed_size := 8 + 8 + 4 + 4 + 4 + 4 + 4 + 4 + 20
	raw_entry_size := fixed_size + filename_length
	cache_entry_size := (raw_entry_size + 7) &^ 7
	return cache_entry_size
}

var ActiveCache []*CacheEntry // correct
var ActiveNR int
var ActiveAlloc uint32
