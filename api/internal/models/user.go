package models

import "fmt"

type User struct {
	ModifiableModel
	Username    string  `json:"displayName" example:"jane_doe" binding:"required" gorm:"uniqueIndex;not null"`
	Email       string  `json:"email" example:"jane_doe@mail.com" binding:"required" gorm:"uniqueIndex;not null"`
	FirebaseUID string  `json:"id" example:"HW4beTVvbTUPRxun9MXZxwKPjmC2" binding:"required" gorm:"uniqueIndex"`
	AvatarUrl   *string `json:"avatar" gorm:"type:varchar(255);" example:"https://example.com/avatar.png"`
} // @name User

func (u User) String() string {
	return fmt.Sprintf("%s, #%d", u.Username, u.ID)
}
