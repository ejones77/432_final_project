package monthly

import (
	"fmt"

	"github.com/ejones77/432_final_project/pkg"
	"gorm.io/gorm"
)

type Taxis struct {
	TripID                   string            `json:"trip_id"`
	TaxiID                   string            `json:"taxi_id"`
	TripStartTimestamp       pkg.CustomTime    `json:"trip_start_timestamp"`
	TripEndTimestamp         pkg.CustomTime    `json:"trip_end_timestamp"`
	TripSeconds              pkg.Float64String `json:"trip_seconds"`
	TripMiles                pkg.Float64String `json:"trip_miles"`
	Fare                     pkg.Float64String `json:"fare"`
	Tips                     pkg.Float64String `json:"tips"`
	Extras                   pkg.Float64String `json:"extras"`
	TripTotal                pkg.Float64String `json:"trip_total"`
	PickupCentroidLatitude   pkg.Float64String `json:"pickup_centroid_latitude"`
	PickupCentroidLongitude  pkg.Float64String `json:"pickup_centroid_longitude"`
	DropoffCentroidLatitude  pkg.Float64String `json:"dropoff_centroid_latitude"`
	DropoffCentroidLongitude pkg.Float64String `json:"dropoff_centroid_longitude"`
}

type Rideshares struct {
	TripID                   string            `json:"trip_id"`
	TripStartTimestamp       pkg.CustomTime    `json:"trip_start_timestamp"`
	TripEndTimestamp         pkg.CustomTime    `json:"trip_end_timestamp"`
	TripSeconds              pkg.Float64String `json:"trip_seconds"`
	TripMiles                pkg.Float64String `json:"trip_miles"`
	Fare                     pkg.Float64String `json:"fare"`
	Tip                      pkg.Float64String `json:"tip"`
	AdditionalCharges        pkg.Float64String `json:"additional_charges"`
	TripTotal                pkg.Float64String `json:"trip_total"`
	PickupCentroidLatitude   pkg.Float64String `json:"pickup_centroid_latitude"`
	PickupCentroidLongitude  pkg.Float64String `json:"pickup_centroid_longitude"`
	DropoffCentroidLatitude  pkg.Float64String `json:"dropoff_centroid_latitude"`
	DropoffCentroidLongitude pkg.Float64String `json:"dropoff_centroid_longitude"`
}

type TaxiRideshares struct {
	TripID                   string            `db:"trip_id"`
	TaxiID                   string            `db:"taxi_id"`
	TripStartTimestamp       pkg.CustomTime    `db:"trip_start_timestamp"`
	TripEndTimestamp         pkg.CustomTime    `db:"trip_end_timestamp"`
	TripSeconds              pkg.Float64String `db:"trip_seconds"`
	TripMiles                pkg.Float64String `db:"trip_miles"`
	Fare                     pkg.Float64String `db:"fare"`
	Tip                      pkg.Float64String `db:"tip"`
	AdditionalCharges        pkg.Float64String `db:"additional_charges"`
	TripTotal                pkg.Float64String `db:"trip_total"`
	PickupCentroidLatitude   pkg.Float64String `db:"pickup_centroid_latitude"`
	PickupCentroidLongitude  pkg.Float64String `db:"pickup_centroid_longitude"`
	DropoffCentroidLatitude  pkg.Float64String `db:"dropoff_centroid_latitude"`
	DropoffCentroidLongitude pkg.Float64String `db:"dropoff_centroid_longitude"`
}

func ExtractTaxis() ([]Taxis, error) {
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
	var results []Taxis
	err := pkg.ConcurrentQuerySample("wrvz-psew",
		columns,
		`trip_start_timestamp >= '2020-04-01' 
		AND trip_start_timestamp < '2020-05-01'`,
		4,
		2000,
		&results)

	if err != nil {
		fmt.Println(err)
	}

	return results, err
}

func ExtractRideshares() ([]Rideshares, error) {
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

	var results []Rideshares
	err := pkg.ConcurrentQuerySample("m6dm-c72p",
		columns,
		`trip_start_timestamp >= '2020-04-01' 
		AND trip_start_timestamp < '2020-05-01'`,
		4,
		2000,
		&results)

	if err != nil {
		fmt.Println(err)
	}

	return results, err
}

func TransformTaxiRideshares() ([]TaxiRideshares, error) {
	/*
		taxi id as null on rideshare
		tips on taxis -> tip
		extras on taxis -> additional_charges

		then it's a simple union all
	*/
	fmt.Println("Extracting data from the taxi and rideshare endpoints")

	taxiData, err := ExtractTaxis()
	if err != nil {
		return nil, err
	}

	rideshareData, err := ExtractRideshares()
	if err != nil {
		return nil, err
	}

	fmt.Println("Data extracted successfully")

	var merged []TaxiRideshares

	// Transform Taxis data
	for _, taxi := range taxiData {
		merged = append(merged, TaxiRideshares{
			TripID:                   taxi.TripID,
			TaxiID:                   taxi.TaxiID,
			TripStartTimestamp:       taxi.TripStartTimestamp,
			TripEndTimestamp:         taxi.TripEndTimestamp,
			TripSeconds:              taxi.TripSeconds,
			TripMiles:                taxi.TripMiles,
			Fare:                     taxi.Fare,
			Tip:                      taxi.Tips,
			AdditionalCharges:        taxi.Extras,
			TripTotal:                taxi.TripTotal,
			PickupCentroidLatitude:   taxi.PickupCentroidLatitude,
			PickupCentroidLongitude:  taxi.PickupCentroidLongitude,
			DropoffCentroidLatitude:  taxi.DropoffCentroidLatitude,
			DropoffCentroidLongitude: taxi.DropoffCentroidLongitude,
		})
	}

	// Transform Rideshares data
	for _, rideshare := range rideshareData {
		merged = append(merged, TaxiRideshares{
			TripID:                   rideshare.TripID,
			TaxiID:                   "",
			TripStartTimestamp:       rideshare.TripStartTimestamp,
			TripEndTimestamp:         rideshare.TripEndTimestamp,
			TripSeconds:              rideshare.TripSeconds,
			TripMiles:                rideshare.TripMiles,
			Fare:                     rideshare.Fare,
			Tip:                      rideshare.Tip,
			AdditionalCharges:        rideshare.AdditionalCharges,
			TripTotal:                rideshare.TripTotal,
			PickupCentroidLatitude:   rideshare.PickupCentroidLatitude,
			PickupCentroidLongitude:  rideshare.PickupCentroidLongitude,
			DropoffCentroidLatitude:  rideshare.DropoffCentroidLatitude,
			DropoffCentroidLongitude: rideshare.DropoffCentroidLongitude,
		})
	}

	return merged, nil
}

func LoadTaxiRideshares(db *gorm.DB) {

	data, err := TransformTaxiRideshares()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Data transformed successfully")

	pkg.LoadToPostgres(db, data)

	fmt.Println("Data loaded successfully")
}
