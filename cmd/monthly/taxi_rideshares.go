package monthly

import (
	"fmt"
	"time"

	"github.com/ejones77/432_final_project/pkg"
	"github.com/go-gota/gota/dataframe"
)

type TaxiRideshares struct {
	TripID                   int       `db:"trip_id"`
	TaxiID                   string    `db:"taxi_id"`
	TripStartTimestamp       time.Time `db:"trip_start_timestamp"`
	TripEndTimestamp         time.Time `db:"trip_end_timestamp"`
	TripSeconds              int       `db:"trip_seconds"`
	TripMiles                float64   `db:"trip_miles"`
	Fare                     float64   `db:"fare"`
	Tip                      float64   `db:"tip"`
	AdditionalCharges        float64   `db:"additional_charges"`
	TripTotal                float64   `db:"trip_total"`
	PickupCentroidLatitude   float64   `db:"pickup_centroid_latitude"`
	PickupCentroidLongitude  float64   `db:"pickup_centroid_longitude"`
	DropoffCentroidLatitude  float64   `db:"dropoff_centroid_latitude"`
	DropoffCentroidLongitude float64   `db:"dropoff_centroid_longitude"`
	PickupZipCode            string    `db:"pickup_zip_code"`
	DropoffZipCode           string    `db:"dropoff_zip_code"`
	PickupCommunityArea      string    `db:"pickup_community_area"`
	DropoffCommunityArea     string    `db:"dropoff_community_area"`
}

func ExtractTaxis() []map[string]interface{} {
	columns := []string{
		"trip_id",
		"taxi_id",
		"trip_start_timestamp",
		"trip_end_timestamp",
		"trip_seconds",
		"trip_miles",
		"fare",
		"tips",
		"extras",
		"trip_total",
		"pickup_centroid_latitude",
		"pickup_centroid_longitude",
		"dropoff_centroid_latitude",
		"dropoff_centroid_longitude",
	}

	data, err := pkg.QuerySample("wrvz-psew",
		"trip_start_timestamp",
		columns,
		`trip_start_timestamp >= '2020-04-01' 
		AND trip_start_timestamp < '2020-04-02'`,
		10,
	)

	if err != nil {
		fmt.Println(err)
	}

	return data
}

func ExtractRideshares() []map[string]interface{} {
	columns := []string{
		"trip_id",
		"trip_start_timestamp",
		"trip_end_timestamp",
		"trip_seconds",
		"trip_miles",
		"fare",
		"tip",
		"additional_charges",
		"trip_total",
		"pickup_centroid_latitude",
		"pickup_centroid_longitude",
		"dropoff_centroid_latitude",
		"dropoff_centroid_longitude",
	}

	data, err := pkg.QuerySample("m6dm-c72p",
		"trip_start_timestamp",
		columns,
		`trip_start_timestamp >= '2020-04-01' 
		AND trip_start_timestamp < '2020-04-02'`,
		10,
	)
	if err != nil {
		fmt.Println(err)
	}

	return data

}

func TransformTaxiRideshares() dataframe.DataFrame {
	/*
		taxi id as null on rideshare
		tips on taxis -> tip
		extras on taxis -> additional_charges

		then it's a simple union all
	*/
	taxiData := ExtractTaxis()
	rideshareData := ExtractRideshares()
	rideshareData = append(rideshareData, map[string]interface{}{"taxi_id": nil})

	df1 := dataframe.LoadMaps(taxiData)
	df2 := dataframe.LoadMaps(rideshareData)

	df1 = df1.Rename("tip", "tips")
	df1 = df1.Rename("additional_charges", "extras")

	df2 = df2.Select(df1.Names())

	merged := df1.RBind(df2)

	return merged
}
