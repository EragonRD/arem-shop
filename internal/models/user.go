package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleSuperAdmin UserRole = "SuperAdmin"
	RoleAdmin      UserRole = "Admin"
)

// User représente un utilisateur interne rattaché à un shop.
type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string    `gorm:"size:120;not null" json:"name"`
	Email     string    `gorm:"size:254;not null;index:idx_users_shop_email,unique,priority:2" json:"email"`
	Password  string    `gorm:"type:text;not null" json:"-"`
	Role      UserRole  `gorm:"type:user_role;not null" json:"role"`
	ShopID    uuid.UUID `gorm:"type:uuid;not null;index:idx_users_shop_email,unique,priority:1;index" json:"shopID"`
	Shop      Shop      `gorm:"foreignKey:ShopID" json:"-"`
	CreatedAt time.Time `gorm:"not null;default:now()" json:"createdAt"`
}

func (User) TableName() string {
	return "users"
}
