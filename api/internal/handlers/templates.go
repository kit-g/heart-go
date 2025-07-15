package handlers

import (
	"github.com/gin-gonic/gin"
	"heart/internal/dbx"
	"heart/internal/models"
)

// GetTemplates godoc
//
//	@Summary		Lists workout templates
//	@Description	Returns all workout templates for the authenticated user
//	@Tags			templates
//	@Accept			json
//	@Produce		json
//	@ID				getTemplates
//	@Success		200	{array}		Template
//	@Failure		401	{object}	ErrorResponse	"Unauthorized"
//	@Failure		500	{object}	ErrorResponse	"Server error"
//	@Router			/templates [get]
//	@Security		BearerAuth
func GetTemplates(_ *gin.Context, userId string) (any, error) {
	var templates []models.Template
	query := dbx.DB.
		Where("user_id = ?", userId).
		Order("id desc")

	if err := query.Find(&templates).Error; err != nil {
		return nil, models.NewServerError(err)
	}

	return models.TemplateResponse{
		Templates: models.NewTemplateArray(templates),
	}, nil
}

// GetTemplate godoc
//
//	@Summary		Get workout template
//	@Description	Returns a specific workout template by ID
//	@Tags			templates
//	@Accept			json
//	@Produce		json
//	@ID				getTemplate
//	@Param			id	path		string	true	"Template ID"
//	@Success		200	{object}	Template
//	@Failure		401	{object}	ErrorResponse	"Unauthorized"
//	@Failure		404	{object}	ErrorResponse	"Not Found"
//	@Failure		500	{object}	ErrorResponse	"Server error"
//	@Router			/templates/{id} [get]
//	@Security		BearerAuth
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

// DeleteTemplate godoc
//
//	@Summary		Delete workout template
//	@Description	Deletes a specific workout template by ID
//	@Tags			templates
//	@Accept			json
//	@Produce		json
//	@ID				deleteTemplate
//	@Param			id	path	string	true	"Template ID"
//	@Success		204	"No Content"
//	@Failure		401	{object}	ErrorResponse	"Unauthorized"
//	@Failure		404	{object}	ErrorResponse	"Not Found"
//	@Failure		500	{object}	ErrorResponse	"Server error"
//	@Router			/templates/{id} [delete]
//	@Security		BearerAuth
func DeleteTemplate(c *gin.Context, userId string) (any, error) {
	return nil, nil
}
