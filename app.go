package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"

	"github.com/ejones77/432_final_project/cmd/daily"
	"github.com/ejones77/432_final_project/cmd/monthly"
	"github.com/ejones77/432_final_project/cmd/once"
	"github.com/ejones77/432_final_project/cmd/weekly"
	"github.com/ejones77/432_final_project/pkg"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func Retry(attempts int, sleep time.Duration, fn func() error) error {
	err := fn()
	if err != nil {
		if attempts--; attempts > 0 {
			// Exponential backoff
			time.Sleep(sleep)
			return Retry(attempts, 2*sleep, fn)
		}
	}
	return err
}

func cronJob(c *cron.Cron, db *gorm.DB, jobtype string, job func(*gorm.DB) error) error {
	_, err := c.AddFunc(jobtype, func() {
		fmt.Println("Starting task:", jobtype)
		err := Retry(3, 1*time.Second, func() error {
			return job(db)
		})
		if err != nil {
			fmt.Println("Error running task:", err)
		} else {
			fmt.Println("Successfully completed task:", jobtype)
		}
	})
	return err
}

func getSecret(secretID string) (map[string]string, error) {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/final-project-424101/secrets/%s/versions/latest", secretID),
	}

	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return nil, err
	}

	// The secret value needs to be parsed as JSON.
	var secretValues map[string]string
	err = json.Unmarshal(result.Payload.Data, &secretValues)
	if err != nil {
		return nil, err
	}

	return secretValues, nil
}

func main() {

	secrets, err := getSecret("POSTGRES_SECRETS")
	if err != nil {
		log.Fatalf("Failed to get secret: %v", err)
	}

	dbname := secrets["POSTGRES_DB"]
	host := secrets["POSTGRES_HOST"]
	user := secrets["POSTGRES_USER"]
	password := secrets["POSTGRES_PASSWORD"]
	port := secrets["POSTGRES_PORT"]

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		user,
		password,
		dbname,
	)
	db := pkg.ConnectToPostgres(dsn)
	c := cron.New()

	// rune the once job if the table is empty
	if pkg.IsEmpty(db, "geographies") {
		err = once.LoadGeographies(db)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Run the daily job immediately on startup
	err = daily.LoadBuildingPermits(db)
	if err != nil {
		log.Fatal(err)
	}

	err = cronJob(c, db, "@daily", func(db *gorm.DB) error {
		if pkg.IsEmpty(db, "building_permits") {
			return weekly.LoadCovid(db)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(10 * time.Second)

	if pkg.IsEmpty(db, "covid_cases") {
		err = weekly.LoadCovid(db)
		if err != nil {
			log.Fatal(err)
		}
	}
	err = cronJob(c, db, "@weekly", func(db *gorm.DB) error {
		if pkg.IsEmpty(db, "covid_cases") {
			return weekly.LoadCovid(db)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(10 * time.Second)

	if pkg.IsEmpty(db, "taxi_rideshares") {
		err = monthly.LoadTaxiRideshares(db)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = cronJob(c, db, "@monthly", func(db *gorm.DB) error {
		if pkg.IsEmpty(db, "taxi_rideshares") {
			return monthly.LoadTaxiRideshares(db)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	c.Start()

	select {}
}
