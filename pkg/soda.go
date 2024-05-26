/*
Used to extract data from an endpoint
*/
package pkg

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/SebastiaanKlippert/go-soda"
)

func QuerySample(apiEndpoint string,
	orderColumn string,
	selectClause []string,
	whereClause string,
	limit uint,
	v interface{}) error {
	/*
		Define the query to pull a specific set of data
		Params:
			apiEndpoint: the endpoint to pull data from
			orderColumn: the column to order the data by

			selectClause: []string the columns to select
			ex: []string{"farm_name", "category", "item", "zipcode"}

			whereClause: string how to filter the data
			ex: `lower(farm_name) like '%sun%farm%' AND (item in('Radishes',
			'Cucumbers') OR lower(item) like '%flower%')`
	*/

	url := fmt.Sprintf("https://data.cityofchicago.org/resource/%s", apiEndpoint)
	sodareq := soda.NewGetRequest(url, "")

	sodareq.Format = "json"
	sodareq.Query.Select = selectClause
	sodareq.Query.Where = whereClause
	sodareq.Query.Limit = limit
	sodareq.Query.AddOrder(orderColumn, soda.DirAsc)

	resp, err := sodareq.Get()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return err
	}

	return nil
}

func ConcurrentQuerySample(apiEndpoint string,
	selectClause []string,
	whereClause string,
	concurrency int,
	pageSize uint,
	v interface{}) error {

	url := fmt.Sprintf("https://data.cityofchicago.org/resource/%s", apiEndpoint)
	gr := soda.NewGetRequest(url, "")
	gr.Format = "json"
	gr.Query.Select = selectClause
	gr.Query.Where = whereClause
	gr.Query.AddOrder("trip_start_timestamp", soda.DirAsc)

	ogr, err := soda.NewOffsetGetRequest(gr)
	if err != nil {
		return err
	}

	resultsVal := reflect.ValueOf(v).Elem()

	for i := 0; i < concurrency; i++ {
		ogr.Add(1)
		go func() {
			defer ogr.Done()

			for {
				resp, err := ogr.Next(pageSize)
				if err == soda.ErrDone {
					break
				}
				if err != nil {
					log.Fatal(err)
				}

				// Create a new slice to hold this batch of results
				sliceType := reflect.SliceOf(resultsVal.Type().Elem())
				data := reflect.New(sliceType).Interface()

				err = json.NewDecoder(resp.Body).Decode(data)
				resp.Body.Close()
				if err != nil {
					log.Fatal(err)
				}

				// Append the data to the results
				resultsVal.Set(reflect.AppendSlice(resultsVal, reflect.ValueOf(data).Elem()))
			}
		}()
	}

	ogr.Wait()

	return nil
}
