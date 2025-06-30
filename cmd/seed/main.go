package main

import (
	"campaign/internal/domain/models"
	"campaign/internal/infrastructure/db"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Campaign data structure
type CampaignData struct {
	ID           string
	Name         string
	ImageURL     string
	CallToAction string
	Status       string
}

// Targeting rule data structure
type TargetingRule struct {
	CampaignID string
	Dimension  string
	Type       string
	Value      string
}

func main() {
	// Load configuration
	cfg, err := models.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Connect to database
	dbConnString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHOST, cfg.DBPORT, cfg.DBUSER, cfg.DBPass, cfg.DBName)

	d, err := db.Connect(dbConnString)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer d.Close()

	log.Println("Connected to database successfully")

	// Seed campaigns
	if err := seedCampaigns(d); err != nil {
		log.Fatalf("Error seeding campaigns: %v", err)
	}

	// Seed targeting rules
	if err := seedTargetingRules(d); err != nil {
		log.Fatalf("Error seeding targeting rules: %v", err)
	}

	// Seed some inactive campaigns for testing
	if err := seedInactiveCampaigns(d); err != nil {
		log.Fatalf("Error seeding inactive campaigns: %v", err)
	}

	log.Println("Database seeded successfully!")
	log.Println("You can now test the API with various combinations:")
	log.Println("- app_id: test_app, country: US, os: android")
	log.Println("- app_id: test_app, country: CA, os: ios")
	log.Println("- app_id: premium_app, country: US, os: android")
}

func seedCampaigns(db *sql.DB) error {
	campaigns := []CampaignData{
		{"camp_001", "Summer Sale Campaign", "https://example.com/images/summer_sale.jpg", "Shop Now", "ACTIVE"},
		{"camp_002", "New App Launch", "https://example.com/images/app_launch.jpg", "Download Now", "ACTIVE"},
		{"camp_003", "Holiday Special", "https://example.com/images/holiday.jpg", "Get 50% Off", "ACTIVE"},
		{"camp_004", "Premium Upgrade", "https://example.com/images/premium.jpg", "Upgrade Now", "ACTIVE"},
		{"camp_005", "Referral Bonus", "https://example.com/images/referral.jpg", "Invite Friends", "ACTIVE"},
		{"camp_006", "Black Friday Deals", "https://example.com/images/black_friday.jpg", "Shop Deals", "ACTIVE"},
		{"camp_007", "Christmas Special", "https://example.com/images/christmas.jpg", "Holiday Offers", "ACTIVE"},
		{"camp_008", "New Year Sale", "https://example.com/images/new_year.jpg", "Start Fresh", "ACTIVE"},
		{"camp_009", "Valentine's Day", "https://example.com/images/valentine.jpg", "Love Deals", "ACTIVE"},
		{"camp_010", "Easter Special", "https://example.com/images/easter.jpg", "Spring Offers", "ACTIVE"},
		{"camp_011", "Back to School", "https://example.com/images/back_to_school.jpg", "Get Ready", "ACTIVE"},
		{"camp_012", "Halloween Spooky Deals", "https://example.com/images/halloween.jpg", "Trick or Treat", "ACTIVE"},
		{"camp_013", "Cyber Monday", "https://example.com/images/cyber_monday.jpg", "Tech Deals", "ACTIVE"},
		{"camp_014", "Spring Collection", "https://example.com/images/spring_collection.jpg", "New Arrivals", "ACTIVE"},
		{"camp_015", "Summer Vacation", "https://example.com/images/summer_vacation.jpg", "Book Now", "ACTIVE"},
		{"camp_016", "Fall Fashion", "https://example.com/images/fall_fashion.jpg", "Style Update", "ACTIVE"},
		{"camp_017", "Winter Warmth", "https://example.com/images/winter_warmth.jpg", "Stay Cozy", "ACTIVE"},
		{"camp_018", "Gaming Tournament", "https://example.com/images/gaming_tournament.jpg", "Join Battle", "ACTIVE"},
		{"camp_019", "Fitness Challenge", "https://example.com/images/fitness_challenge.jpg", "Get Fit", "ACTIVE"},
		{"camp_020", "Food Festival", "https://example.com/images/food_festival.jpg", "Taste Now", "ACTIVE"},
		{"camp_021", "Music Concert", "https://example.com/images/music_concert.jpg", "Get Tickets", "ACTIVE"},
		{"camp_022", "Travel Adventure", "https://example.com/images/travel_adventure.jpg", "Explore More", "ACTIVE"},
		{"camp_023", "Tech Innovation", "https://example.com/images/tech_innovation.jpg", "Discover Tech", "ACTIVE"},
		{"camp_024", "Art Exhibition", "https://example.com/images/art_exhibition.jpg", "View Gallery", "ACTIVE"},
		{"camp_025", "Sports Championship", "https://example.com/images/sports_championship.jpg", "Watch Live", "ACTIVE"},
	}

	for _, campaign := range campaigns {
		query := `
			INSERT INTO campaigns (campaign_id, campaign_name, image_url, call_to_action, campaign_status)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (campaign_id) DO UPDATE SET
				campaign_name = EXCLUDED.campaign_name,
				image_url = EXCLUDED.image_url,
				call_to_action = EXCLUDED.call_to_action,
				campaign_status = EXCLUDED.campaign_status,
				udate = NOW()
		`
		_, err := db.Exec(query, campaign.ID, campaign.Name, campaign.ImageURL, campaign.CallToAction, campaign.Status)
		if err != nil {
			return fmt.Errorf("error inserting campaign %s: %v", campaign.ID, err)
		}
		log.Printf("Inserted/updated campaign: %s - %s", campaign.ID, campaign.Name)
	}

	return nil
}

