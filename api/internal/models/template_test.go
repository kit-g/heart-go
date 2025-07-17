package models

import (
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
)

func TestTemplate_StructFields(t *testing.T) {
	t.Run("Template with all fields", func(t *testing.T) {
		id := ksuid.New()
		exercises := []TemplateExercise{
			{
				ExerciseID:    "push_up",
				ExerciseOrder: 1,
				Sets: []SetIn{
					{
						Completed: true,
						Weight:    80.0,
						Reps:      10,
					},
				},
			},
		}

		template := Template{
			ModifiableModel: ModifiableModel{
				Model: Model{ID: id},
			},
			Name:          "Upper Body Workout",
			UserID:        "user123",
			OrderInParent: 1,
			Exercises:     exercises,
		}

		assert.Equal(t, id, template.ID)
		assert.Equal(t, "Upper Body Workout", template.Name)
		assert.Equal(t, "user123", template.UserID)
		assert.Equal(t, 1, template.OrderInParent)
		assert.Len(t, template.Exercises, 1)
		assert.Equal(t, "push_up", template.Exercises[0].ExerciseID)
	})

	t.Run("Template with minimal fields", func(t *testing.T) {
		template := Template{
			Name:   "Minimal Template",
			UserID: "user456",
		}

		assert.Equal(t, "Minimal Template", template.Name)
		assert.Equal(t, "user456", template.UserID)
		assert.Equal(t, 0, template.OrderInParent)
		assert.Len(t, template.Exercises, 0)
	})

	t.Run("Template with empty exercises", func(t *testing.T) {
		template := Template{
			Name:      "Empty Template",
			UserID:    "user789",
			Exercises: []TemplateExercise{},
		}

		assert.Equal(t, "Empty Template", template.Name)
		assert.Equal(t, "user789", template.UserID)
		assert.NotNil(t, template.Exercises)
		assert.Len(t, template.Exercises, 0)
	})
}

func TestTemplateExercise_StructFields(t *testing.T) {
	t.Run("TemplateExercise with sets", func(t *testing.T) {
		sets := []SetIn{
			{
				Completed: true,
				Weight:    100.0,
				Reps:      12,
				Duration:  30.0,
				Distance:  0.0,
			},
			{
				Completed: false,
				Weight:    105.0,
				Reps:      10,
				Duration:  35.0,
				Distance:  0.0,
			},
		}

		exercise := TemplateExercise{
			ExerciseID:    "bench_press",
			ExerciseOrder: 2,
			Sets:          sets,
		}

		assert.Equal(t, "bench_press", exercise.ExerciseID)
		assert.Equal(t, 2, exercise.ExerciseOrder)
		assert.Len(t, exercise.Sets, 2)

		assert.True(t, exercise.Sets[0].Completed)
		assert.Equal(t, 100.0, exercise.Sets[0].Weight)
		assert.Equal(t, 12, exercise.Sets[0].Reps)

		assert.False(t, exercise.Sets[1].Completed)
		assert.Equal(t, 105.0, exercise.Sets[1].Weight)
		assert.Equal(t, 10, exercise.Sets[1].Reps)
	})

	t.Run("TemplateExercise with no sets", func(t *testing.T) {
		exercise := TemplateExercise{
			ExerciseID:    "plank",
			ExerciseOrder: 1,
			Sets:          []SetIn{},
		}

		assert.Equal(t, "plank", exercise.ExerciseID)
		assert.Equal(t, 1, exercise.ExerciseOrder)
		assert.Len(t, exercise.Sets, 0)
	})

	t.Run("TemplateExercise with zero order", func(t *testing.T) {
		exercise := TemplateExercise{
			ExerciseID:    "squat",
			ExerciseOrder: 0,
			Sets: []SetIn{
				{
					Completed: true,
					Weight:    50.0,
					Reps:      15,
				},
			},
		}

		assert.Equal(t, "squat", exercise.ExerciseID)
		assert.Equal(t, 0, exercise.ExerciseOrder)
		assert.Len(t, exercise.Sets, 1)
	})

	t.Run("TemplateExercise with cardio sets", func(t *testing.T) {
		exercise := TemplateExercise{
			ExerciseID:    "running",
			ExerciseOrder: 3,
			Sets: []SetIn{
				{
					Completed: true,
					Weight:    0.0,
					Reps:      0,
					Duration:  1800.0, // 30 minutes
					Distance:  5.0,    // 5 km
				},
			},
		}

		assert.Equal(t, "running", exercise.ExerciseID)
		assert.Equal(t, 3, exercise.ExerciseOrder)
		assert.Len(t, exercise.Sets, 1)
		assert.True(t, exercise.Sets[0].Completed)
		assert.Equal(t, 0.0, exercise.Sets[0].Weight)
		assert.Equal(t, 0, exercise.Sets[0].Reps)
		assert.Equal(t, 1800.0, exercise.Sets[0].Duration)
		assert.Equal(t, 5.0, exercise.Sets[0].Distance)
	})
}

