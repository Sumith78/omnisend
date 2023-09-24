package types
import (
	"time"

	
	"github.com/google/uuid"

	
)

type Checkpoint struct {
	ID             uuid.UUID      `gorm:"primaryKey;type:uuid;column:checkpoint_id"`
	ShipmentID     uuid.UUID      `gorm:"foreignKey:shipment_id;type:uuid;column:checkpoint_shipment_id"`
	Status         ShipmentStatus `gorm:"column:checkpoint_status"`
	Title          string         `gorm:"column:checkpoint_title"`
	Timestamp      time.Time      `gorm:"column:checkpoint_timestamp"`
	TrackerEventID string         `gorm:"column:checkpoint_tracker_event_id"`
	CreatedAt      time.Time      `gorm:"column:checkpoint_created_at"`
}
