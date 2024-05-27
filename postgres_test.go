package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/ejones77/432_final_project/pkg"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestConnectToPostgres(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	os.Setenv("POSTGRES_DB", "test_db")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))
	db := pkg.ConnectToPostgres(dsn)

	assert.NotNil(t, db)

	os.Unsetenv("POSTGRES_DB")
}

func connectToTestDB(t *testing.T) *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	os.Setenv("POSTGRES_DB", "test_db")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))
	db := pkg.ConnectToPostgres(dsn)

	assert.NotNil(t, db)

	os.Unsetenv("POSTGRES_DB")

	return db
}

func TestLoadToPostgres(t *testing.T) {
	db := connectToTestDB(t)

	// create a table
	db.Exec("DROP TABLE IF EXISTS load_tests;")
	db.Exec("CREATE TABLE load_tests (id INT, type VARCHAR(255));")

	// insert data
	type LoadTest struct {
		ID   int    `db:"id"`
		Type string `db:"type"`
	}

	LoadTests := []LoadTest{
		{ID: 1, Type: "A"},
		{ID: 2, Type: "B"},
		{ID: 3, Type: "C"},
	}

	pkg.LoadToPostgres(db, LoadTests)

	// check if data was inserted
	var count int
	db.Raw("SELECT COUNT(*) FROM load_tests").Scan(&count)
	assert.Equal(t, 3, count)

	// clean up
	db.Exec("DROP TABLE load_tests")
}