func seedTargetingRules(db *sql.DB) error {
	rules := []TargetingRule{
		// App targeting - test_app
		{"camp_001", "app_id", "include", "test_app"},
		{"camp_002", "app_id", "include", "test_app"},
		{"camp_003", "app_id", "include", "test_app"},
		{"camp_004", "app_id", "include", "test_app"},
		{"camp_005", "app_id", "include", "test_app"},
		{"camp_006", "app_id", "include", "test_app"},
		{"camp_007", "app_id", "include", "test_app"},
		{"camp_008", "app_id", "include", "test_app"},
		{"camp_009", "app_id", "include", "test_app"},
		{"camp_010", "app_id", "include", "test_app"},
		{"camp_011", "app_id", "include", "test_app"},
		{"camp_012", "app_id", "include", "test_app"},
		{"camp_013", "app_id", "include", "test_app"},
		{"camp_014", "app_id", "include", "test_app"},
		{"camp_015", "app_id", "include", "test_app"},
		{"camp_016", "app_id", "include", "test_app"},
		{"camp_017", "app_id", "include", "test_app"},
		{"camp_018", "app_id", "include", "test_app"},
		{"camp_019", "app_id", "include", "test_app"},
		{"camp_020", "app_id", "include", "test_app"},
		{"camp_021", "app_id", "include", "test_app"},
		{"camp_022", "app_id", "include", "test_app"},
		{"camp_023", "app_id", "include", "test_app"},
		{"camp_024", "app_id", "include", "test_app"},
		{"camp_025", "app_id", "include", "test_app"},

		// App targeting - premium_app
		{"camp_001", "app_id", "include", "premium_app"},
		{"camp_002", "app_id", "include", "premium_app"},
		{"camp_004", "app_id", "include", "premium_app"},
		{"camp_005", "app_id", "include", "premium_app"},
		{"camp_006", "app_id", "include", "premium_app"},
		{"camp_011", "app_id", "include", "premium_app"},
		{"camp_013", "app_id", "include", "premium_app"},
		{"camp_015", "app_id", "include", "premium_app"},
		{"camp_017", "app_id", "include", "premium_app"},
		{"camp_019", "app_id", "include", "premium_app"},
		{"camp_021", "app_id", "include", "premium_app"},
		{"camp_023", "app_id", "include", "premium_app"},
		{"camp_025", "app_id", "include", "premium_app"},

		// App targeting - gaming_app
		{"camp_018", "app_id", "include", "gaming_app"},
		{"camp_019", "app_id", "include", "gaming_app"},
		{"camp_025", "app_id", "include", "gaming_app"},

		// App targeting - travel_app
		{"camp_015", "app_id", "include", "travel_app"},
		{"camp_022", "app_id", "include", "travel_app"},

		// App targeting - food_app
		{"camp_020", "app_id", "include", "food_app"},

		// App targeting - music_app
		{"camp_021", "app_id", "include", "music_app"},

		// App targeting - art_app
		{"camp_024", "app_id", "include", "art_app"},

		// Country targeting - US
		{"camp_001", "country", "include", "US"},
		{"camp_002", "country", "include", "US"},
		{"camp_004", "country", "include", "US"},
		{"camp_005", "country", "include", "US"},
		{"camp_006", "country", "include", "US"},
		{"camp_007", "country", "include", "US"},
		{"camp_008", "country", "include", "US"},
		{"camp_009", "country", "include", "US"},
		{"camp_010", "country", "include", "US"},
		{"camp_011", "country", "include", "US"},
		{"camp_012", "country", "include", "US"},
		{"camp_013", "country", "include", "US"},
		{"camp_014", "country", "include", "US"},
		{"camp_015", "country", "include", "US"},
		{"camp_016", "country", "include", "US"},
		{"camp_017", "country", "include", "US"},
		{"camp_018", "country", "include", "US"},
		{"camp_019", "country", "include", "US"},
		{"camp_020", "country", "include", "US"},
		{"camp_021", "country", "include", "US"},
		{"camp_022", "country", "include", "US"},
		{"camp_023", "country", "include", "US"},
		{"camp_024", "country", "include", "US"},
		{"camp_025", "country", "include", "US"},

		// Country targeting - CA
		{"camp_001", "country", "include", "CA"},
		{"camp_003", "country", "include", "CA"},
		{"camp_004", "country", "include", "CA"},
		{"camp_005", "country", "include", "CA"},
		{"camp_007", "country", "include", "CA"},
		{"camp_008", "country", "include", "CA"},
		{"camp_011", "country", "include", "CA"},
		{"camp_012", "country", "include", "CA"},
		{"camp_015", "country", "include", "CA"},
		{"camp_017", "country", "include", "CA"},
		{"camp_020", "country", "include", "CA"},
		{"camp_022", "country", "include", "CA"},
		{"camp_024", "country", "include", "CA"},

		// Country targeting - UK
		{"camp_002", "country", "include", "UK"},
		{"camp_003", "country", "include", "UK"},
		{"camp_006", "country", "include", "UK"},
		{"camp_009", "country", "include", "UK"},
		{"camp_010", "country", "include", "UK"},
		{"camp_011", "country", "include", "UK"},
		{"camp_013", "country", "include", "UK"},
		{"camp_014", "country", "include", "UK"},
		{"camp_016", "country", "include", "UK"},
		{"camp_018", "country", "include", "UK"},
		{"camp_021", "country", "include", "UK"},
		{"camp_023", "country", "include", "UK"},
		{"camp_025", "country", "include", "UK"},

		// OS targeting - android
		{"camp_001", "os", "include", "android"},
		{"camp_003", "os", "include", "android"},
		{"camp_005", "os", "include", "android"},
		{"camp_006", "os", "include", "android"},
		{"camp_007", "os", "include", "android"},
		{"camp_008", "os", "include", "android"},
		{"camp_009", "os", "include", "android"},
		{"camp_010", "os", "include", "android"},
		{"camp_011", "os", "include", "android"},
		{"camp_012", "os", "include", "android"},
		{"camp_013", "os", "include", "android"},
		{"camp_014", "os", "include", "android"},
		{"camp_015", "os", "include", "android"},
		{"camp_016", "os", "include", "android"},
		{"camp_017", "os", "include", "android"},
		{"camp_018", "os", "include", "android"},
		{"camp_019", "os", "include", "android"},
		{"camp_020", "os", "include", "android"},
		{"camp_021", "os", "include", "android"},
		{"camp_022", "os", "include", "android"},
		{"camp_023", "os", "include", "android"},
		{"camp_024", "os", "include", "android"},
		{"camp_025", "os", "include", "android"},

		// OS targeting - ios
		{"camp_002", "os", "include", "ios"},
		{"camp_004", "os", "include", "ios"},
		{"camp_005", "os", "include", "ios"},
		{"camp_006", "os", "include", "ios"},
		{"camp_007", "os", "include", "ios"},
		{"camp_008", "os", "include", "ios"},
		{"camp_009", "os", "include", "ios"},
		{"camp_010", "os", "include", "ios"},
		{"camp_011", "os", "include", "ios"},
		{"camp_012", "os", "include", "ios"},
		{"camp_013", "os", "include", "ios"},
		{"camp_014", "os", "include", "ios"},
		{"camp_015", "os", "include", "ios"},
		{"camp_016", "os", "include", "ios"},
		{"camp_017", "os", "include", "ios"},
		{"camp_018", "os", "include", "ios"},
		{"camp_019", "os", "include", "ios"},
		{"camp_020", "os", "include", "ios"},
		{"camp_021", "os", "include", "ios"},
		{"camp_022", "os", "include", "ios"},
		{"camp_023", "os", "include", "ios"},
		{"camp_024", "os", "include", "ios"},
		{"camp_025", "os", "include", "ios"},

		// OS targeting - windows
		{"camp_001", "os", "include", "windows"},
		{"camp_002", "os", "include", "windows"},
		{"camp_004", "os", "include", "windows"},
		{"camp_006", "os", "include", "windows"},
		{"camp_008", "os", "include", "windows"},
		{"camp_011", "os", "include", "windows"},
		{"camp_013", "os", "include", "windows"},
		{"camp_015", "os", "include", "windows"},
		{"camp_017", "os", "include", "windows"},
		{"camp_018", "os", "include", "windows"},
		{"camp_019", "os", "include", "windows"},
		{"camp_021", "os", "include", "windows"},
		{"camp_023", "os", "include", "windows"},
		{"camp_025", "os", "include", "windows"},

		// OS targeting - macos
		{"camp_013", "os", "include", "macos"},
		{"camp_014", "os", "include", "macos"},
		{"camp_016", "os", "include", "macos"},
		{"camp_020", "os", "include", "macos"},
		{"camp_021", "os", "include", "macos"},
		{"camp_023", "os", "include", "macos"},
		{"camp_024", "os", "include", "macos"},

		// Exclude rules for testing
		{"camp_001", "country", "exclude", "JP"},
		{"camp_002", "os", "exclude", "linux"},
		{"camp_003", "app_id", "exclude", "old_app"},
		{"camp_004", "country", "exclude", "RU"},
		{"camp_005", "os", "exclude", "blackberry"},
		{"camp_011", "country", "exclude", "BR"},
		{"camp_012", "os", "exclude", "webos"},
		{"camp_013", "app_id", "exclude", "legacy_app"},
		{"camp_014", "country", "exclude", "CN"},
		{"camp_015", "os", "exclude", "symbian"},
		{"camp_016", "country", "exclude", "IN"},
		{"camp_017", "os", "exclude", "tizen"},
		{"camp_018", "country", "exclude", "KR"},
		{"camp_019", "os", "exclude", "firefoxos"},
		{"camp_020", "country", "exclude", "MX"},
		{"camp_021", "os", "exclude", "ubuntu"},
		{"camp_022", "country", "exclude", "AR"},
		{"camp_023", "os", "exclude", "chromeos"},
		{"camp_024", "country", "exclude", "ZA"},
		{"camp_025", "os", "exclude", "kaios"},
	}

	for _, rule := range rules {
		query := `
			INSERT INTO targeting_rules (campaign_id, dimension, type, value)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (campaign_id, dimension, type, value) DO UPDATE SET
				udate = NOW()
		`
		_, err := db.Exec(query, rule.CampaignID, rule.Dimension, rule.Type, rule.Value)
		if err != nil {
			return fmt.Errorf("error inserting targeting rule for campaign %s: %v", rule.CampaignID, err)
		}
		log.Printf("Inserted/updated targeting rule: %s %s %s %s", rule.CampaignID, rule.Dimension, rule.Type, rule.Value)
	}

	return nil
}

