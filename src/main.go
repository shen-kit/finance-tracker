package main

import (
	"github.com/shen-kit/finance-tracker/backend"
	"github.com/shen-kit/finance-tracker/frontend"
)

func main() {
	backend.SetupDb("./test.db")
	// backend.CreateDummyData()

	frontend.CreateTUI()
}
