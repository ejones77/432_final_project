package once

import (
	"fmt"

	"github.com/ejones77/432_final_project/pkg"
	"github.com/go-gota/gota/dataframe"
)

type Geographies struct {
	GeographyType                    string  `db:"geography_type"`
	CommunityAreaOrZip               string  `db:"community_area_or_zip"`
	CommunityAreaName                string  `db:"community_area_name"`
	CcviScore                        float64 `db:"ccvi_score"`
	CcviCategory                     string  `db:"ccvi_category"`
	RankSocioeconomicStatus          float64 `db:"rank_socioeconomic_status"`
	RankAdultsNoPcp                  float64 `db:"rank_adults_no_pcp"`
	RankCumulativeMobilityRatio      float64 `db:"rank_cumulative_mobility_ratio"`
	RankFrontlineEssentialWorkers    float64 `db:"rank_frontline_essential_workers"`
	RankAge65Plus                    float64 `db:"rank_age_65_plus"`
	RankComorbidConditions           float64 `db:"rank_comorbid_conditions"`
	RankCovid19IncidenceRate         float64 `db:"rank_covid_19_incidence_rate"`
	RankCovid19HospitalAdmissionRate float64 `db:"rank_covid_19_hospital_admission_rate"`
	RankCovid19CrudeMortalityRate    float64 `db:"rank_covid_19_crude_mortality_rate"`
	BelowPovertyLevel                float64 `db:"below_poverty_level"`
	CrowdedHousing                   float64 `db:"crowded_housing"`
	NoHighSchoolDiploma              float64 `db:"no_high_school_diploma"`
	PerCapitaIncome                  float64 `db:"per_capita_income"`
	Unemployment                     float64 `db:"unemployment"`
}

func extractPubHealth() []map[string]interface{} {
	columns := []string{
		"community_area",
		"below_poverty_level",
		"crowded_housing",
		"no_high_school_diploma",
		"per_capita_income",
		"unemployment",
	}

	data, err := pkg.QuerySample("iqnk-2tcu",
		"community_area",
		columns,
		"",
		200)
	if err != nil {
		fmt.Println(err)
	}

	return data
}

func extractCCVI() []map[string]interface{} {

	columns := []string{
		"geography_type",
		"community_area_or_zip",
		"community_area_name",
		"ccvi_score",
		"ccvi_category",
		"rank_socioeconomic_status",
		"rank_adults_no_pcp",
		"rank_cumulative_mobility_ratio",
		"rank_frontline_essential_workers",
		"rank_age_65_plus",
		"rank_comorbid_conditions",
		"rank_covid_19_incidence_rate",
		"rank_covid_19_hospital_admission_rate",
		"rank_covid_19_crude_mortality_rate",
	}

	data, err := pkg.QuerySample("xhc6-88s9",
		"geography_type",
		columns,
		"",
		200)

	if err != nil {
		fmt.Println(err)
	}

	return data
}

func TransformGeographies() dataframe.DataFrame {
	// Extract the data from the geographies endpoint
	pubHealthData := extractPubHealth()
	ccviData := extractCCVI()

	df1 := dataframe.LoadMaps(pubHealthData)
	df2 := dataframe.LoadMaps(ccviData)

	df1 = df1.Rename("community_area_or_zip", "community_area")

	merged := df2.LeftJoin(df1, "community_area_or_zip")

	return merged
}

func LoadGeographies() {
	db := pkg.ConnectToPostgres()
	defer db.Close()

	geographies := TransformGeographies()

	numFields := geographies.Ncol()
	placeholders := pkg.GeneratePlaceholders(numFields)
	query_string := fmt.Sprintf("INSERT INTO geographies VALUES (%s)", placeholders)
	pkg.LoadToPostgres(geographies, db, numFields, query_string, &Geographies{})
}