func seedInactiveCampaigns(db *sql.DB) error {
	inactiveCampaigns := []CampaignData{
		{"camp_inactive_001", "Expired Campaign", "https://example.com/images/expired.jpg", "Expired", "INACTIVE"},
		{"camp_inactive_002", "Paused Campaign", "https://example.com/images/paused.jpg", "Paused", "INACTIVE"},
		{"camp_inactive_003", "Draft Campaign", "https://example.com/images/draft.jpg", "Draft", "INACTIVE"},
	}

	for _, campaign := range inactiveCampaigns {
		query := `
			INSERT INTO campaigns (campaign_id, campaign_name, image_url, call_to_action, campaign_status)
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (campaign_id) DO UPDATE SET
				campaign_name = EXCLUDED.campaign_name,
				image_url = EXCLUDED.image_url,
				call_to_action = EXCLUDED.call_to_action,
				campaign_status = EXCLUDED.campaign_status,
				udate = NOW()
		`
		_, err := db.Exec(query, campaign.ID, campaign.Name, campaign.ImageURL, campaign.CallToAction, campaign.Status)
		if err != nil {
			return fmt.Errorf("error inserting inactive campaign %s: %v", campaign.ID, err)
		}
		log.Printf("Inserted/updated inactive campaign: %s - %s", campaign.ID, campaign.Name)
	}

	return nil
}
