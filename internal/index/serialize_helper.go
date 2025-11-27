package index

import(
	"bytes"
	"encoding/binary"
	"log"
	
	"github.com/willboyle18/gogit/internal/cache"
)

func turn_cache_entry_into_raw_bytes(cache_entry *cache.CacheEntry) []byte {
	buffer := new(bytes.Buffer)
	err := binary.Write(buffer, binary.BigEndian, cache_entry.Ctime.Sec)
	if err != nil {
		log.Fatal(err)
	}
	err = binary.Write(buffer, binary.BigEndian, cache_entry.Ctime.Nsec)
	if err != nil {
		log.Fatal(err)
	}
	err = binary.Write(buffer, binary.BigEndian, cache_entry.Mtime.Sec)
	if err != nil {
		log.Fatal(err)
	}
	err = binary.Write(buffer, binary.BigEndian, cache_entry.Mtime.Nsec)
	if err != nil {
		log.Fatal(err)
	}
	err = binary.Write(buffer, binary.BigEndian, cache_entry.Dev)
	if err != nil {
		log.Fatal(err)
	}
	err = binary.Write(buffer, binary.BigEndian, cache_entry.Ino)
	if err != nil {
		log.Fatal(err)
	}
	err = binary.Write(buffer, binary.BigEndian, cache_entry.Mode)
	if err != nil {
		log.Fatal(err)
	}
	err = binary.Write(buffer, binary.BigEndian, cache_entry.Uid)
	if err != nil {
		log.Fatal(err)
	}
	err = binary.Write(buffer, binary.BigEndian, cache_entry.Gid)
	if err != nil {
		log.Fatal(err)
	}
	err = binary.Write(buffer, binary.BigEndian, cache_entry.Size)
	if err != nil {
		log.Fatal(err)
	}
	buffer.Write(cache_entry.Sha1[:])
	name_length := uint16(len(cache_entry.Name))
	err = binary.Write(buffer, binary.BigEndian, name_length)
	if err != nil {
		log.Fatal(err)
	}
	err = binary.Write(buffer, binary.BigEndian, []byte(cache_entry.Name))
	buffer.Write([]byte(cache_entry.Name))

	raw := buffer.Bytes()

	padding := (8 - (len(raw) % 8)) % 8

	for i := 0; i < padding; i++ {
		buffer.WriteByte(0)
	}
	return buffer.Bytes()
}