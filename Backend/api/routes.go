package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yourusername/yourproject/api/handlers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/api/orders", handlers.GetOrdersHandler)
	r.POST("/api/update-status", handlers.UpdateShipmentStatusHandler)
	r.POST("/api/receive-notification", handlers.ReceiveNotificationHandler)

	return r
}
