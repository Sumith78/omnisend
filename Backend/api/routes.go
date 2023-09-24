package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/api/orders", GetOrdersHandler)
	r.POST("/api/update-status", UpdateShipmentStatusHandler)
	r.POST("/api/receive-notification", ReceiveNotificationHandler)

	return r
}

func GetOrdersHandler(c *gin.Context) {
	orders := GetOrders()
	c.JSON(http.StatusOK, orders)
}

func UpdateShipmentStatusHandler(c *gin.Context) {
	var request struct {
		OrderID   string `json:"order_id"`
		NewStatus string `json:"new_status"`
		Email     string `json:"email"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderID := request.OrderID
	newStatus := ShipmentStatus(request.NewStatus)

	if UpdateStatus(orderID, newStatus) {
		NotifyShipmentStatusChange(request.Email, orderID)
		c.JSON(http.StatusOK, gin.H{"message": "Shipment status updated successfully"})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
	}
}

func ReceiveNotificationHandler(c *gin.Context) {
	var notification struct {
		// Define the fields of the notification struct based on the actual data
	}

	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification handled successfully"})
}
