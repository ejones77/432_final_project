package weekly

import (
	"fmt"
	"time"

	"github.com/ejones77/432_final_project/pkg"
	"gorm.io/gorm"
)

type CovidCases struct {
	ZipCode                     string            `json:"zip_code" db:"zip_code"`
	WeekStart                   pkg.CustomTime    `json:"week_start" db:"week_start"`
	WeekEnd                     pkg.CustomTime    `json:"week_end" db:"week_end"`
	CasesWeekly                 pkg.Float64String `json:"cases_weekly" db:"cases_weekly"`
	CaseRateWeekly              pkg.Float64String `json:"case_rate_weekly" db:"case_rate_weekly"`
	TestsWeekly                 pkg.Float64String `json:"tests_weekly" db:"tests_weekly"`
	PercentTestedPositiveWeekly pkg.Float64String `json:"percent_tested_positive_weekly" db:"percent_tested_positive_weekly"`
	DeathsWeekly                pkg.Float64String `json:"deaths_weekly" db:"deaths_weekly"`
	DeathRateWeekly             pkg.Float64String `json:"death_rate_weekly" db:"death_rate_weekly"`
	Population                  pkg.Float64String `json:"population" db:"population"`
}

func ExtractCovid(db *gorm.DB) ([]CovidCases, error) {

	columns := []string{
		"zip_code",
		"week_start",
		"week_end",
		"cases_weekly",
		"case_rate_weekly",
		"tests_weekly",
		"percent_tested_positive_weekly",
		"deaths_weekly",
		"death_rate_weekly",
		"population",
	}

	var results []CovidCases
	// Check if the table exists

	var startDate, endDate string
	if !pkg.IsEmpty(db, "covid_cases") {
		// Fetch the maximum date from the database
		var maxDate time.Time
		db.Table("covid_cases").Select("max(week_start)").Scan(&maxDate)

		// Calculate start and end dates based on max date
		startDate = maxDate.Format("2006-01-02")
		endDate = maxDate.AddDate(0, 0, 7).Format("2006-01-02")
	} else {
		// Use a static sample
		startDate = "2020-04-01"
		endDate = "2020-10-18"
	}

	err := pkg.QuerySample("yhhz-zm2v",
		"week_start",
		columns,
		fmt.Sprintf(`week_start >= '%s' AND week_start < '%s'`, startDate, endDate),
		2000,
		&results)

	if err != nil {
		fmt.Println(err)
	}

	return results, err
}

func LoadCovid(db *gorm.DB) error {
	data, err := ExtractCovid(db)
	if err != nil {
		return err
	}
	fmt.Println("Extracted data from Covid endpoint")

	pkg.LoadToPostgres(db, data)

	fmt.Println("Loaded data into the database")

	return nil
}
