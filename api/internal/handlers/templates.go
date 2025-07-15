package handlers

import (
	"github.com/gin-gonic/gin"
	"heart/internal/dbx"
	"heart/internal/models"
)

func GetTemplates(c *gin.Context, userId string) (any, error) {
	return nil, nil
}

func GetTemplate(c *gin.Context, userId string) (any, error) {
	return nil, nil
}

// MakeTemplate godoc
//
//	@Summary		Creates a workout template
//	@Description	Validates, saves and returns a workout template
//	@Tags			templates
//	@Accept			json
//	@Produce		json
//	@ID				makeTemplate
//	@Param			input	body		TemplateIn	true	"Template request"
//	@Success		200		{object}	Template
//	@Failure		401		{object}	ErrorResponse	"Unauthorized"
//	@Failure		500		{object}	ErrorResponse	"Server error"
//	@Router			/templates [post]
//	@Security		BearerAuth
func MakeTemplate(c *gin.Context, userId string) (any, error) {
	var template models.TemplateIn
	if err := c.BindJSON(&template); err != nil {
		return nil, models.NewValidationError(err)
	}

	created := models.NewTemplate(&template, userId)

	if err := dbx.DB.Create(&created).Error; err != nil {
		return nil, models.NewServerError(err)
	}

	return models.NewTemplateOut(&created), nil
}

func DeleteTemplate(c *gin.Context, userId string) (any, error) {
	return nil, nil
}
