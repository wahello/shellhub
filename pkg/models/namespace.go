package models

import (
	"time"
)

type Namespace struct {
	Name         string             `json:"name"  validate:"required,hostname_rfc1123,excludes=."`
	Owner        string             `json:"owner"`
	APITokens    []Token            `json:"api_tokens" bson:"api_tokens"`
	TenantID     string             `json:"tenant_id" bson:"tenant_id,omitempty"`
	Members      []interface{}      `json:"members" bson:"members"`
	Settings     *NamespaceSettings `json:"settings"`
	Devices      int                `json:"devices" bson:",omitempty"`
	Sessions     int                `json:"sessions" bson:",omitempty"`
	MaxDevices   int                `json:"max_devices" bson:"max_devices"`
	DevicesCount int                `json:"devices_count" bson:"devices_count,omitempty"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	Billing      *Billing           `json:"billing" bson:"billing,omitempty"`
}

type NamespaceSettings struct {
	SessionRecord bool `json:"session_record" bson:"session_record,omitempty"`
}

type Member struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name,omitempty" bson:"-"`
}

type Token struct {
	ID       string `json:"id" bson:"id"`
	TenantID string `json:"tenant_id" bson:"tenant_id"`
	ReadOnly bool   `json:"read_only" bson:"read_only"`
}

type APITokenAuthClaims struct {
	ID       string `json:"id"`
	TenantID string `json:"tenant_id"`
	ReadOnly bool   `json:"read_only"`

	AuthClaims         `mapstruct:",squash"`
	jwt.StandardClaims `mapstruct:",squash"`
}

type APITokenAuthRequest struct {
	TenantID string `json:"tenant_id"`
}

type APITokenAuthResponse struct {
	ID       string `json:"id"`
	APIToken string `json:"api_token"`
	TenantID string `json:"tenant_id"`
	ReadOnly bool   `json:"read_only"`
}

type TokenFields struct {
	ReadOnly bool `json:"read_only"`
}

type APITokenUpdate struct {
	TokenFields `bson:",inline"`
}
