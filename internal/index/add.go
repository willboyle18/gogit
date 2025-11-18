package index

import (
	"fmt"
	"log"
	"os"
)

func Verify_Path(path string) bool {
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

func Add(args []string){
	fmt.Println(args)

	// Block 1: Load the existing index (we can skip for now because we are assuming the cache is empty for now)

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
		verified := Verify_Path(path)
		if !verified{
			fmt.Fprintf(os.Stderr, "ignoring path " + path)
			continue
		}

		// Block 5: Add the file to the index
	}
	// Block 6: Write the new index
	fmt.Println("finished parsing paths")

	// Block 7: Cleanup on failure
}