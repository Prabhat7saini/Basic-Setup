package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID        int       `gorm:"primaryKey"`
	UUID      uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();not null"`
	Name      *string   `gorm:"type:varchar;not null"`
	Email     string    `gorm:"type:varchar;not null"`
	Password  string    `gorm:"type:varchar;not null"`
	Status    bool
	CreatedAt time.Time `gorm:"autoCreateTime"`
	CreatedBy *int
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	UpdatedBy string
}