func TestTemplateIn_StructFields(t *testing.T) {
	t.Run("TemplateIn with all fields", func(t *testing.T) {
		exercises := []TemplateExercise{
			{
				ExerciseID:    "push_up",
				ExerciseOrder: 1,
				Sets: []SetIn{
					{
						Completed: true,
						Weight:    0.0,
						Reps:      20,
					},
				},
			},
			{
				ExerciseID:    "pull_up",
				ExerciseOrder: 2,
				Sets: []SetIn{
					{
						Completed: false,
						Weight:    0.0,
						Reps:      10,
					},
				},
			},
		}

		templateIn := TemplateIn{
			Name:      "Upper Body Template",
			Order:     1,
			Exercises: exercises,
		}

		assert.Equal(t, "Upper Body Template", templateIn.Name)
		assert.Equal(t, 1, templateIn.Order)
		assert.Len(t, templateIn.Exercises, 2)
		assert.Equal(t, "push_up", templateIn.Exercises[0].ExerciseID)
		assert.Equal(t, "pull_up", templateIn.Exercises[1].ExerciseID)
	})

	t.Run("TemplateIn with minimal fields", func(t *testing.T) {
		templateIn := TemplateIn{
			Name: "Minimal Template",
		}

		assert.Equal(t, "Minimal Template", templateIn.Name)
		assert.Equal(t, 0, templateIn.Order)
		assert.Len(t, templateIn.Exercises, 0)
	})

	t.Run("TemplateIn with empty name", func(t *testing.T) {
		templateIn := TemplateIn{
			Name:  "",
			Order: 5,
			Exercises: []TemplateExercise{
				{
					ExerciseID:    "test_exercise",
					ExerciseOrder: 1,
					Sets:          []SetIn{},
				},
			},
		}

		assert.Equal(t, "", templateIn.Name)
		assert.Equal(t, 5, templateIn.Order)
		assert.Len(t, templateIn.Exercises, 1)
	})
}

func TestNewTemplate(t *testing.T) {
	t.Run("Create template from TemplateIn with exercises", func(t *testing.T) {
		exercises := []TemplateExercise{
			{
				ExerciseID:    "squat",
				ExerciseOrder: 1,
				Sets: []SetIn{
					{
						Completed: true,
						Weight:    80.0,
						Reps:      12,
					},
					{
						Completed: false,
						Weight:    85.0,
						Reps:      10,
					},
				},
			},
			{
				ExerciseID:    "deadlift",
				ExerciseOrder: 2,
				Sets: []SetIn{
					{
						Completed: true,
						Weight:    120.0,
						Reps:      8,
					},
				},
			},
		}

		templateIn := &TemplateIn{
			Name:      "Leg Day",
			Order:     2,
			Exercises: exercises,
		}

		result := NewTemplate(templateIn, "user123")

		assert.Equal(t, "Leg Day", result.Name)
		assert.Equal(t, "user123", result.UserID)
		assert.Equal(t, 0, result.OrderInParent) // Order field is not mapped to OrderInParent
		assert.Len(t, result.Exercises, 2)

		// Check first exercise
		assert.Equal(t, "squat", result.Exercises[0].ExerciseID)
		assert.Equal(t, 1, result.Exercises[0].ExerciseOrder)
		assert.Len(t, result.Exercises[0].Sets, 2)

		// Check second exercise
		assert.Equal(t, "deadlift", result.Exercises[1].ExerciseID)
		assert.Equal(t, 2, result.Exercises[1].ExerciseOrder)
		assert.Len(t, result.Exercises[1].Sets, 1)
	})

	t.Run("Create template from TemplateIn with no exercises", func(t *testing.T) {
		templateIn := &TemplateIn{
			Name:      "Empty Template",
			Order:     1,
			Exercises: []TemplateExercise{},
		}

		result := NewTemplate(templateIn, "user456")

		assert.Equal(t, "Empty Template", result.Name)
		assert.Equal(t, "user456", result.UserID)
		assert.Len(t, result.Exercises, 0)
	})

	t.Run("Create template from TemplateIn with nil exercises", func(t *testing.T) {
		templateIn := &TemplateIn{
			Name:      "Nil Exercises Template",
			Order:     3,
			Exercises: nil,
		}

		result := NewTemplate(templateIn, "user789")

		assert.Equal(t, "Nil Exercises Template", result.Name)
		assert.Equal(t, "user789", result.UserID)
		assert.Nil(t, result.Exercises)
	})

	t.Run("Create template with empty user ID", func(t *testing.T) {
		templateIn := &TemplateIn{
			Name: "Template with Empty User",
		}

		result := NewTemplate(templateIn, "")

		assert.Equal(t, "Template with Empty User", result.Name)
		assert.Equal(t, "", result.UserID)
	})
}

