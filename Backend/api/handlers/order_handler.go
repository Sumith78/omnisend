package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

// ... (Business, Product, Order, Shipment, Checkpoint structs)

type ProductArray []Product

func (p *ProductArray) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte), &p)
}

func (p ProductArray) Value() (driver.Value, error) {
	if p == nil {
		return nil, nil
	}
	return json.Marshal(p)
}

type ShipmentStatus string
type TrackerServiceCode string

// ... (GetOrders, FetchOrders, UpdateStatus, NotifyShipmentStatusChange functions)

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

	// Find and update the shipment status
	for i, order := range orders {
		if order.ID.String() == request.OrderID {
			orders[i].ShipmentStatus = request.NewStatus
			orders[i].Email = request.Email
			NotifyShipmentStatusChange(orders[i].Email, orders[i].ID.String())
			break
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Shipment status updated successfully"})
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
