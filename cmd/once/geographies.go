package once

import (
	"fmt"

	"github.com/ejones77/432_final_project/pkg"
	"gorm.io/gorm"
)

type PubHealth struct {
	CommunityArea       string            `json:"community_area"`
	BelowPovertyLevel   pkg.Float64String `json:"below_poverty_level"`
	CrowdedHousing      pkg.Float64String `json:"crowded_housing"`
	NoHighSchoolDiploma pkg.Float64String `json:"no_high_school_diploma"`
	PerCapitaIncome     pkg.Float64String `json:"per_capita_income"`
	Unemployment        pkg.Float64String `json:"unemployment"`
}

type CCVI struct {
	GeographyType                    string            `json:"geography_type"`
	CommunityAreaOrZip               pkg.Float64String `json:"community_area_or_zip"`
	CommunityAreaName                string            `json:"community_area_name"`
	CcviScore                        pkg.Float64String `json:"ccvi_score"`
	CcviCategory                     string            `json:"ccvi_category"`
	RankSocioeconomicStatus          pkg.Float64String `json:"rank_socioeconomic_status"`
	RankAdultsNoPcp                  pkg.Float64String `json:"rank_adults_no_pcp"`
	RankCumulativeMobilityRatio      pkg.Float64String `json:"rank_cumulative_mobility_ratio"`
	RankFrontlineEssentialWorkers    pkg.Float64String `json:"rank_frontline_essential_workers"`
	RankAge65Plus                    pkg.Float64String `json:"rank_age_65_plus"`
	RankComorbidConditions           pkg.Float64String `json:"rank_comorbid_conditions"`
	RankCovid19IncidenceRate         pkg.Float64String `json:"rank_covid_19_incidence_rate"`
	RankCovid19HospitalAdmissionRate pkg.Float64String `json:"rank_covid_19_hospital_admission_rate"`
	RankCovid19CrudeMortalityRate    pkg.Float64String `json:"rank_covid_19_crude_mortality_rate"`
}

type Geographies struct {
	GeographyType                    string            `db:"geography_type"`
	CommunityAreaOrZip               string            `db:"community_area_or_zip"`
	CommunityAreaName                string            `db:"community_area_name"`
	CcviScore                        pkg.Float64String `db:"ccvi_score"`
	CcviCategory                     string            `db:"ccvi_category"`
	RankSocioeconomicStatus          pkg.Float64String `db:"rank_socioeconomic_status"`
	RankAdultsNoPcp                  pkg.Float64String `db:"rank_adults_no_pcp"`
	RankCumulativeMobilityRatio      pkg.Float64String `db:"rank_cumulative_mobility_ratio"`
	RankFrontlineEssentialWorkers    pkg.Float64String `db:"rank_frontline_essential_workers"`
	RankAge65Plus                    pkg.Float64String `db:"rank_age_65_plus"`
	RankComorbidConditions           pkg.Float64String `db:"rank_comorbid_conditions"`
	RankCovid19IncidenceRate         pkg.Float64String `db:"rank_covid_19_incidence_rate"`
	RankCovid19HospitalAdmissionRate pkg.Float64String `db:"rank_covid_19_hospital_admission_rate"`
	RankCovid19CrudeMortalityRate    pkg.Float64String `db:"rank_covid_19_crude_mortality_rate"`
	BelowPovertyLevel                pkg.Float64String `db:"below_poverty_level"`
	CrowdedHousing                   pkg.Float64String `db:"crowded_housing"`
	NoHighSchoolDiploma              pkg.Float64String `db:"no_high_school_diploma"`
	PerCapitaIncome                  pkg.Float64String `db:"per_capita_income"`
	Unemployment                     pkg.Float64String `db:"unemployment"`
}

func extractPubHealth() ([]PubHealth, error) {
	columns := []string{
		"community_area",
		"below_poverty_level",
		"crowded_housing",
		"no_high_school_diploma",
		"per_capita_income",
		"unemployment",
	}

	var data []PubHealth
	err := pkg.QuerySample("iqnk-2tcu", "community_area", columns, "", 200, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func extractCCVI() ([]CCVI, error) {
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

	var data []CCVI
	err := pkg.QuerySample("xhc6-88s9", "geography_type", columns, "", 200, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func TransformGeographies() ([]Geographies, error) {
	// Extract the data from the geographies endpoint
	fmt.Println("Extracting data from the geographies endpoint")
	pubHealthData, err := extractPubHealth()
	if err != nil {
		return nil, err
	}

	ccviData, err := extractCCVI()
	if err != nil {
		return nil, err
	}

	fmt.Println("Data extracted successfully")

	pubHealthMap := make(map[string]PubHealth)
	for _, ph := range pubHealthData {
		pubHealthMap[ph.CommunityArea] = ph
	}

	var geographies []Geographies
	for _, ccvi := range ccviData {
		geo := Geographies{
			// Copy fields from ccvi
			GeographyType:                    ccvi.GeographyType,
			CommunityAreaOrZip:               ccvi.CommunityAreaOrZip.String(),
			CommunityAreaName:                ccvi.CommunityAreaName,
			CcviScore:                        ccvi.CcviScore,
			CcviCategory:                     ccvi.CcviCategory,
			RankSocioeconomicStatus:          ccvi.RankSocioeconomicStatus,
			RankAdultsNoPcp:                  ccvi.RankAdultsNoPcp,
			RankCumulativeMobilityRatio:      ccvi.RankCumulativeMobilityRatio,
			RankFrontlineEssentialWorkers:    ccvi.RankFrontlineEssentialWorkers,
			RankAge65Plus:                    ccvi.RankAge65Plus,
			RankComorbidConditions:           ccvi.RankComorbidConditions,
			RankCovid19IncidenceRate:         ccvi.RankCovid19IncidenceRate,
			RankCovid19HospitalAdmissionRate: ccvi.RankCovid19HospitalAdmissionRate,
			RankCovid19CrudeMortalityRate:    ccvi.RankCovid19CrudeMortalityRate,
		}

		// If there is corresponding PubHealth data, copy it
		if ph, ok := pubHealthMap[ccvi.CommunityAreaOrZip.String()]; ok {
			geo.BelowPovertyLevel = ph.BelowPovertyLevel
			geo.CrowdedHousing = ph.CrowdedHousing
			geo.NoHighSchoolDiploma = ph.NoHighSchoolDiploma
			geo.PerCapitaIncome = ph.PerCapitaIncome
			geo.Unemployment = ph.Unemployment
		}

		geographies = append(geographies, geo)
	}

	return geographies, nil
}

func LoadGeographies(db *gorm.DB) error {
	geographies, err := TransformGeographies()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Printf("Data transformed successfully")
	pkg.LoadToPostgres(db, geographies)
	fmt.Printf("Data loaded successfully")

	return nil
}
