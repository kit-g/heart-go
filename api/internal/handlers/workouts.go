package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"heart/internal/dbx"
	"heart/internal/models"
	"strconv"
)

// GetWorkouts godoc
//
//	@Summary		Returns user workouts
//	@Description	Returns paginated list of user workouts with exercises and sets
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@ID				getWorkouts
//	@Param			pageSize	query		integer	false	"Page size for pagination"
//	@Param			cursor		query		string	false	"Cursor for pagination"
//	@Success		200			{object}	WorkoutResponse
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		500			{object}	ErrorResponse	"Server error"
//	@Router			/workouts [get]
//	@Security		BearerAuth
func GetWorkouts(c *gin.Context, userId string) (any, error) {
	pageSize := 10 // default page size
	if size := c.Query("pageSize"); size != "" {
		if parsed, err := strconv.Atoi(size); err == nil && parsed > 0 {
			pageSize = parsed
		}
	}

	var workouts []models.Workout
	query := dbx.DB.
		Preload("Exercises.Sets").
		Preload("Exercises.Exercise").
		Where("user_id = ?", userId).
		Order("id desc").
		Limit(pageSize)

	if cursor := c.Query("cursor"); cursor != "" {
		query = query.Where("id < ?", cursor)
	}

	if err := query.Find(&workouts).Error; err != nil {
		return nil, models.NewServerError(err)
	}

	var nextCursor string
	if len(workouts) == pageSize {
		nextCursor = workouts[len(workouts)-1].ID.String()
	}

	return models.WorkoutResponse{
		Workouts: models.NewWorkoutsArray(workouts),
		Cursor:   nextCursor,
	}, nil
}

// GetWorkout godoc
//
//	@Summary		Returns a workout
//	@Description	Returns a single workout by ID with exercises and sets
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@ID				getWorkout
//	@Param			workoutId	path		string	true	"Workout ID"
//	@Success		200			{object}	Workout
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Not Found"
//	@Failure		500			{object}	ErrorResponse	"Server error"
//	@Router			/workouts/{workoutId} [get]
//	@Security		BearerAuth
func GetWorkout(c *gin.Context, userId string) (any, error) {
	workoutId := c.Param("workoutId")

	var workout models.Workout
	if err := dbx.DB.
		Preload("Exercises.Sets").
		Preload("Exercises.Exercise").
		Where("id = ? AND user_id = ?", workoutId, userId).
		First(&workout).Error; err != nil {
		if models.Is(err, gorm.ErrRecordNotFound) {
			return nil, models.NewNotFoundError("Workout not found", err)
		}
		return nil, models.NewServerError(err)
	}

	return models.NewWorkoutOut(&workout), nil
}

// MakeWorkout godoc
//
//	@Summary		Creates a workout
//	@Description	Validates, saves and returns a workout
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@ID				makeWorkout
//	@Param			input	body		WorkoutIn	true	"Workout request"
//	@Success		200		{object}	Workout
//	@Failure		401		{object}	ErrorResponse	"Unauthorized"
//	@Failure		500		{object}	ErrorResponse	"Server error"
//	@Router			/workouts [post]
//	@Security		BearerAuth
func MakeWorkout(c *gin.Context, userID string) (any, error) {
	var workoutIn models.WorkoutIn
	if err := c.BindJSON(&workoutIn); err != nil {
		return nil, models.NewValidationError(err)
	}

	workout := models.NewWorkout(&workoutIn, userID)

	if err := dbx.DB.Create(&workout).Error; err != nil {
		return nil, models.NewServerError(err)
	}

	return models.NewWorkoutOut(&workout), nil
}

// DeleteWorkout godoc
//
//	@Summary		Deletes a workout
//	@Description	Deletes a workout by ID
//	@Tags			workouts
//	@Accept			json
//	@Produce		json
//	@ID				deleteWorkout
//	@Param			workoutId	path		string	true	"Workout ID"
//	@Success		204		"No Content"
//	@Failure		401			{object}	ErrorResponse	"Unauthorized"
//	@Failure		404			{object}	ErrorResponse	"Not Found"
//	@Failure		500			{object}	ErrorResponse	"Server error"
//	@Router			/workouts/{workoutId} [delete]
//	@Security		BearerAuth
func DeleteWorkout(c *gin.Context, userId string) (any, error) {
	workoutId := c.Param("workoutId")

	result := dbx.DB.
		Where("id = ? AND user_id = ?", workoutId, userId).
		Delete(&models.Workout{})

	if result.Error != nil {
		return nil, models.NewServerError(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, models.NewNotFoundError("Workout not found", errors.New("workout not found"))
	}

	return models.NoContent, nil
}
