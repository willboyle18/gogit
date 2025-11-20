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
	StDev uint32
	StIno uint32
	StMode uint32
	StUid uint32
	StGid uint32
	StSize uint32
	Sha1 [20]byte
	Namelen uint16
	Name string
}

func Cache_Entry_Size(filename_length int) int {
	fixed_size := 8 + 8 + 4 + 4 + 4 + 4 + 4 + 4 + 20 + 2
	raw_entry_size := fixed_size + filename_length
	cache_entry_size := (raw_entry_size + 8) &^ 7
	return cache_entry_size
}

var ActiveCache []*CacheEntry
var ActiveNR uint32
var ActiveAlloc uint32
