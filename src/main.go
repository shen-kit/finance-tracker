package main

import (
	"fmt"

	"github.com/shen-kit/finance-tracker/backend"
)

func main() {
	backend.SetupDb("./test.db")
	backend.CreateDummyData()
	fmt.Println("Created dummy data.")
}
