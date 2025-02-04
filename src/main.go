package main

import (
	"fmt"

	"github.com/shen-kit/finance-tracker/backend"
)

func main() {
	backend.SetupDb("./test.db")
	// backend.CreateDummyData()
	// res, err := backend.GetInvestmentsRecent(0)
	res, err := backend.GetInvestmentsFilter(backend.NewFilterOpts().WithCode("VG"))
	if err != nil {
		println(err)
	}
	fmt.Printf("result length: %d\n", len(res))
	for _, r := range res {
		fmt.Printf("%+v\n", r)
	}
}
