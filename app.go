package main

import (
	"fmt"

	"github.com/ejones77/432_final_project/cmd/monthly"
)

func main() {
	fmt.Println("---------------------------------")
	data := monthly.TransformTaxiRideshares()
	fmt.Println(data.Describe())
}
