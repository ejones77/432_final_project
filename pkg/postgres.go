/*
Handling utils necessary for loading into postgres.
*/
package pkg

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strconv"

	"github.com/go-gota/gota/dataframe"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectToPostgres() *sql.DB {

	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	var (
		POSTGRES_HOST     = os.Getenv("POSTGRES_HOST")
		POSTGRES_PORT     = os.Getenv("POSTGRES_PORT")
		POSTGRES_USER     = os.Getenv("POSTGRES_USER")
		POSTGRES_PASSWORD = os.Getenv("POSTGRES_PASSWORD")
		POSTGRES_DB       = os.Getenv("POSTGRES_DB")
	)

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
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

func ConvertTypes(value interface{}, targetType reflect.Type) (interface{}, error) {
	/*
		enforces if a detected type doesn't match the target type
		starts with null value handling, then converts to the target type
	*/

	if value == nil {
		switch targetType.Kind() {
		case reflect.String:
			return "", nil
		case reflect.Int:
			return 0, nil
		case reflect.Float64:
			return float64(0), nil
		default:
			return nil, fmt.Errorf("unsupported type: %v", targetType)
		}
	}

	switch targetType.Kind() {
	case reflect.String:
		return fmt.Sprintf("%v", value), nil
	case reflect.Int:
		return strconv.Atoi(fmt.Sprintf("%v", value))
	case reflect.Float64:
		return strconv.ParseFloat(fmt.Sprintf("%v", value), 64)
	default:
		return nil, fmt.Errorf("unsupported type: %v", targetType)
	}
}

func GetFieldNames(structType interface{}) []string {
	/*
		This is to obtain the correct order of fields in the struct
	*/
	t := reflect.TypeOf(structType).Elem()
	fieldNames := make([]string, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		fieldNames[i] = t.Field(i).Tag.Get("db")
	}
	return fieldNames
}

func LoadToPostgres(df dataframe.DataFrame, db *sql.DB, numFields int, query string, structType interface{}) {
	/*
		Reorders & converts types to load into postgres
		Params:
			df: dataframe.DataFrame the dataframe to load
			db: *sql.DB the database connection
			numFields: int the number of fields in the dataframe
			query: string the query to execute
			structType: interface{} the struct to load into
	*/

	fieldNames := GetFieldNames(structType)
	df = df.Select(fieldNames)

	for i := 0; i < df.Nrow(); i++ {
		row := df.Maps()[i]
		var values []interface{}
		for j := 0; j < numFields; j++ {
			value := row[df.Names()[j]]
			field := reflect.TypeOf(structType).Elem().Field(j)
			convertedValue, err := ConvertTypes(value, field.Type)
			if err != nil {
				fmt.Println(err)
				continue
			}
			values = append(values, convertedValue)
		}
		_, err := db.Exec(query, values...)
		if err != nil {
			fmt.Println(err)
		}
	}
}
