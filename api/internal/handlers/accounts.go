package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"heart/internal/dbx"
	"heart/internal/models"
)

func GetAccount(c *gin.Context, userId string) (any, error) {
	return nil, nil
}

// RegisterAccount godoc
//
//	@Summary		Creates an account record
//	@Description	Accounts are managed by Firebase so we just need to store them
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@ID				registerAccount
//	@Param			input	body		User	true	"User request"
//	@Success		201		{object}	User
//	@Failure		401		{object}	ErrorResponse	"Unauthorized"
//	@Failure		500		{object}	ErrorResponse	"Server error"
//	@Router			/accounts [post]
//	@Security		BearerAuth
func RegisterAccount(c *gin.Context, userId string) (any, error) {
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		return nil, models.NewValidationError(err)
	}

	if user.FirebaseUID != userId {
		return nil, models.NewValidationError(fmt.Errorf("user id mismatch"))
	}

	if err := dbx.DB.Create(&user).Error; err != nil {
		return nil, models.NewServerError(err)
	}

	return user, nil
}

func EditAccount(c *gin.Context, userId string) (any, error) {
	return nil, nil
}

func DeleteAccount(c *gin.Context, userId string) (any, error) {
	return nil, nil
}
