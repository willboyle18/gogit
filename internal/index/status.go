package index

import(
	"fmt"
	"os"

	"github.com/willboyle18/gogit/internal/cache"
)

func Status(){
	entries := cache.Read_Cache()
	if entries < 0 {
		fmt.Fprintf(os.Stderr, "Cache currupted")
		return
	}
	if cache.ActiveNR == 0 {
		fmt.Println("No files staged for commit")
		return
	}

	fmt.Println("Files staged:")
	for i := 0; i < cache.ActiveNR; i++ {
		fmt.Printf("\t%d. %s\n", i + 1, cache.ActiveCache[i].Name)
	}
}