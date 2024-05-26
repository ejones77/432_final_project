CREATE TABLE geographies (
    geography_type VARCHAR(255),
    community_area_or_zip VARCHAR(255),
    community_area_name VARCHAR(255),
    ccvi_score NUMERIC,
    ccvi_category VARCHAR(255),
    rank_socioeconomic_status NUMERIC,
    rank_adults_no_pcp NUMERIC,
    rank_cumulative_mobility_ratio NUMERIC,
    rank_frontline_essential_workers NUMERIC,
    rank_age_65_plus NUMERIC,
    rank_comorbid_conditions NUMERIC,
    rank_covid_19_incidence_rate NUMERIC,
    rank_covid_19_hospital_admission_rate NUMERIC,
    rank_covid_19_crude_mortality_rate NUMERIC,
    below_poverty_level NUMERIC,
    crowded_housing NUMERIC,
    no_high_school_diploma NUMERIC,
    per_capita_income NUMERIC,
    unemployment NUMERIC
);

CREATE TABLE taxi_rideshares (
    trip_id VARCHAR(255) PRIMARY KEY,
    taxi_id VARCHAR(255),
    trip_start_timestamp TIMESTAMP,
    trip_end_timestamp TIMESTAMP,
    trip_seconds NUMERIC,
    trip_miles NUMERIC,
    fare NUMERIC,
    tip NUMERIC,
    additional_charges NUMERIC,
    trip_total NUMERIC,
    pickup_centroid_latitude NUMERIC,
    pickup_centroid_longitude NUMERIC,
    dropoff_centroid_latitude NUMERIC,
    dropoff_centroid_longitude NUMERIC
);

CREATE TABLE covid_cases (
    zip_code VARCHAR(255),
    week_start DATE,
    week_end DATE,
    cases_weekly NUMERIC,
    case_rate_weekly NUMERIC,
    tests_weekly NUMERIC,
    percent_tested_positive_weekly NUMERIC,
    deaths_weekly NUMERIC,
    death_rate_weekly NUMERIC,
    population NUMERIC
);

CREATE TABLE building_permits (
    id VARCHAR(255) PRIMARY KEY,
    permit_number VARCHAR(255),
    permit_status VARCHAR(255),
    permit_milestone VARCHAR(255),
    permit_type VARCHAR(255),
    review_type VARCHAR(255),
    application_start_date DATE,
    issue_date DATE,
    work_description TEXT,
    building_fee_paid NUMERIC,
    zoning_fee_paid NUMERIC,
    other_fee_paid NUMERIC,
    building_fee_subtotal NUMERIC,
    zoning_fee_subtotal NUMERIC,
    other_fee_subtotal NUMERIC,
    building_fee_waived NUMERIC,
    zoning_fee_waived NUMERIC,
    other_fee_waived NUMERIC,
    subtotal_waived NUMERIC,
    total_fee NUMERIC,
    community_area VARCHAR(255),
    latitude NUMERIC,
    longitude NUMERIC
);

CREATE TABLE traffic_estimates (
    time TIMESTAMP,
    segment_id VARCHAR(255),
    speed NUMERIC,
    direction VARCHAR(255),
    length NUMERIC,
    bus_count NUMERIC,
    hour NUMERIC,
    day_of_week VARCHAR(255),
    month VARCHAR(255),
    start_latitude NUMERIC,
    start_longitude NUMERIC,
    end_latitude NUMERIC,
    end_longitude NUMERIC,
    start_zip_code VARCHAR(255),
    end_zip_code VARCHAR(255),
    start_community_area VARCHAR(255),
    end_community_area VARCHAR(255)
);