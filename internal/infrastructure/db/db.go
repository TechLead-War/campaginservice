package db

import (
	"campaign/internal/domain/models"
	"campaign/pkg/utils"
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/prometheus/client_golang/prometheus"
)

func Connect(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		log.Printf("failed to open DB connection: %v", err)
		return nil, err
	}
	if err = db.Ping(); err != nil {
		log.Printf("DB not reachable: %v", err)
		return nil, err
	}

	log.Println("DB connection established successfully")
	return db, nil
}

// TargetingDimension represents a targeting dimension with its value
type TargetingDimension struct {
	Dimension string
	Value     string
}

// GetTargetedCampaignsDynamic is a scalable version that supports any number of targeting dimensions
func GetTargetedCampaignsDynamic(db *sql.DB, dimensions []TargetingDimension, limit int, offset int) ([]models.Campaign, error) {
	timer := prometheus.NewTimer(utils.DBOperationDuration.WithLabelValues("GetTargetedCampaignsDynamic"))
	defer timer.ObserveDuration()

	if limit <= 0 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	if len(dimensions) == 0 {
		// If no dimensions provided, return all active campaigns
		return getAllActiveCampaigns(db, limit, offset)
	}

	// Build dynamic query
	query, args := buildDynamicTargetingQuery(dimensions, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("DB query failed for GetTargetedCampaignsDynamic: %v\nQuery: %s\n Args:%v", err, query, args)
		return nil, err
	}
	defer rows.Close()

	var campaigns []models.Campaign
	for rows.Next() {
		var campaign models.Campaign
		err := rows.Scan(&campaign.CampaignID, &campaign.CampaignName, &campaign.ImageURL, &campaign.CallToAction)
		if err != nil {
			log.Printf("Error scanning campaign row: %v", err)
			return nil, err
		}
		campaigns = append(campaigns, campaign)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error after iterating campaign row: %v", err)
		return nil, err
	}

	// Log dimensions for debugging
	dimStrs := make([]string, len(dimensions))
	for i, dim := range dimensions {
		dimStrs[i] = fmt.Sprintf("%s=%s", dim.Dimension, dim.Value)
	}
	log.Printf("GetTargetedCampaign found %d campaigns for dimensions: %s", len(campaigns), strings.Join(dimStrs, ", "))

	return campaigns, nil
}

// GetTargetedCampaigns is the original function for backward compatibility
func GetTargetedCampaigns(db *sql.DB, appID, country, os string, limit int, offset int) ([]models.Campaign, error) {
	dimensions := []TargetingDimension{
		{Dimension: "app_id", Value: appID},
		{Dimension: "country", Value: country},
		{Dimension: "os", Value: os},
	}
	return GetTargetedCampaignsDynamic(db, dimensions, limit, offset)
}

// buildDynamicTargetingQuery builds a dynamic SQL query based on provided dimensions
func buildDynamicTargetingQuery(dimensions []TargetingDimension, limit, offset int) (string, []interface{}) {
	var cteParts []string
	var joinParts []string
	var whereParts []string
	var args []interface{}
	argIndex := 1

	// Build CTEs for each dimension
	for _, dim := range dimensions {
		cteName := fmt.Sprintf("%s_targeting", dim.Dimension)
		cteParts = append(cteParts, fmt.Sprintf(`
		%s AS (
			SELECT campaign_id, 
				   bool_or(type = 'include' AND value = $%d) as has_include,
				   bool_or(type = 'exclude' AND value = $%d) as has_exclude,
				   count(*) FILTER (WHERE type = 'include') as include_count
			FROM targeting_rules 
			WHERE dimension = '%s' 
			GROUP BY campaign_id
		)`, cteName, argIndex, argIndex, dim.Dimension))

		joinParts = append(joinParts, fmt.Sprintf("LEFT JOIN %s %s ON c.campaign_id = %s.campaign_id",
			cteName, cteName, cteName))

		whereParts = append(whereParts, fmt.Sprintf("(%s.campaign_id IS NULL OR (%s.include_count = 0 OR %s.has_include) AND NOT %s.has_exclude)",
			cteName, cteName, cteName, cteName))

		args = append(args, dim.Value)
		argIndex++
	}

	// Add limit and offset
	args = append(args, limit, offset)

	// Build the complete query
	query := fmt.Sprintf(`
		WITH %s
		SELECT DISTINCT c.campaign_id, c.campaign_name, c.image_url, c.call_to_action
		FROM campaigns c
		%s
		WHERE c.campaign_status = 'ACTIVE'
		  AND %s
		ORDER BY c.campaign_id
		LIMIT $%d OFFSET $%d;
	`, strings.Join(cteParts, ","), strings.Join(joinParts, "\n\t\t"), strings.Join(whereParts, "\n\t\t  AND "), argIndex, argIndex+1)

	return query, args
}

// getAllActiveCampaigns returns all active campaigns when no targeting dimensions are provided
func getAllActiveCampaigns(db *sql.DB, limit, offset int) ([]models.Campaign, error) {
	query := `
		SELECT campaign_id, campaign_name, image_url, call_to_action
		FROM campaigns
		WHERE campaign_status = 'ACTIVE'
		ORDER BY campaign_id
		LIMIT $1 OFFSET $2;
	`

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		log.Printf("DB query failed for getAllActiveCampaigns: %v", err)
		return nil, err
	}
	defer rows.Close()

	var campaigns []models.Campaign
	for rows.Next() {
		var campaign models.Campaign
		err := rows.Scan(&campaign.CampaignID, &campaign.CampaignName, &campaign.ImageURL, &campaign.CallToAction)
		if err != nil {
			log.Printf("Error scanning campaign row: %v", err)
			return nil, err
		}
		campaigns = append(campaigns, campaign)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error after iterating campaign row: %v", err)
		return nil, err
	}

	log.Printf("GetAllActiveCampaigns found %d campaigns", len(campaigns))
	return campaigns, nil
}

// GetAvailableDimensions returns all available targeting dimensions from the database
func GetAvailableDimensions(db *sql.DB) ([]string, error) {
	query := `
		SELECT DISTINCT dimension 
		FROM targeting_rules 
		ORDER BY dimension;
	`

	rows, err := db.Query(query)
	if err != nil {
		log.Printf("DB query failed for GetAvailableDimensions: %v", err)
		return nil, err
	}
	defer rows.Close()

	var dimensions []string
	for rows.Next() {
		var dimension string
		err := rows.Scan(&dimension)
		if err != nil {
			log.Printf("Error scanning dimension row: %v", err)
			return nil, err
		}
		dimensions = append(dimensions, dimension)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error after iterating dimension row: %v", err)
		return nil, err
	}

	return dimensions, nil
}

// GetAvailableValuesForDimension returns all available values for a specific dimension
func GetAvailableValuesForDimension(db *sql.DB, dimension string) ([]string, error) {
	query := `
		SELECT DISTINCT value 
		FROM targeting_rules 
		WHERE dimension = $1 
		ORDER BY value;
	`

	rows, err := db.Query(query, dimension)
	if err != nil {
		log.Printf("DB query failed for GetAvailableValuesForDimension: %v", err)
		return nil, err
	}
	defer rows.Close()

	var values []string
	for rows.Next() {
		var value string
		err := rows.Scan(&value)
		if err != nil {
			log.Printf("Error scanning value row: %v", err)
			return nil, err
		}
		values = append(values, value)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Error after iterating value row: %v", err)
		return nil, err
	}

	return values, nil
}
