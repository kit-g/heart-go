package models

type FeedbackRequest struct {
	Message string `json:"message" example:"Good job!" binding:"required"`
}
