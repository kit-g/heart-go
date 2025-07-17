package models

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

type ExercisesResponse struct {
	Exercises []ExerciseOut `json:"exercises"`
} // @name ExercisesResponse
