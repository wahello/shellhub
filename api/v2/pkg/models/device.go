package models

import (
	"time"
)

type Device struct {
	UID       string          `json:"uid"`
	Name      string          `json:"name" bson:"name,omitempty" validate:"required,hostname_rfc1123"`
	Identity  *DeviceIdentity `json:"identity"`
	Info      *DeviceInfo     `json:"info"`
	PublicKey string          `json:"public_key" bson:"public_key"`
	TenantID  string          `json:"tenant_id" bson:"tenant_id"`
	LastSeen  time.Time       `json:"last_seen" bson:"last_seen"`
	Online    bool            `json:"online" bson:",omitempty"`
	Namespace string          `json:"namespace" bson:",omitempty"`
	Status    string          `json:"status" bson:"status,omitempty" validate:"oneof=accepted rejected pending unused`
}

type DeviceIdentity struct {
	MAC string `json:"mac"`
}

type DeviceInfo struct {
	ID         string `json:"id"`
	PrettyName string `json:"pretty_name"`
	Version    string `json:"version"`
	Arch       string `json:"arch"`
	Platform   string `json:"platform"`
}
