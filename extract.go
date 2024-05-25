package main

import (
	"fmt"

	"github.com/ejones77/432_final_project/pkg"
)

func main() {
	// Get all data from the geographies endpoint
	data, err := pkg.GetAllData("xhc6-88s9", "geography_type")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Print the data
	fmt.Println(data)
}
