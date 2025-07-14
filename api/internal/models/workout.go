package models

import (
	"github.com/segmentio/ksuid"
	"log"
	"time"
)

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
	Start     time.Time         `gorm:"not null"`
	End       time.Time         `gorm:"not null"`
	UserID    string            `gorm:"type:text;index"`
	User      User              `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL"`
	Name      string            `gorm:"type:text"`
	Exercises []WorkoutExercise `gorm:"foreignKey:WorkoutID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
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
	Sets          []Set
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

func NewExerciseOut(e *Exercise) ExerciseOut {
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

type SetIn struct {
	Completed bool    `json:"completed" binding:"required" example:"true"`
	Weight    float64 `json:"weight,omitempty" example:"100"`
	Reps      int     `json:"reps,omitempty" example:"10"`
	Duration  float64 `json:"duration,omitempty" example:"10"`
	Distance  float64 `json:"distance,omitempty" example:"10"`
}

type WorkoutExerciseIn struct {
	Exercise string  `json:"exercise" binding:"required" example:"Push Up"`
	Sets     []SetIn `json:"sets"`
	Order    int     `json:"order" example:"1"`
}

type WorkoutIn struct {
	ID        string              `json:"id" binding:"required" example:"2zsp6iMWgOx9n6qQxZm0GmeXog1"`
	Name      string              `json:"name" example:"Legs"`
	Start     time.Time           `json:"start" binding:"required" example:"2023-01-01T12:00:00Z"`
	End       time.Time           `json:"end" example:"2023-01-01T12:00:00Z"`
	Exercises []WorkoutExerciseIn `json:"exercises" binding:"required"`
} // @name WorkoutIn

type SetOut struct {
	ID        string  `json:"id" binding:"required" example:"1234567890"`
	Completed bool    `json:"completed" binding:"required" example:"true"`
	Weight    float64 `json:"weight" example:"100"`
	Reps      int     `json:"reps" example:"10"`
	Duration  float64 `json:"duration" example:"10"`
	Distance  float64 `json:"distance" example:"10"`
} // @name Set

type WorkoutExerciseOut struct {
	Exercise *string  `json:"exercise" example:"Push Up"`
	Sets     []SetOut `json:"sets"`
} // @name WorkoutExercise

type WorkoutOut struct {
	ID        string               `json:"id" example:"2zsp6iMWgOx9n6qQxZm0GmeXog1"`
	Name      string               `json:"name" example:"Legs"`
	Start     time.Time            `json:"start" example:"2023-01-01T12:00:00Z"`
	End       time.Time            `json:"end" example:"2023-01-01T12:00:00Z"`
	Exercises []WorkoutExerciseOut `json:"exercises"`
} // @name Workout

func NewSetOut(s *Set) SetOut {
	return SetOut{
		ID:        s.ID.String(),
		Completed: s.Completed,
		Weight:    s.Weight,
		Reps:      s.Reps,
		Duration:  s.Duration,
		Distance:  s.Distance,
	}
}

func NewWorkoutExerciseOut(e *WorkoutExercise) WorkoutExerciseOut {
	sets := make([]SetOut, len(e.Sets))

	for i, s := range e.Sets {
		sets[i] = NewSetOut(&s)
	}

	return WorkoutExerciseOut{
		Exercise: &e.ExerciseID,
		Sets:     sets,
	}
}

func NewWorkout(w *WorkoutIn, userId string) Workout {
	id, err := ksuid.Parse(w.ID)

	if err != nil {
		log.Fatal(err)
	}

	workout := Workout{
		SoftDeleteModel: SoftDeleteModel{
			ModifiableModel: ModifiableModel{
				Model: Model{ID: id},
			},
		},
		Name:      w.Name,
		Start:     w.Start,
		End:       w.End,
		UserID:    userId,
		Exercises: make([]WorkoutExercise, len(w.Exercises)),
	}

	for i, exercise := range w.Exercises {
		workout.Exercises[i] = WorkoutExercise{
			ExerciseID:    exercise.Exercise,
			ExerciseOrder: exercise.Order,
			Sets:          make([]Set, len(exercise.Sets)),
		}

		for j, set := range exercise.Sets {
			workout.Exercises[i].Sets[j] = Set{
				Completed: set.Completed,
				Weight:    set.Weight,
				Reps:      set.Reps,
				Duration:  set.Duration,
				Distance:  set.Distance,
			}
		}
	}

	return workout
}

func NewWorkoutOut(w *Workout) WorkoutOut {
	exercises := make([]WorkoutExerciseOut, len(w.Exercises))

	for i, e := range w.Exercises {
		exercises[i] = NewWorkoutExerciseOut(&e)
	}

	return WorkoutOut{
		ID:        w.ID.String(),
		Name:      w.Name,
		Start:     w.Start,
		End:       w.End,
		Exercises: exercises,
	}
}

type ExercisesResponse struct {
	Exercises []ExerciseOut `json:"exercises"`
} // @name ExercisesResponse

type WorkoutResponse struct {
	Workouts []WorkoutOut `json:"workouts"`
} // @name WorkoutResponse
