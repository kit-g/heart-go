package models

import (
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        ksuid.KSUID `json:"id" gorm:"type:text;primaryKey;"`
	CreatedAt time.Time   `json:"created_at"`
}

func (m *Model) BeforeCreate(_ *gorm.DB) error {
	if m.ID == ksuid.Nil {
		m.ID = ksuid.New()
	}
	return nil
}

type ModifiableModel struct {
	Model
	UpdatedAt time.Time `json:"updated_at"`
}

type SoftDeleteModel struct {
	ModifiableModel
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
