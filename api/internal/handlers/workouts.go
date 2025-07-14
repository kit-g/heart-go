package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"heart/internal/dbx"
	"heart/internal/models"
)

func GetWorkouts(c *gin.Context, userId string) (any, error) {
	//page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	//limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	//offset := (page - 1) * limit
	//
	//var workouts []models.Workout
	//
	//query := dbx.DB.
	//	Model(&models.Workout{}).
	//	Joins("JOIN users ON users.id = workouts.user_id").
	//	Select("workouts.*").
	//	Where("workouts.user_id = ?", userId).
	//	Preload("Attachments").
	//	Order("created_at desc").
	//	Limit(limit).
	//	Offset(offset)
	//
	//if err := query.Find(&workouts).Error; err != nil {
	//	return nil, models.NewServerError(err)
	//}
	//
	//var out []models.NoteOut
	//total := 0
	//
	//for _, n := range workouts {
	//	if total == 0 {
	//		total = n.NotesCount
	//	}
	//	out = append(out, models.NewNoteOut(&n.Note))
	//}
	//
	//return models.NotesResponse{
	//	Notes: out,
	//	Total: total,
	//}, nil
	return nil, nil
}

func GetWorkout(_ *gin.Context, _ string) (any, error) {
	return nil, nil
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
		return nil, fmt.Errorf("invalid request body: %w", err)
	}

	workout := models.NewWorkout(&workoutIn, userID)

	if err := dbx.DB.Create(&workout).Error; err != nil {
		return nil, models.NewServerError(err)
	}

	return models.NewWorkoutOut(&workout), nil
}

func DeleteWorkout(_ *gin.Context, _ string) (any, error) {
	return nil, nil
}
