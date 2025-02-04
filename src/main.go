package main

import (
	"fmt"

	"github.com/shen-kit/finance-tracker/backend"
	"github.com/shen-kit/finance-tracker/helper"
)

func main() {
	backend.SetupDb("./test.db")
	// backend.CreateDummyData()

	// testQueryFunc(backend.GetInvestmentsRecent(0))
	// testQueryFunc(backend.GetInvestmentsFilter(backend.NewFilterOpts().WithCode("IV").WithMinCost(4000)))
	// testQueryFunc(backend.GetRecordsRecent(0))
	// testQueryFunc(backend.GetRecordsFilter(backend.NewFilterOpts().WithMinCost(0)))

	startDate, _ := helper.MakeDate(2000, 1, 1)
	endDate, _ := helper.MakeDate(3000, 1, 1)
	val, err := backend.GetIncomeSum(startDate, endDate)
	// val, err := backend.GetExpenditureSum(startDate, endDate)
	// val, err := backend.GetCategorySum(2, startDate, endDate)
	if err != nil {
		fmt.Println("error")
	}
	println(val)
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
