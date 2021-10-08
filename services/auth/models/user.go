package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deletedAt"`
	Avatar    string         `json:"avatar"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Password  string         `json:"-"`
	Version   int            `json:"version"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	id := uuid.Must(uuid.NewV4())
	u.ID = id.String()
	return nil
}
