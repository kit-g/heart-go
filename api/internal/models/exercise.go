package models

type Exercise struct {
	PK           string            `dynamodbav:"pk"` // always "EXERCISE"
	Name         string            `gorm:"primaryKey;not null" dynamodbav:"sk"`
	Category     string            `gorm:"not null" dynamodbav:"category"`
	Target       string            `gorm:"not null" dynamodbav:"target"`
	Asset        *ImageDescription `gorm:"type:text" dynamodbav:"asset,omitempty"`
	Thumbnail    *ImageDescription `gorm:"type:text" dynamodbav:"thumbnail,omitempty"`
	Instructions *string           `gorm:"type:text" dynamodbav:"instructions,omitempty"`
	UserID       string            `gorm:"type:text" dynamodbav:"userId,omitempty"`
}

func (e *Exercise) String() string {
	return e.Name
}

type ExerciseOut struct {
	Name         string            `json:"name" example:"Push Up"`
	Category     string            `json:"category" example:"Body weight"`
	Target       string            `json:"target" example:"Chest"`
	Asset        *ImageDescription `json:"asset,omitempty"`
	Thumbnail    *ImageDescription `json:"thumbnail,omitempty"`
	Instructions *string           `json:"instructions,omitempty" example:"Keep your body straight and lower yourself until your chest almost touches the ground."`
} // @name Exercise

func NewExerciseOut(e *Exercise) ExerciseOut {
	return ExerciseOut{
		Name:         e.Name,
		Category:     e.Category,
		Target:       e.Target,
		Asset:        e.Asset,
		Thumbnail:    e.Thumbnail,
		Instructions: e.Instructions,
	}
}

type ExercisesResponse struct {
	Exercises []ExerciseOut `json:"exercises"`
} // @name ExercisesResponse
