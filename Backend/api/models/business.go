package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Business struct {
	ID                   uuid.UUID      `gorm:"primaryKey;type:uuid;column:business_id"`
	Name                 string         `gorm:"column:business_name"`
	Email                string         `gorm:"column:business_email"`
	Phone                string         `gorm:"column:business_phone"` 
	Onboarded            bool           `gorm:"column:business_onboarded"`
	Fqdns                pq.StringArray `gorm:"index;type:text[];column:business_fqdns"`
	TrackingPageURL      string         `gorm:"column:business_tracking_page_url"`
	ReplaceTrackingLinks *bool          `gorm:"column:business_replace_tracking_links;default:FALSE"`
	CurrentBillingID     uuid.UUID      `gorm:"column:business_current_billing_id"`
	IsFreeTrialUsed      bool           `gorm:"column:business_is_free_trial_used"`
	AuthProvider         string         `gorm:"column:business_auth_provider"`
	CountryCode          string         `gorm:"column:business_country_code"`
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
