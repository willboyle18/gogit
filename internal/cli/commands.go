package cli

import (
	"fmt"

	"github.com/willboyle18/gogit/internal/index"
	"github.com/willboyle18/gogit/internal/repo"
)

func Run(args []string){
	fmt.Println("parsing arguments")
	if args[1] == "init"{
		repo.Init()
	} else if args[1] == "add" {
		index.Add(args)
	}
}