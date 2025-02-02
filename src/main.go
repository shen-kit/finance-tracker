package main

import (
	"github.com/shen-kit/finance-tracker/backend"
)

func main() {
	backend.SetupDb("./test.db")
}
