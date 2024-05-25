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

type Data map[string]interface{}

func GetAllData(apiEndpoint string, orderColumn string) ([]Data, error) {
	url := fmt.Sprintf("https://data.cityofchicago.org/resource/%s", apiEndpoint)
	gr := soda.NewGetRequest(url, "")
	gr.Format = "json"
	gr.Query.AddOrder(orderColumn, soda.DirAsc)

	ogr, err := soda.NewOffsetGetRequest(gr)
	if err != nil {
		return nil, err
	}

	var allData []Data

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

				var results []Data
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

func QuerySample(apiEndpoint string, orderColumn string, selectClause []string, whereClause string) ([]Data, error) {
	/*
		Define the query to pull a specific set of data
		Params:
			apiEndpoint: the endpoint to pull data from
			orderColumn: the column to order the data by
			selectClause: []string the columns to select
			ex: []string{"farm_name", "category", "item", "zipcode"}
			whereClause: string how to filter the data
	*/

	url := fmt.Sprintf("https://data.cityofchicago.org/resource/%s", apiEndpoint)
	sodareq := soda.NewGetRequest(url, "")

	// count all records
	count, err := sodareq.Count()
	if err != nil {
		return nil, err
	}
	fmt.Println(count)

	// get dataset last updated time
	modified, err := sodareq.Modified()
	if err != nil {
		return nil, err
	}
	fmt.Println(modified)

	// list all fields/columns
	fields, err := sodareq.Fields()
	if err != nil {
		return nil, err
	}
	fmt.Println(fields)

	// get some JSON data using a complex query
	sodareq.Format = "json"
	sodareq.Query.Select = selectClause
	sodareq.Query.Where = whereClause
	sodareq.Query.AddOrder(orderColumn, soda.DirAsc)

	// count this result first
	querycount, err := sodareq.Count()
	if err != nil {
		return nil, err
	}
	fmt.Println(querycount)

	// get the results
	resp, err := sodareq.Get()
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var results []Data
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}
