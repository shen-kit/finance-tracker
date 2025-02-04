package main

import (
	"fmt"

	"github.com/shen-kit/finance-tracker/backend"
)

func main() {
	backend.SetupDb("./test.db")
	// backend.CreateDummyData()
	res, err := backend.GetInvestmentsRecent(0)
	if err != nil {
		println("error: ", err.Error())
	}
	fmt.Printf("result length: %d\n", len(res))
	for _, r := range res {
		fmt.Printf("%+v\n", r)
	}
	res, err = backend.GetInvestmentsFilter(backend.NewFilterOpts().WithCode("VG"))
	if err != nil {
		println("error: ", err.Error())
	}
	fmt.Printf("result length: %d\n", len(res))
	for _, r := range res {
		fmt.Printf("%+v\n", r)
	}

	var recRes []backend.Record
	recRes, err = backend.GetRecordsRecent(0)
	if err != nil {
		println("error: ", err.Error())
	}
	fmt.Printf("result length: %d\n", len(recRes))
	for _, r := range recRes {
		fmt.Printf("%+v\n", r)
	}
	recRes, err = backend.GetRecordsFilter(backend.NewFilterOpts().WithCatId([]int{1, 2}))
	if err != nil {
		println("error: ", err.Error())
	}
	fmt.Printf("result length: %d\n", len(recRes))
	for _, r := range recRes {
		fmt.Printf("%+v\n", r)
	}
}
