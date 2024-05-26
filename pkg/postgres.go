/*
Handling utils necessary for loading into postgres.
*/
package pkg

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToPostgres() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	fmt.Println("Successfully connected!")
	return db
}

func GeneratePlaceholders(numFields int) string {
	/*
		Helps format the query string to not have to write $1, $2, $3, etc.
	*/
	placeholders := ""
	for i := 1; i <= numFields; i++ {
		placeholders += fmt.Sprintf("$%d", i)
		if i != numFields {
			placeholders += ", "
		}
	}
	return placeholders
}

func LoadToPostgres(db *gorm.DB, data interface{}) {
	batchSize := 500

	err := db.CreateInBatches(data, batchSize).Error
	if err != nil {
		fmt.Println(err)
	}
}
