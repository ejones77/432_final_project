package once

import (
	"fmt"
	"os"

	"github.com/ejones77/432_final_project/pkg"
	"github.com/go-gota/gota/dataframe"
)

func extractPubHealth() []map[string]interface{} {
	// Get all data from the geographies endpoint
	data, err := pkg.GetAllData("iqnk-2tcu", "community_area")
	if err != nil {
		fmt.Println(err)
	}

	return data
}

func extractCCVI() []map[string]interface{} {
	// Get all data from the geographies endpoint
	data, err := pkg.GetAllData("xhc6-88s9", "geography_type")
	if err != nil {
		fmt.Println(err)
	}

	return data
}

func TransformGeographies() {
	// Extract the data from the geographies endpoint
	pubHealthData := extractPubHealth()
	ccviData := extractCCVI()

	df1 := dataframe.LoadMaps(pubHealthData)
	df2 := dataframe.LoadMaps(ccviData)

	df1 = df1.Rename("community_area_or_zip", "community_area")

	merged := df2.LeftJoin(df1, "community_area_or_zip")

	file, err := os.Create("geographies.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	merged.WriteCSV(file)
}
