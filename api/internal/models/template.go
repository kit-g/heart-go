package models

type Template struct {
	ModifiableModel
	Name          string
	UserID        string             `gorm:"type:text;index"`
	User          User               `json:"user" gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	OrderInParent int                `gorm:"type:integer;default:0"`
	Exercises     []TemplateExercise `json:"exercises" gorm:"type:jsonb;serializer:json"`
}

type TemplateExercise struct {
	ExerciseID    string  `json:"exercise"`
	ExerciseOrder int     `json:"order"`
	Sets          []SetIn `json:"sets"`
} // @name TemplateExercise

type TemplateIn struct {
	Name      string             `json:"name"`
	Order     int                `json:"order"`
	Exercises []TemplateExercise `json:"exercises"`
} // @name TemplateIn

func NewTemplate(t *TemplateIn, userId string) Template {
	return Template{
		Name:      t.Name,
		UserID:    userId,
		Exercises: t.Exercises,
	}
}

type TemplateOut struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	Order     int                `json:"order"`
	Exercises []TemplateExercise `json:"exercises"`
} // @name Template

func NewTemplateOut(t *Template) TemplateOut {
	return TemplateOut{
		ID:        t.ID.String(),
		Name:      t.Name,
		Order:     t.OrderInParent,
		Exercises: t.Exercises,
	}
}
