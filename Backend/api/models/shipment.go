package models

import (
	"encoding/json"

	"time"
	"database/sql/driver"

	
	"github.com/google/uuid"

	
)
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

type ShipmentStatus string

type TrackerServiceCode string

var orders []Order