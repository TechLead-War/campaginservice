package main

import (
	"campaign/internal/api/handler"
	"campaign/internal/infrastructure/cache"
	"database/sql"

	"github.com/gin-gonic/gin"
)

func Delivery(router *gin.RouterGroup, db *sql.DB, memCache *cache.MemoryCache) {
	deliveryHandler := handler.NewDeliveryHandler(db, memCache)

	// Main delivery endpoint
	router.GET("/delivery", deliveryHandler.DeliveryHandler)

	// Discovery endpoints for available targeting options
	router.GET("/dimensions", deliveryHandler.GetAvailableDimensions)
	router.GET("/dimensions/:dimension/values", deliveryHandler.GetAvailableValues)
}
