package handler

import (
	"campaign/internal/domain/models"
	"campaign/internal/infrastructure/cache"
	"campaign/internal/infrastructure/db"
	"campaign/pkg/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	cacheTTLE = 5 * time.Minute
)

type DeliveryHandler struct {
	db        *sql.DB
	memeCache *cache.MemoryCache
	//redis *redis.Client
}

func NewDeliveryHandler(db *sql.DB, memCache *cache.MemoryCache) *DeliveryHandler {
	return &DeliveryHandler{
		db:        db,
		memeCache: memCache,
		//redis: redis,
	}
}

func (h *DeliveryHandler) DeliveryHandler(c *gin.Context) {
	// Start timing for latency measurement
	start := time.Now()
	defer func() {
		// Record delivery API latency
		duration := time.Since(start)
		utils.DeliveryAPILatency.Observe(duration.Seconds())
	}()

	// Extract all query parameters for dynamic targeting
	targetingParams := h.extractTargetingParams(c)

	// Validate required parameters (keeping backward compatibility)
	if err := h.validateRequiredParams(targetingParams); err != nil {
		utils.ErrorJSONGin(c, http.StatusBadRequest, err.Error())
		return
	}

	// Parse and validate pagination parameters
	page, limit, err := h.parsePaginationParams(c)
	if err != nil {
		utils.ErrorJSONGin(c, http.StatusBadRequest, err.Error())
		return
	}

	offset := (page - 1) * limit

	// Generate cache key
	cacheKey := h.generateCacheKey(targetingParams, page, limit)

	// Try to get from cache first
	if cachedData, found := h.memeCache.Get(cacheKey); found {
		log.Printf("In-memory Cache HIT for key: %s", cacheKey)
		c.Header("X-Cache-Type", "IN_MEMORY_HIT")
		c.Data(http.StatusOK, "application/json", cachedData)
		utils.RecordCacheHit()
		return
	}

	log.Printf("In-memory Cache MISS for key: %s", cacheKey)
	utils.RecordCacheMiss()
	c.Header("X-Cache-Type", "MISS")

	// Convert targeting params to dimensions
	dimensions := h.convertToTargetingDimensions(targetingParams)

	// Get data from database using dynamic approach
	dbCampaigns, err := db.GetTargetedCampaignsDynamic(h.db, dimensions, limit, offset)
	if err != nil {
		log.Printf("Database error: %v", err)
		utils.ErrorJSONGin(c, http.StatusInternalServerError, utils.InternalServerError)
		return
	}

	// Build response
	response := h.buildResponse(dbCampaigns)

	// Marshal response for caching and sending
	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error marshalling delivery response: %v", err)
		utils.ErrorJSONGin(c, http.StatusInternalServerError, utils.InternalServerError)
		return
	}

	// Cache the response
	h.memeCache.Set(cacheKey, responseBytes, cacheTTLE)
	log.Printf("Successfully stored response in in-memory cache for key: %s", cacheKey)

	c.JSON(http.StatusOK, response)
}

// extractTargetingParams extracts all targeting parameters from the request
func (h *DeliveryHandler) extractTargetingParams(c *gin.Context) map[string]string {
	params := make(map[string]string)

	// Extract all query parameters that could be targeting dimensions
	for key, values := range c.Request.URL.Query() {
		// Skip pagination parameters
		if key == "page" || key == "limit" {
			continue
		}

		if len(values) > 0 && strings.TrimSpace(values[0]) != "" {
			params[key] = strings.TrimSpace(values[0])
		}
	}

	return params
}

// validateRequiredParams validates the required parameters (backward compatibility)
func (h *DeliveryHandler) validateRequiredParams(params map[string]string) error {
	requiredParams := []string{"app_id", "country", "os"}

	for _, param := range requiredParams {
		if value, exists := params[param]; !exists || value == "" {
			switch param {
			case "app_id":
				return fmt.Errorf(utils.ErrMissingApp)
			case "country":
				return fmt.Errorf(utils.ErrMissingCountry)
			case "os":
				return fmt.Errorf(utils.ErrMissingOS)
			}
		}
	}

	return nil
}

// convertToTargetingDimensions converts query parameters to targeting dimensions
func (h *DeliveryHandler) convertToTargetingDimensions(params map[string]string) []db.TargetingDimension {
	var dimensions []db.TargetingDimension

	// Define known targeting dimensions (can be extended)
	knownDimensions := []string{"app_id", "country", "os", "device_type", "language", "timezone", "age_group", "gender"}

	for _, dim := range knownDimensions {
		if value, exists := params[dim]; exists && value != "" {
			dimensions = append(dimensions, db.TargetingDimension{
				Dimension: dim,
				Value:     value,
			})
		}
	}

	return dimensions
}

func (h *DeliveryHandler) parsePaginationParams(c *gin.Context) (page, limit int, err error) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", strconv.Itoa(utils.DefaultApiPageLimit))

	page, err = strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return 1, utils.DefaultApiPageLimit, fmt.Errorf("invalid page parameter: %s", pageStr)
	}

	limit, err = strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		return page, utils.DefaultApiPageLimit, fmt.Errorf("invalid limit parameter: %s", limitStr)
	}

	return page, limit, nil
}

func (h *DeliveryHandler) generateCacheKey(params map[string]string, page, limit int) string {
	// Sort parameters for consistent cache keys
	var keys []string
	for k := range params {
		keys = append(keys, k)
	}

	// Actually sort the keys for consistent cache keys
	sort.Strings(keys)

	// Build sorted parameter string
	var paramParts []string
	for _, key := range keys {
		paramParts = append(paramParts, fmt.Sprintf("%s:%s", key, params[key]))
	}

	paramString := strings.Join(paramParts, ":")
	return fmt.Sprintf("delivery:%s:page%d:limit%d", paramString, page, limit)
}

func (h *DeliveryHandler) buildResponse(campaigns []models.Campaign) []models.DeliveryResponse {
	response := make([]models.DeliveryResponse, 0, len(campaigns))

	for _, campaign := range campaigns {
		response = append(response, models.DeliveryResponse{
			CampaignID:   campaign.CampaignID,
			ImageURL:     campaign.ImageURL,
			CallToAction: campaign.CallToAction,
		})
	}

	return response
}

// GetAvailableDimensions returns all available targeting dimensions
func (h *DeliveryHandler) GetAvailableDimensions(c *gin.Context) {
	dimensions, err := db.GetAvailableDimensions(h.db)
	if err != nil {
		log.Printf("Error getting available dimensions: %v", err)
		utils.ErrorJSONGin(c, http.StatusInternalServerError, utils.InternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"dimensions": dimensions,
	})
}

// GetAvailableValues returns all available values for a specific dimension
func (h *DeliveryHandler) GetAvailableValues(c *gin.Context) {
	dimension := c.Param("dimension")
	if dimension == "" {
		utils.ErrorJSONGin(c, http.StatusBadRequest, "dimension parameter is required")
		return
	}

	values, err := db.GetAvailableValuesForDimension(h.db, dimension)
	if err != nil {
		log.Printf("Error getting available values for dimension %s: %v", dimension, err)
		utils.ErrorJSONGin(c, http.StatusInternalServerError, utils.InternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"dimension": dimension,
		"values":    values,
	})
}
