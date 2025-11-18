package cache

import (
	"crypto/sha1"
	"compress/zlib"
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