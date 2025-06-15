package handler

import (
	"campaignservice/db"
	"campaignservice/models"
	"campaignservice/utils"
	"encoding/json"
	"net/http"
)

func DeliveryHandler(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	app := q.Get("app")
	osParam := q.Get("os")
	country := q.Get("country")

	switch {
	case app == "":
		utils.ErrorJSON(w, http.StatusBadRequest, utils.ErrMissingApp)
		return
	case osParam == "":
		utils.ErrorJSON(w, http.StatusBadRequest, utils.ErrMissingOS)
		return
	case country == "":
		utils.ErrorJSON(w, http.StatusBadRequest, utils.ErrMissingCountry)
		return
	}

	dbConn := db.Connect()
	defer dbConn.Close()

	campaigns, err := db.GetTargetedCampaigns(dbConn, app, osParam, country)
	if err != nil {
		utils.ErrorJSON(w, http.StatusInternalServerError, utils.InternalServerError)
		return
	}

	var response []models.DeliveryResponse
	for _, c := range campaigns {
		response = append(response, models.DeliveryResponse{
			CID: c.CampaignID,
			Img: c.ImageURL,
			CTA: c.CallToAction,
		})
	}

	if len(response) == 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
