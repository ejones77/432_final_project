/*
Used to extract data from an endpoint
*/
package pkg

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/SebastiaanKlippert/go-soda"
)

func GetAllData(apiEndpoint string, orderColumn string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("https://data.cityofchicago.org/resource/%s", apiEndpoint)
	gr := soda.NewGetRequest(url, "")
	gr.Format = "json"
	gr.Query.AddOrder(orderColumn, soda.DirAsc)

	ogr, err := soda.NewOffsetGetRequest(gr)
	if err != nil {
		return nil, err
	}

	var allData []map[string]interface{}

	// goroutines to fetch data -- defined in documentation of go-soda
	for i := 0; i < 4; i++ {

		ogr.Add(1)
		go func() {
			defer ogr.Done()

			for {
				resp, err := ogr.Next(2000)
				if err == soda.ErrDone {
					break
				}
				if err != nil {
					log.Fatal(err)
				}

				var results []map[string]interface{}
				err = json.NewDecoder(resp.Body).Decode(&results)
				resp.Body.Close()
				if err != nil {
					log.Fatal(err)
				}
				allData = append(allData, results...)
			}
		}()

	}
	ogr.Wait()

	return allData, nil
}

func QuerySample(apiEndpoint string,
	orderColumn string,
	selectClause []string,
	whereClause string,
	limit uint) ([]map[string]interface{}, error) {
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

	// get some JSON data using a complex query
	sodareq.Format = "json"
	sodareq.Query.Select = selectClause
	sodareq.Query.Where = whereClause
	sodareq.Query.Limit = limit
	sodareq.Query.AddOrder(orderColumn, soda.DirAsc)

	// get the results
	resp, err := sodareq.Get()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var results []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}
