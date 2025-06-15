package db

import (
	"campaignservice/models"
	"campaignservice/utils"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

func LoadDBConfig() models.DBConfig {
	return models.DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
	}
}

func Connect() *sql.DB {
	cfg := LoadDBConfig()
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("failed to connect to DB:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("DB not reachable:", err)
	}
	return db
}

func GetTargetedCampaigns(db *sql.DB, appID, country, os string) ([]models.Campaign, error) {
	dimensionValues := map[string]string{
		"app":     appID,
		"country": country,
		"os":      os,
	}

	var conditions []string
	var args []interface{}
	argIndex := 1

	for _, dim := range utils.TargetingDimensions {
		val := dimensionValues[dim]
		conditions = append(conditions,
			fmt.Sprintf("(tr.dimension = '%s' AND tr.type = 'include' AND tr.value != $%d)", dim, argIndex),
			fmt.Sprintf("(tr.dimension = '%s' AND tr.type = 'exclude' AND tr.value = $%d)", dim, argIndex),
		)
		args = append(args, val)
		argIndex++
	}

	whereClause := strings.Join(conditions, " OR ")

	query := fmt.Sprintf(`
        SELECT c.campaign_id, c.campaign_name, c.image_url, c.call_to_action
        FROM campaigns c
        WHERE c.campaign_status = 'ACTIVE'
        AND NOT EXISTS (
            SELECT 1 FROM targeting_rules tr
            WHERE tr.campaign_id = c.campaign_id
            AND (%s)
        )`, whereClause)

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Printf("DB query failed: %v\n", err)
		return nil, err
	}
	defer rows.Close()

	var campaigns []models.Campaign
	for rows.Next() {
		var c models.Campaign
		err := rows.Scan(&c.CampaignID, &c.CampaignName, &c.ImageURL, &c.CallToAction)
		if err != nil {
			return nil, err
		}
		campaigns = append(campaigns, c)
	}
	return campaigns, nil
}
