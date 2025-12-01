package cli

import (
	"fmt"
	"os"

	"github.com/willboyle18/gogit/internal/commit"
	"github.com/willboyle18/gogit/internal/index"
	"github.com/willboyle18/gogit/internal/repo"
)

func Run(args []string){
	fmt.Println("parsing arguments")
	fmt.Println(args)
	if args[1] == "init"{
		repo.Init()
	} else if args[1] == "add" {
		index.Add(args)
	} else if args[1] == "status" {
		index.Status()
	} else if args[1] == "commit" {
		if len(args) != 4 {
			fmt.Fprintf(os.Stderr, "Incorrect arguments\n")
			os.Exit(1)
		}
		if args[2] != "-m"{
			fmt.Fprintf(os.Stderr, "No commit message written\n")
			os.Exit(1)
		}
		if args[3] == ""{
			fmt.Fprintf(os.Stderr, "Incorrectly formatted commit message\n")
			os.Exit(1)
		}
		commit.Commit(args[3])
	} else if args[1] == "latest" {
		commit.Latest()
	}
}