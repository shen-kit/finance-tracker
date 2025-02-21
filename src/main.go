package main

import (
	"fmt"
	"os"

	"github.com/shen-kit/finance-tracker/backend"
	"github.com/shen-kit/finance-tracker/frontend"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Run using `$finance-tracker <path_to_db>`")
		os.Exit(0)
	}

	backend.SetupDb(os.Args[1])
	frontend.CreateTUI()
}
