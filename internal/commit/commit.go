package commit

import (
	"fmt"
	"os"

	"github.com/willboyle18/gogit/internal/cache"
)

func Commit() {
	// Step 1: Load the current index into memory
	entries := cache.Read_Cache()
	if entries < 0 {
		fmt.Fprintf(os.Stderr, "Cache currupted")
		return
	}
}