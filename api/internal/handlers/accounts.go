package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"heart/internal/awsx"
	"heart/internal/config"
	"heart/internal/dbx"
	"heart/internal/models"
)

// GetAccount godoc
//
//	@Summary		Get user account
//	@Description	Returns user account information for the authenticated user
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@ID				getAccount
//	@Param			accountId	path		string	true	"Account ID"
//	@Success		200			{object}	User
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Not Found"
//	@Failure		500			{object}	ErrorResponse	"Server error"
//	@Router			/accounts/{accountId} [get]
//	@Security		BearerAuth
func GetAccount(_ *gin.Context, userId string) (any, error) {
	var user models.User

	if err := dbx.DB.
		Where("firebase_uid = ?", userId).
		First(&user).
		Error; err != nil {
		return nil, models.NewServerError(err)
	}

	return user, nil
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
	var request models.EditAccountRequest
	if err := c.BindJSON(&request); err != nil {
		return nil, models.NewValidationError(err)
	}

	switch request.Action {
	case "undoAccountDeletion":

	case "removeAvatar":
		_, err := awsx.DeleteObject(
			c.Request.Context(),
			config.App.MediaBucket,
			config.App.AvatarKey(userId),
		)

		if err != nil {
			return nil, models.NewServerError(err)
		}

		if err := dbx.DB.
			Model(&models.User{}).
			Where("firebase_uid = ?", userId).
			Update("avatar_url", nil).
			Error; err != nil {
			return nil, models.NewServerError(err)
		}

	case "uploadAvatar":
		var mimeType string
		if request.MimeType == nil || *request.MimeType == "" {
			mimeType = defaultMimeType
		} else {
			mimeType = *request.MimeType
		}

		response, err := awsx.GeneratePresignedPostURL(
			c.Request.Context(),
			config.App.UploadBucket,
			config.App.AvatarKey(userId),
			mimeType,
			config.App.UploadDestinationTag(),
		)

		if err != nil {
			return nil, models.NewServerError(err)
		}

		return models.PresignedUrlResponse{
			URL:    response.URL,
			Fields: response.Values,
		}, nil
	}
	return nil, nil
}

// DeleteAccount godoc
//
//	@Summary		Delete user account
//	@Description	Schedules account deletion for the authenticated user
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@ID				deleteAccount
//	@Param			accountId	path	string	true	"Account ID"
//	@Success		204			"No Content"
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Not Found"
//	@Failure		500			{object}	ErrorResponse	"Server error"
//	@Router			/accounts [delete]
//	@Security		BearerAuth
func DeleteAccount(c *gin.Context, userId string) (any, error) {
	when, schedule, err := awsx.CreateAccountDeletionSchedule(c.Request.Context(), userId)

	if err != nil {
		return nil, models.NewServerError(err)
	}

	if when != nil && schedule != nil {
		columns := map[string]interface{}{
			"account_deletion_schedule": schedule,
			"scheduled_for_deletion_at": when,
		}

		if err := dbx.DB.
			Model(&models.User{}).
			Where("firebase_uid = ?", userId).
			Updates(columns).
			Error; err != nil {
			return nil, models.NewServerError(err)
		}
	}

	return models.NoContent, nil
}

const defaultMimeType = "image/png"
