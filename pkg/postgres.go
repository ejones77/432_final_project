/*
Handling utils necessary for loading into postgres.
*/
package pkg

import (
	"fmt"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ConnectToPostgres(dsn string) *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	fmt.Println("Successfully connected!")
	return db
}

func IsEmpty(db *gorm.DB, tableName string) bool {
	var count int64
	db.Table(tableName).Count(&count)
	return count == 0
}

func IsDupe(db *gorm.DB, tableName string, idField string, idValue interface{}) bool {
	var count int64
	db.Table(tableName).Where(idField+" = ?", idValue).Count(&count)
	return count > 0
}

func LoadToPostgres(db *gorm.DB, data interface{}) {
	batchSize := 500

	err := db.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(data, batchSize).Error
	if err != nil {
		fmt.Println(err)
	}
}
