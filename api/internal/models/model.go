package models

import (
	"github.com/segmentio/ksuid"
	"gorm.io/gorm"
	"time"
)

type Model struct {
	ID        ksuid.KSUID `json:"id" example:"1bXoQzZ5H3F4uA1ErL5UyoKAtqW" swaggertype:"string" gorm:"type:text;primaryKey;"`
	CreatedAt time.Time   `json:"createdAt" example:"2025-07-12T12:11:54.450476-04:00"`
}

func (m *Model) BeforeCreate(_ *gorm.DB) error {
	if m.ID == ksuid.Nil {
		m.ID = ksuid.New()
	}
	return nil
}

type ModifiableModel struct {
	Model
	UpdatedAt time.Time `json:"updatedAt" example:"2025-07-15T12:11:54.450476-04:00"`
}

type SoftDeleteModel struct {
	ModifiableModel
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
