package cli

import (
	"fmt"
	"github.com/willboyle18/gogit/internal/repo"
)

func Run(args []string){
	if args[1] == "init"{
		fmt.Println("Calling repo.init")
		repo.Init()
		fmt.Println("done");
	}
}