func TestTemplateOut_StructFields(t *testing.T) {
	t.Run("TemplateOut with all fields", func(t *testing.T) {
		exercises := []TemplateExercise{
			{
				ExerciseID:    "bench_press",
				ExerciseOrder: 1,
				Sets: []SetIn{
					{
						Completed: true,
						Weight:    100.0,
						Reps:      8,
					},
				},
			},
		}

		templateOut := TemplateOut{
			ID:        "2ztgx4cIWnxtt95klKnYGGtIfb1",
			Name:      "Chest Workout",
			Order:     1,
			Exercises: exercises,
		}

		assert.Equal(t, "2ztgx4cIWnxtt95klKnYGGtIfb1", templateOut.ID)
		assert.Equal(t, "Chest Workout", templateOut.Name)
		assert.Equal(t, 1, templateOut.Order)
		assert.Len(t, templateOut.Exercises, 1)
		assert.Equal(t, "bench_press", templateOut.Exercises[0].ExerciseID)
	})

	t.Run("TemplateOut with zero values", func(t *testing.T) {
		templateOut := TemplateOut{}

		assert.Equal(t, "", templateOut.ID)
		assert.Equal(t, "", templateOut.Name)
		assert.Equal(t, 0, templateOut.Order)
		assert.Len(t, templateOut.Exercises, 0)
	})
}

func TestNewTemplateOut(t *testing.T) {
	t.Run("Create TemplateOut from Template with exercises", func(t *testing.T) {
		id := ksuid.New()
		exercises := []TemplateExercise{
			{
				ExerciseID:    "push_up",
				ExerciseOrder: 1,
				Sets: []SetIn{
					{
						Completed: true,
						Weight:    0.0,
						Reps:      15,
					},
				},
			},
			{
				ExerciseID:    "dip",
				ExerciseOrder: 2,
				Sets: []SetIn{
					{
						Completed: false,
						Weight:    0.0,
						Reps:      12,
					},
				},
			},
		}

		template := &Template{
			ModifiableModel: ModifiableModel{
				Model: Model{ID: id},
			},
			Name:          "Upper Body",
			UserID:        "user123",
			OrderInParent: 2,
			Exercises:     exercises,
		}

		result := NewTemplateOut(template)

		assert.Equal(t, id.String(), result.ID)
		assert.Equal(t, "Upper Body", result.Name)
		assert.Equal(t, 2, result.Order)
		assert.Len(t, result.Exercises, 2)

		// Check exercises are preserved
		assert.Equal(t, "push_up", result.Exercises[0].ExerciseID)
		assert.Equal(t, 1, result.Exercises[0].ExerciseOrder)
		assert.Len(t, result.Exercises[0].Sets, 1)

		assert.Equal(t, "dip", result.Exercises[1].ExerciseID)
		assert.Equal(t, 2, result.Exercises[1].ExerciseOrder)
		assert.Len(t, result.Exercises[1].Sets, 1)
	})

	t.Run("Create TemplateOut from Template with no exercises", func(t *testing.T) {
		id := ksuid.New()
		template := &Template{
			ModifiableModel: ModifiableModel{
				Model: Model{ID: id},
			},
			Name:          "Empty Template",
			UserID:        "user456",
			OrderInParent: 1,
			Exercises:     []TemplateExercise{},
		}

		result := NewTemplateOut(template)

		assert.Equal(t, id.String(), result.ID)
		assert.Equal(t, "Empty Template", result.Name)
		assert.Equal(t, 1, result.Order)
		assert.Len(t, result.Exercises, 0)
	})

	t.Run("Create TemplateOut from Template with nil exercises", func(t *testing.T) {
		id := ksuid.New()
		template := &Template{
			ModifiableModel: ModifiableModel{
				Model: Model{ID: id},
			},
			Name:          "Nil Exercises Template",
			UserID:        "user789",
			OrderInParent: 0,
			Exercises:     nil,
		}

		result := NewTemplateOut(template)

		assert.Equal(t, id.String(), result.ID)
		assert.Equal(t, "Nil Exercises Template", result.Name)
		assert.Equal(t, 0, result.Order)
		assert.Nil(t, result.Exercises)
	})
}

