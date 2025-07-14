package handlers

import (
	"github.com/gin-gonic/gin"
	"heart/internal/dbx"
	"heart/internal/models"
)

// GetExercises godoc
//
//	@Summary		List all exercises
//	@Description	Returns all exercises in a single page
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@ID				getExercises
//	@Success		200	{object}	ExercisesResponse
//	@Failure		401	{object}	ErrorResponse	"Unauthorized"
//	@Failure		500	{object}	ErrorResponse	"Server error"
//	@Router			/exercises [get]
//	@Security		BearerAuth
func GetExercises(_ *gin.Context, _ string) (any, error) {
	var exercises []models.Exercise

	query := dbx.DB.Model(&models.Exercise{})

	if err := query.Find(&exercises).Error; err != nil {
		return nil, models.NewServerError(err)
	}

	response := models.ExercisesResponse{
		Exercises: make([]models.ExerciseOut, len(exercises)),
	}

	for i, e := range exercises {
		response.Exercises[i] = models.NewExerciseOut(&e)
	}

	return response, nil
}
