package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/google/uuid"
)

type Order struct {
	ID                  uuid.UUID `gorm:"primaryKey;type:uuid;column:order_id"`
	BusinessID          uuid.UUID `gorm:"foreignKey:business_id;type:uuid;column:order_business_id"`
	OrderTotal          float32   `gorm:"column:order_order_total"`
	Currency            string    `gorm:"column:order_currency"`
	OrderNumber         string    `gorm:"column:order_order_number"`
	CustomerFirstName   string    `gorm:"column:order_customer_first_name"`
	CustomerLastName    string    `gorm:"column:order_customer_last_name"`
	CustomerEmail       string    `gorm:"column:order_customer_email"`
	CustomerOrderNumber string    `gorm:"column:order_customer_order_number"`
	CustomerPhone       string    `gorm:"column:order_customer_phone"`
	CustomerAddress     string    `gorm:"column:order_customer_address"`
	ProviderOrderID     string    `gorm:"column:order_provider_order_id"`
	CreatedAt           time.Time `gorm:"column:order_created_at"`
	ProviderCreatedAt   time.Time `gorm:"column:order_provider_created_at"`
	Products            ProductArray `gorm:"column:order_products"`
	UniqueID            string    `gorm:"column:order_unique_id;uniqueIndex:idx_order_unique_id"`
	ShipmentStatus      ShipmentStatus `gorm:"column:order_shipment_status"`
	Email               string    `gorm:"column:order_email"`
}

var orders []Order

var sampleOrders = `
[
	{
		"order_id": "1",
		"business_id": "101",
		"order_total": 50.0,
		"currency": "USD",
		"order_number": "ORD12345",
		"customer_first_name": "John",
		"customer_last_name": "Doe",
		"customer_email": "john.doe@example.com",
		"customer_order_number": "CUST54321",
		"customer_phone": "+1234567890",
		"customer_address": "123 Main St, Cityville, USA",
		"provider_order_id": "PROV67890",
		"order_created_at": "2023-09-21T12:34:56Z",
		"provider_created_at": "2023-09-21T13:45:56Z",
		"order_products": [
			{
				"product_id": "P101",
				"product_name": "Product A",
				"quantity": 2,
				"price": 25.0
			}
		],
		"unique_id": "UNIQ12345",
		"shipment_status": "Pending",
		"email": "john.doe@example.com"
	},
	{
		"order_id": "2",
		"business_id": "102",
		"order_total": 70.0,
		"currency": "EUR",
		"order_number": "ORD54321",
		"customer_first_name": "Jane",
		"customer_last_name": "Smith",
		"customer_email": "jane.smith@example.com",
		"customer_order_number": "CUST67890",
		"customer_phone": "+9876543210",
		"customer_address": "456 Oak St, Townsville, UK",
		"provider_order_id": "PROV09876",
		"order_created_at": "2023-09-20T10:11:12Z",
		"provider_created_at": "2023-09-20T11:12:13Z",
		"order_products": [
			{
				"product_id": "P102",
				"product_name": "Product B",
				"quantity": 1,
				"price": 70.0
			}
		],
		"unique_id": "UNIQ54321",
		"shipment_status": "Shipped",
		"email": "jane.smith@example.com"
	}
]
`

type ShipmentStatusUpdate struct {
	NewStatus string `json:"new_status"`
}

func GetOrders() []Order {
	var orders []Order
	if err := json.Unmarshal([]byte(sampleOrders), &orders); err != nil {
		log.Println("Error unmarshalling sample orders:", err)
	}
	return orders
}

func FetchOrders() []Order {
	return orders
}

func UpdateStatus(orderID string, status ShipmentStatusUpdate) bool {
	for i, order := range orders {
		if order.ID.String() == orderID {
			orders[i].ShipmentStatus = status.NewStatus
			return true
		}
	}
	return false
}

func main() {
	r := gin.Default()

	r.GET("/api/orders", GetOrdersHandler)
	r.POST("/api/update-status", UpdateShipmentStatusHandler)
	r.POST("/api/receive-notification", ReceiveNotificationHandler)

	r.Run(":8080")
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

func NotifyShipmentStatusChange(email string, orderID string) {
	from := mail.NewEmail("Your Name", "your@example.com")
	subject := "Shipment Status Update"
	to := mail.NewEmail("", email)
	plainTextContent := fmt.Sprintf("Your order with ID %s has a new shipment status.", orderID)
	htmlContent := fmt.Sprintf("<strong>Your order with ID %s has a new shipment status.</strong>", orderID)
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))

	response, err := client.Send(message)
	if err != nil {
		log.Println("Error sending email:", err)
		return
	}

	fmt.Println("Email sent successfully. Response:", response.StatusCode)
}