func TestTemplateResponse_StructFields(t *testing.T) {
	t.Run("TemplateResponse with templates", func(t *testing.T) {
		templates := []TemplateOut{
			{
				ID:    "1",
				Name:  "Template 1",
				Order: 1,
			},
			{
				ID:    "2",
				Name:  "Template 2",
				Order: 2,
			},
		}

		response := TemplateResponse{
			Templates: templates,
		}

		assert.Len(t, response.Templates, 2)
		assert.Equal(t, "Template 1", response.Templates[0].Name)
		assert.Equal(t, "Template 2", response.Templates[1].Name)
	})

	t.Run("TemplateResponse with empty templates", func(t *testing.T) {
		response := TemplateResponse{
			Templates: []TemplateOut{},
		}

		assert.Len(t, response.Templates, 0)
	})

	t.Run("TemplateResponse with nil templates", func(t *testing.T) {
		response := TemplateResponse{
			Templates: nil,
		}

		assert.Nil(t, response.Templates)
	})
}

func TestNewTemplateArray(t *testing.T) {
	t.Run("Create array from multiple templates", func(t *testing.T) {
		id1 := ksuid.New()
		id2 := ksuid.New()
		id3 := ksuid.New()

		templates := []Template{
			{
				ModifiableModel: ModifiableModel{
					Model: Model{ID: id1},
				},
				Name:          "Template 1",
				UserID:        "user123",
				OrderInParent: 1,
				Exercises: []TemplateExercise{
					{
						ExerciseID:    "push_up",
						ExerciseOrder: 1,
						Sets: []SetIn{
							{
								Completed: true,
								Weight:    0.0,
								Reps:      10,
							},
						},
					},
				},
			},
			{
				ModifiableModel: ModifiableModel{
					Model: Model{ID: id2},
				},
				Name:          "Template 2",
				UserID:        "user123",
				OrderInParent: 2,
				Exercises:     []TemplateExercise{},
			},
			{
				ModifiableModel: ModifiableModel{
					Model: Model{ID: id3},
				},
				Name:          "Template 3",
				UserID:        "user456",
				OrderInParent: 1,
				Exercises: []TemplateExercise{
					{
						ExerciseID:    "squat",
						ExerciseOrder: 1,
						Sets: []SetIn{
							{
								Completed: false,
								Weight:    60.0,
								Reps:      12,
							},
						},
					},
					{
						ExerciseID:    "deadlift",
						ExerciseOrder: 2,
						Sets: []SetIn{
							{
								Completed: true,
								Weight:    100.0,
								Reps:      8,
							},
						},
					},
				},
			},
		}

		result := NewTemplateArray(templates)

		assert.Len(t, result, 3)

		// Check first template
		assert.Equal(t, id1.String(), result[0].ID)
		assert.Equal(t, "Template 1", result[0].Name)
		assert.Equal(t, 1, result[0].Order)
		assert.Len(t, result[0].Exercises, 1)

		// Check second template
		assert.Equal(t, id2.String(), result[1].ID)
		assert.Equal(t, "Template 2", result[1].Name)
		assert.Equal(t, 2, result[1].Order)
		assert.Len(t, result[1].Exercises, 0)

		// Check third template
		assert.Equal(t, id3.String(), result[2].ID)
		assert.Equal(t, "Template 3", result[2].Name)
		assert.Equal(t, 1, result[2].Order)
		assert.Len(t, result[2].Exercises, 2)
		assert.Equal(t, "squat", result[2].Exercises[0].ExerciseID)
		assert.Equal(t, "deadlift", result[2].Exercises[1].ExerciseID)
	})

	t.Run("Create array from empty templates", func(t *testing.T) {
		templates := []Template{}
		result := NewTemplateArray(templates)

		assert.Len(t, result, 0)
	})

	t.Run("Create array from single template", func(t *testing.T) {
		id := ksuid.New()
		templates := []Template{
			{
				ModifiableModel: ModifiableModel{
					Model: Model{ID: id},
				},
				Name:          "Single Template",
				UserID:        "user123",
				OrderInParent: 1,
				Exercises: []TemplateExercise{
					{
						ExerciseID:    "plank",
						ExerciseOrder: 1,
						Sets: []SetIn{
							{
								Completed: true,
								Weight:    0.0,
								Reps:      0,
								Duration:  60.0,
							},
						},
					},
				},
			},
		}

		result := NewTemplateArray(templates)

		assert.Len(t, result, 1)
		assert.Equal(t, id.String(), result[0].ID)
		assert.Equal(t, "Single Template", result[0].Name)
		assert.Equal(t, 1, result[0].Order)
		assert.Len(t, result[0].Exercises, 1)
		assert.Equal(t, "plank", result[0].Exercises[0].ExerciseID)
	})
}

