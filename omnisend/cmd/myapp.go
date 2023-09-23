package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type Business struct {
	ID          string `json:"business_id"`
	Name        string `json:"business_name"`
	// Add other fields as needed
}

type Order struct {
	ID                  string    `json:"order_id"`
	BusinessID          string    `json:"business_id"`
	OrderTotal          float32   `json:"order_order_total"`
	Currency            string    `json:"order_currency"`
	OrderNumber         string    `json:"order_order_number"`
	CustomerFirstName   string    `json:"order_customer_first_name"`
	CustomerLastName    string    `json:"order_customer_last_name"`
	CustomerEmail       string    `json:"order_customer_email"`
	CustomerOrderNumber string    `json:"order_customer_order_number"`
	CustomerPhone       string    `json:"order_customer_phone"`
	CustomerAddress     string    `json:"order_customer_address"`
	ProviderOrderID     string    `json:"order_provider_order_id"`
	OrderCreatedAt      string    `json:"order_created_at"`
	ProviderCreatedAt   string    `json:"order_provider_created_at"`
	Products            []Product `json:"order_products"`
	UniqueID            string    `json:"order_unique_id"`
	ShipmentStatus      string    `json:"shipment_status"` // Added shipment status
	Email 				string	  `json:"email"`          // Added email
	// Add other fields as needed
}

type Product struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float32 `json:"price"`
	// Add other fields as needed
}

type Shipment struct {
	ID                      string  `json:"shipment_id"`
	OrderID                 string  `json:"shipment_order_id"`
	BusinessID              string  `json:"shipment_business_id"`
	Status                  string  `json:"shipment_status"`
	ProviderFulfillmentID   string  `json:"shipment_provider_fulfillment_id"`
	// Add other fields as needed
}

var orders []Order

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

    // Find and update the shipment status
    for i, order := range orders {
        if order.ID == request.OrderID {
            orders[i].ShipmentStatus = request.NewStatus
            orders[i].Email = request.Email
            NotifyShipmentStatusChange(orders[i].Email, orders[i].ID)
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

func GetOrdersHandler(c *gin.Context) {
    orders := GetOrders()
    c.JSON(http.StatusOK, orders)
}

func ReceiveNotificationHandler(c *gin.Context) {
    var notification struct {
        // Define the fields of the notification struct based on the actual data
    }

    if err := c.ShouldBindJSON(&notification); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Process the notification data
    // For example, update your database based on the notification content
    // Or trigger some other actions based on the received notification

    c.JSON(http.StatusOK, gin.H{"message": "Notification handled successfully"})
}


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
		"unique_id": "UNIQ12345"
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
		"unique_id": "UNIQ54321"
	}
]
`

func GetOrders() []Order {
	var orders []Order
	if err := json.Unmarshal([]byte(sampleOrders), &orders); err != nil {
		log.Println("Error unmarshalling sample orders:", err)
	}
	return orders
}