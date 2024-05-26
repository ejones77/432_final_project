package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/ejones77/432_final_project/cmd/daily"
	"github.com/ejones77/432_final_project/pkg"
)

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	db := pkg.ConnectToPostgres()
	//once.LoadGeographies(db)
	//monthly.LoadTaxiRideshares(db)
	//weekly.LoadCovid(db)
	daily.LoadBuildingPermits(db)
}
