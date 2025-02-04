package main

import (
	"fmt"

	"github.com/shen-kit/finance-tracker/backend"
)

func main() {
	backend.SetupDb("./test.db")
	// backend.CreateDummyData()

	testQueryFunc(backend.GetInvestmentsRecent(0))
	testQueryFunc(backend.GetInvestmentsFilter(backend.NewFilterOpts().WithCode("IV").WithMinCost(4000)))
	testQueryFunc(backend.GetRecordsRecent(0))
	testQueryFunc(backend.GetRecordsFilter(backend.NewFilterOpts().WithMinCost(0)))
}

func testQueryFunc[T any](res []T, err error) {
	fmt.Printf("\n")
	if err != nil {
		println("error: ", err.Error())
	}
	fmt.Printf("result length: %d\n", len(res))
	for _, r := range res {
		fmt.Printf("%+v\n", r)
	}
}
