package index

import(
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


func cache_name_compare(name1 string, name2 string) int {
	len1 := len(name1)
	len2 := len(name2)

	var min_len int
	if len1 < len2{
		min_len = len1
	} else{
		min_len = len2
	}

	for i := 0; i < min_len ; i++ {
		if name1[i] < name2[i]{
			return -1
		} else if name2[i] < name1[i]{
			return 1
		} else{
			continue
		}
	}
	if len1 < len2 {
		return -1
	} else if len2 < len1 {
		return 1
	}
	return 0
}

func cache_name_pos(name string) int {
	first := 0
	last := int(cache.ActiveNR)

	for first < last {
		next := (first + last) >> 1
		cache_entry := cache.ActiveCache[next]
		cmp := cache_name_compare(name, cache_entry.Name)
		if cmp == 0 {
			return -next-1
		}
		if cmp < 0 {
			last = next
			continue
		}
		first = next + 1
	}
	return first
}