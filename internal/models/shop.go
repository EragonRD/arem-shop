package models

import (
	"time"

	"github.com/google/uuid"
)

// Shop représente un tenant isolé du système.
type Shop struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name           string    `gorm:"size:120;not null" json:"name"`
	Active         bool      `gorm:"not null;default:true" json:"active"`
	WhatsAppNumber string    `gorm:"size:20;not null" json:"whatsAppNumber"`
	CreatedAt      time.Time `gorm:"not null;default:now()" json:"createdAt"`
}

func (Shop) TableName() string {
	return "shops"
}
