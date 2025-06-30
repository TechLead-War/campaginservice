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

type DeliveryRequest struct {
	AppID   string `json:"app_id"`
	Country string `json:"country"`
	OS      string `json:"os"`
}

type DeliveryResponse struct {
	CampaignID   string `json:"campaign_id"`
	ImageURL     string `json:"image_url"`
	CallToAction string `json:"call_to_action"`
}
