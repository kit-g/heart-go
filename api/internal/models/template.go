package models

type Template struct {
	Model
	Name          string
	UserID        string `gorm:"type:text;index"`
	User          User   `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	OrderInParent int    `gorm:"type:integer;default:0"`
}

type TemplateExercise struct {
	Model
	TemplateID int      `gorm:"not null;index"`
	Template   Template `gorm:"foreignKey:TemplateID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ExerciseID string   `gorm:"not null;index"`
	Exercise   Exercise `gorm:"foreignKey:ExerciseID;references:Name;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
