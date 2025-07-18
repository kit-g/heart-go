package models

import (
	"time"
)

type User struct {
	Username                string     `json:"displayName" example:"jane_doe" binding:"required" gorm:"uniqueIndex;not null"`
	Email                   string     `json:"email" example:"jane_doe@mail.com" binding:"required" gorm:"uniqueIndex;not null"`
	FirebaseUID             string     `json:"id" example:"HW4beTVvbTUPRxun9MXZxwKPjmC2" binding:"required" gorm:"uniqueIndex"`
	AvatarUrl               *string    `json:"avatar" gorm:"type:varchar(255);" example:"https://example.com/avatar.png"`
	AccountDeletionSchedule *string    `json:"accountDeletionSchedule,omitempty" gorm:"type:varchar(255);"`
	ScheduledForDeletionAt  *time.Time `json:"scheduledForDeletionAt,omitempty"`
} // @name User

type EditAccountRequest struct {
	Action   string  `json:"action" example:"removeAvatar" binding:"required"`
	MimeType *string `json:"mimeType,omitempty" example:"image/png"`
} // @name EditAccountRequest

type PresignedUrlResponse struct {
	URL    string            `json:"url"`
	Fields map[string]string `json:"fields"`
} // @name PresignedUrlResponse
