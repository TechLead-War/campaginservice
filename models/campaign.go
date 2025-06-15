package models

type Campaign struct {
	CampaignID     string
	CampaignName   string
	ImageURL       string
	CallToAction   string
	CampaignStatus string
	CDate          string
	UDate          string
}

type TargetingRule struct {
	CampaignID string
	Dimension  string
	Type       string
	Value      string
	CDate      string
	UDate      string
}
