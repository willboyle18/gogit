package repo

import (
	"fmt"
	"log"
	"os"
)

func Init() {
	err := os.Mkdir(".gogit", 0700)
	if err != nil && !os.IsExist(err) {
		fmt.Fprintln(os.Stderr, "unable to create .gogit")
		log.Fatal(err)
		os.Exit(1)
	}

	var sha1_dir string = ".gogit/objects"
	err = os.Mkdir(sha1_dir, 0700)
	if err != nil && !os.IsExist(err){
		fmt.Fprintf(os.Stderr, sha1_dir)
		log.Fatal(err)
		os.Exit(1)
	}

	for i := 0; i < 256; i++ {
		hex_digits := fmt.Sprintf("%02x", i)
		path := sha1_dir + "/" + hex_digits
		err := os.Mkdir(path, 0700)
		if err != nil && !os.IsExist(err){
			fmt.Fprintf(os.Stderr, path)
			log.Fatal(err)
			os.Exit(1)
		}
	}

	fmt.Println("gogit repo intialized successfully")
}