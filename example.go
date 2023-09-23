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

type Business struct {
	ID                   uuid.UUID      `gorm:"primaryKey;type:uuid;column:business_id"`
	Name                 string         `gorm:"column:business_name"`
	Email                string         `gorm:"column:business_email"`
	Phone                string         `gorm:"column:business_phone"` // E.164 format
	Onboarded            bool           `gorm:"column:business_onboarded"`
	Fqdns                pq.StringArray `gorm:"index;type:text[];column:business_fqdns"` // List of FQDNs where the seller will load the order tracking widget
	TrackingPageURL      string         `gorm:"column:business_tracking_page_url"`       // URL of page that has the order tracking widget.
	ReplaceTrackingLinks *bool          `gorm:"column:business_replace_tracking_links;default:FALSE"`
	CurrentBillingID     uuid.UUID      `gorm:"column:business_current_billing_id"`
	IsFreeTrialUsed      bool           `gorm:"column:business_is_free_trial_used"`
	AuthProvider         string         `gorm:"column:business_auth_provider"`
	CountryCode          string         `gorm:"column:business_country_code"` // ISO-3166-1 alpha-2 format
	Currency             string         `gorm:"column:business_currency"`
	Category             string         `gorm:"column:business_category"`
	Verified             bool           `gorm:"column:business_verified"`
	EmailOptIn           *bool          `gorm:"type:boolean;default:FALSE;column:tracking_email_opt_in"`
	SMSOptIn             *bool          `gorm:"type:boolean;default:FALSE;column:tracking_sms_opt_in"`
	ShopifyAppScopes     string         `gorm:"column:business_shopify_app_scopes"`
	ShopifyAppAuthToken  string         `gorm:"column:business_shopify_app_auth_token"`
	ShopifyShopDomain    string         `gorm:"column:business_shopify_shop_domain"`
	DropshippingMode     *bool          `gorm:"column:business_dropshipping_mode;default:FALSE"`
	CreatedAt            time.Time      `gorm:"column:business_created_at"`
}

type Product struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float32 `json:"price"`
}

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

type Shipment struct {
	ID                      uuid.UUID          `gorm:"primaryKey;type:uuid;column:shipment_id"`
	OrderID                 uuid.UUID          `gorm:"foreignKey:order_id;type:uuid;column:shipment_order_id"`
	BusinessID              uuid.UUID          `gorm:"foreignKey:business_id;type:uuid;column:shipment_business_id"`
	Status                  ShipmentStatus     `gorm:"column:shipment_status"`
	ProviderFulfillmentID   string             `gorm:"column:shipment_provider_fulfillment_id"`
	CarrierName             string             `gorm:"column:shipment_carrier_name"`
	CarrierCode             string             `gorm:"column:shipment_carrier_code"`
	CarrierTrackingNumber   string             `gorm:"column:shipment_carrier_tracking_number"`
	EstimatedDeliveryDate   time.Time          `gorm:"column:shipment_estimated_delivery_date"`
	ActualDeliveryDate      time.Time          `gorm:"column:shipment_actual_delivery_date"`
	CreatedAt               time.Time          `gorm:"column:shipment_created_at"`
	ProviderCreatedAt       time.Time          `gorm:"column:shipment_provider_created_at"`
	Summary                 string             `gorm:"column:shipment_summary"`
	DestinationZipCode      string             `gorm:"column:shipment_destination_zip_code"`
	DestinationCity         string             `gorm:"column:shipment_destination_city"`
	DestinationState        string             `gorm:"column:shipment_destination_state"`
	DestinationCountryCode  string             `gorm:"column:shipment_destination_country_code"`
	DestinationLatitude     float64            `gorm:"column:shipment_destination_latitude"`
	DestinationLongitude    float64            `gorm:"column:shipment_destination_longitude"`
	OriginZipCode           string             `gorm:"column:shipment_origin_zip_code"`
	OriginCity              string             `gorm:"column:shipment_origin_city"`
	OriginState             string             `gorm:"column:shipment_origin_state"`
	OriginCountryCode       string             `gorm:"column:shipment_origin_country_code"`
	ReturnShipmentID        uuid.UUID          `gorm:"foreignKey:shipment_id;type:uuid;column:shipment_return_shipment_id"`
	Weight                  string             `gorm:"column:shipment_weight"`
	TrackerID               string             `gorm:"column:shipment_tracker_id"`
	ProviderTrackingPageURL string             `gorm:"column:shipment_provider_tracking_url"`
	TrackerServiceCode      TrackerServiceCode `gorm:"shipment_tracker_service_code"`
	UniqueID                string             `gorm:"column:shipment_unique_id;uniqueIndex:idx_shipment_unique_id"`
}

type Checkpoint struct {
	ID             uuid.UUID      `gorm:"primaryKey;type:uuid;column:checkpoint_id"`
	ShipmentID     uuid.UUID      `gorm:"foreignKey:shipment_id;type:uuid;column:checkpoint_shipment_id"`
	Status         ShipmentStatus `gorm:"column:checkpoint_status"`
	Title          string         `gorm:"column:checkpoint_title"`
	Timestamp      time.Time      `gorm:"column:checkpoint_timestamp"`
	TrackerEventID string         `gorm:"column:checkpoint_tracker_event_id"`
	CreatedAt      time.Time      `gorm:"column:checkpoint_created_at"`
}

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

var orders []Order

// Dummy data for testing
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
	// Update status in data source
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

	c.JSON(http.StatusOK, gin.H{"message": "Notification handled successfully"})
}
