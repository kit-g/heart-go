package models

import "time"

type Exercise struct {
	Name            string `gorm:"primaryKey;not null"`
	Category        string `gorm:"not null"`
	Target          string `gorm:"not null"`
	Asset           string `gorm:"type:text"`
	AssetWidth      int    `gorm:"type:integer"`
	AssetHeight     int    `gorm:"type:integer"`
	Thumbnail       string `gorm:"type:text"`
	ThumbnailWidth  int    `gorm:"type:integer"`
	ThumbnailHeight int    `gorm:"type:integer"`
	Instructions    string `gorm:"type:text"`
	UserID          string `gorm:"type:text"`
}

func (e *Exercise) String() string {
	return e.Name
}

type Workout struct {
	SoftDeleteModel
	Start  time.Time `gorm:"not null"`
	End    time.Time `gorm:"not null"`
	UserID string    `gorm:"type:text;index"`
	User   User      `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Name   string    `gorm:"type:text"`
}

func (w *Workout) String() string {
	return w.Name
}

type WorkoutExercise struct {
	Model
	WorkoutID     string   `gorm:"not null;index"`
	Workout       Workout  `gorm:"foreignKey:WorkoutID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ExerciseID    string   `gorm:"not null;index"`
	Exercise      Exercise `gorm:"foreignKey:ExerciseID;references:Name;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ExerciseOrder int      `gorm:"type:integer"`
}

func (we *WorkoutExercise) String() string {
	return we.Exercise.Name + " in " + we.Workout.Name
}

type Set struct {
	ModifiableModel
	WorkoutExerciseID string          `gorm:"not null;index"`
	WorkoutExercise   WorkoutExercise `gorm:"foreignKey:WorkoutExerciseID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Completed         bool            `gorm:"not null;default:false"`
	Weight            float64         `gorm:"type:real"` // kgs
	Reps              int             `gorm:"type:integer"`
	Duration          float64         `gorm:"type:real"` // seconds
	Distance          float64         `gorm:"type:real"` // kilometers
}

type ImageDescription struct {
	Link   *string `json:"link" example:"https://example.com/image.jpg"`
	Width  *int    `json:"width" example:"100"`
	Height *int    `json:"height" example:"100"`
} // @name ImageDescription

type ExerciseOut struct {
	Name         string            `json:"name" example:"Push Up"`
	Category     string            `json:"category" example:"Body weight"`
	Target       string            `json:"target" example:"Chest"`
	Asset        *ImageDescription `json:"asset,omitempty"`
	Thumbnail    *ImageDescription `json:"thumbnail,omitempty"`
	Instructions *string           `json:"instructions,omitempty" example:"Keep your body straight and lower yourself until your chest almost touches the ground."`
} // @name Exercise

func NewExerciseOut(e Exercise) ExerciseOut {
	var asset *ImageDescription
	if e.Asset != "" {
		asset = &ImageDescription{
			Link:   &e.Asset,
			Width:  &e.AssetWidth,
			Height: &e.AssetHeight,
		}
	}

	var thumbnail *ImageDescription
	if e.Thumbnail != "" {
		thumbnail = &ImageDescription{
			Link:   &e.Thumbnail,
			Width:  &e.ThumbnailWidth,
			Height: &e.ThumbnailHeight,
		}
	}

	return ExerciseOut{
		Name:         e.Name,
		Category:     e.Category,
		Target:       e.Target,
		Asset:        asset,
		Thumbnail:    thumbnail,
		Instructions: &e.Instructions,
	}
}

type ExercisesResponse struct {
	Exercises []ExerciseOut `json:"exercises"`
} // @name ExercisesResponse
