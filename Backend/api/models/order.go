package models

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	ID                  uuid.UUID     `gorm:"primaryKey;type:uuid;column:order_id"`
	BusinessID          uuid.UUID     `gorm:"foreignKey:business_id;type:uuid;column:order_business_id"`
	OrderTotal          float32       `gorm:"column:order_order_total"`
	Currency            string        `gorm:"column:order_currency"`
	OrderNumber         string        `gorm:"column:order_order_number"`
	CustomerFirstName   string        `gorm:"column:order_customer_first_name"`
	CustomerLastName    string        `gorm:"column:order_customer_last_name"`
	CustomerEmail       string        `gorm:"column:order_customer_email"`
	CustomerOrderNumber string        `gorm:"column:order_customer_order_number"`
	CustomerPhone       string        `gorm:"column:order_customer_phone"`
	CustomerAddress     string        `gorm:"column:order_customer_address"`
	ProviderOrderID     string        `gorm:"column:order_provider_order_id"`
	CreatedAt           time.Time     `gorm:"column:order_created_at"`
	ProviderCreatedAt   time.Time     `gorm:"column:order_provider_created_at"`
	Products            ProductArray  `gorm:"column:order_products"`
	UniqueID            string        `gorm:"column:order_unique_id;uniqueIndex:idx_order_unique_id"`
	ShipmentStatus      ShipmentStatus `gorm:"column:order_shipment_status"`
	Email               string        `gorm:"column:order_email"`
}
