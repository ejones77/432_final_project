package daily

import (
	"fmt"
	"time"

	"github.com/ejones77/432_final_project/pkg"
	"gorm.io/gorm"
)

type BuildingPermits struct {
	ID                   string            `json:"id" db:"id"`
	PermitNumber         string            `json:"permit_" db:"permit_number"`
	PermitStatus         string            `json:"permit_status" db:"permit_status"`
	PermitMilestone      string            `json:"permit_milestone" db:"permit_milestone"`
	PermitType           string            `json:"permit_type" db:"permit_type"`
	ReviewType           string            `json:"review_type" db:"review_type"`
	ApplicationStartDate pkg.CustomTime    `json:"application_start_date" db:"application_start_date"`
	IssueDate            pkg.CustomTime    `json:"issue_date" db:"issue_date"`
	WorkDescription      string            `json:"work_description" db:"work_description"`
	BuildingFeePaid      pkg.Float64String `json:"building_fee_paid" db:"building_fee_paid"`
	ZoningFeePaid        pkg.Float64String `json:"zoning_fee_paid" db:"zoning_fee_paid"`
	OtherFeePaid         pkg.Float64String `json:"other_fee_paid" db:"other_fee_paid"`
	BuildingFeeSubtotal  pkg.Float64String `json:"building_fee_subtotal" db:"building_fee_subtotal"`
	ZoningFeeSubtotal    pkg.Float64String `json:"zoning_fee_subtotal" db:"zoning_fee_subtotal"`
	OtherFeeSubtotal     pkg.Float64String `json:"other_fee_subtotal" db:"other_fee_subtotal"`
	BuildingFeeWaived    pkg.Float64String `json:"building_fee_waived" db:"building_fee_waived"`
	ZoningFeeWaived      pkg.Float64String `json:"zoning_fee_waived" db:"zoning_fee_waived"`
	OtherFeeWaived       pkg.Float64String `json:"other_fee_waived" db:"other_fee_waived"`
	SubtotalWaived       pkg.Float64String `json:"subtotal_waived" db:"subtotal_waived"`
	TotalFee             pkg.Float64String `json:"total_fee" db:"total_fee"`
	CommunityArea        string            `json:"community_area" db:"community_area"`
	Latitude             pkg.Float64String `json:"latitude" db:"latitude"`
	Longitude            pkg.Float64String `json:"longitude" db:"longitude"`
}

func ExtractBuildingPermits(db *gorm.DB) ([]BuildingPermits, error) {

	columns := []string{
		"id",
		"permit_",
		"permit_status",
		"permit_milestone",
		"permit_type",
		"review_type",
		"application_start_date",
		"issue_date",
		"work_description",
		"building_fee_paid",
		"zoning_fee_paid",
		"other_fee_paid",
		"building_fee_subtotal",
		"zoning_fee_subtotal",
		"other_fee_subtotal",
		"building_fee_waived",
		"zoning_fee_waived",
		"other_fee_waived",
		"subtotal_waived",
		"total_fee",
		"community_area",
		"latitude",
		"longitude",
	}

	var results []BuildingPermits

	// Check if the table exists
	var startDate, endDate string
	if !pkg.IsEmpty(db, "building_permits") {
		// Calculate start and end dates based on max date
		var maxDate time.Time
		db.Table("building_permits").Select("max(application_start_date)").Scan(&maxDate)
		startDate = maxDate.Format("2006-01-02")
		endDate = maxDate.AddDate(0, 0, 1).Format("2006-01-02")
	} else {
		// Use this first month as a sample
		startDate = "2020-04-01"
		endDate = "2020-05-01"
	}

	err := pkg.ConcurrentQuerySample("ydr8-5enu",
		"application_start_date",
		columns,
		fmt.Sprintf(`application_start_date >= '%s' AND application_start_date < '%s'`, startDate, endDate),
		4,
		2000,
		&results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func LoadBuildingPermits(db *gorm.DB) error {
	data, err := ExtractBuildingPermits(db)
	if err != nil {
		return err
	}

	pkg.LoadToPostgres(db, data)
	return nil
}

func UpdateBuildingPermits(db *gorm.DB, permits []BuildingPermits) error {
	for _, permit := range permits {
		result := db.Model(&BuildingPermits{}).Where("id = ?", permit.ID).Updates(permit)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