func TestTemplateExercise_WithComplexSets(t *testing.T) {
	t.Run("TemplateExercise with mixed set types", func(t *testing.T) {
		exercise := TemplateExercise{
			ExerciseID:    "compound_exercise",
			ExerciseOrder: 1,
			Sets: []SetIn{
				// Strength set
				{
					Completed: true,
					Weight:    100.0,
					Reps:      8,
					Duration:  0.0,
					Distance:  0.0,
				},
				// Cardio set
				{
					Completed: false,
					Weight:    0.0,
					Reps:      0,
					Duration:  300.0, // 5 minutes
					Distance:  1.0,   // 1 km
				},
				// Bodyweight set
				{
					Completed: true,
					Weight:    0.0,
					Reps:      20,
					Duration:  0.0,
					Distance:  0.0,
				},
			},
		}

		assert.Equal(t, "compound_exercise", exercise.ExerciseID)
		assert.Equal(t, 1, exercise.ExerciseOrder)
		assert.Len(t, exercise.Sets, 3)

		// Check strength set
		assert.True(t, exercise.Sets[0].Completed)
		assert.Equal(t, 100.0, exercise.Sets[0].Weight)
		assert.Equal(t, 8, exercise.Sets[0].Reps)

		// Check cardio set
		assert.False(t, exercise.Sets[1].Completed)
		assert.Equal(t, 0.0, exercise.Sets[1].Weight)
		assert.Equal(t, 0, exercise.Sets[1].Reps)
		assert.Equal(t, 300.0, exercise.Sets[1].Duration)
		assert.Equal(t, 1.0, exercise.Sets[1].Distance)

		// Check bodyweight set
		assert.True(t, exercise.Sets[2].Completed)
		assert.Equal(t, 0.0, exercise.Sets[2].Weight)
		assert.Equal(t, 20, exercise.Sets[2].Reps)
	})
}

func TestTemplate_EdgeCases(t *testing.T) {
	t.Run("Template with very long name", func(t *testing.T) {
		longName := "This is a very long template name that might be used in some edge cases where users input extremely long names for their workout templates"

		template := Template{
			Name:   longName,
			UserID: "user123",
		}

		assert.Equal(t, longName, template.Name)
		assert.Equal(t, "user123", template.UserID)
	})

	t.Run("Template with negative order", func(t *testing.T) {
		template := Template{
			Name:          "Negative Order Template",
			UserID:        "user123",
			OrderInParent: -1,
		}

		assert.Equal(t, "Negative Order Template", template.Name)
		assert.Equal(t, -1, template.OrderInParent)
	})

	t.Run("Template with large order number", func(t *testing.T) {
		template := Template{
			Name:          "Large Order Template",
			UserID:        "user123",
			OrderInParent: 999999,
		}

		assert.Equal(t, "Large Order Template", template.Name)
		assert.Equal(t, 999999, template.OrderInParent)
	})
}
