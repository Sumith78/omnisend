package types


import (
	

	"time"

	"database/sql/driver"
	"encoding/json"

	
	"github.com/google/uuid"
	
)

type Product struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float32 `json:"price"`
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