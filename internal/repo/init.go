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

	fmt.Println(".gogit created successfully")
}