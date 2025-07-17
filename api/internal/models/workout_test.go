package models

import (
	"testing"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWorkout_String(t *testing.T) {
	tests := []struct {
		name     string
		workout  Workout
		expected string
	}{
		{
			name: "Workout with name",
			workout: Workout{
				Name: "Leg Day",
			},
			expected: "Leg Day",
		},
		{
			name: "Workout with empty name",
			workout: Workout{
				Name: "",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.workout.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestWorkoutExercise_String(t *testing.T) {
	tests := []struct {
		name     string
		we       WorkoutExercise
		expected string
	}{
		{
			name: "WorkoutExercise with exercise and workout names",
			we: WorkoutExercise{
				Exercise: Exercise{Name: "Push Up"},
				Workout:  Workout{Name: "Upper Body"},
			},
			expected: "Push Up in Upper Body",
		},
		{
			name: "WorkoutExercise with empty names",
			we: WorkoutExercise{
				Exercise: Exercise{Name: ""},
				Workout:  Workout{Name: ""},
			},
			expected: " in ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.we.String()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNewSetOut(t *testing.T) {
	t.Run("Set with all fields", func(t *testing.T) {
		id := ksuid.New()
		set := &Set{
			ModifiableModel: ModifiableModel{
				Model: Model{ID: id},
			},
			Completed: true,
			Weight:    100.5,
			Reps:      12,
			Duration:  45.5,
			Distance:  2.5,
		}

		result := NewSetOut(set)

		assert.Equal(t, id.String(), result.ID)
		assert.True(t, result.Completed)
		assert.Equal(t, 100.5, result.Weight)
		assert.Equal(t, 12, result.Reps)
		assert.Equal(t, 45.5, result.Duration)
		assert.Equal(t, 2.5, result.Distance)
	})

	t.Run("Set with minimal fields", func(t *testing.T) {
		id := ksuid.New()
		set := &Set{
			ModifiableModel: ModifiableModel{
				Model: Model{ID: id},
			},
			Completed: false,
		}

		result := NewSetOut(set)

		assert.Equal(t, id.String(), result.ID)
		assert.False(t, result.Completed)
		assert.Equal(t, 0.0, result.Weight)
		assert.Equal(t, 0, result.Reps)
		assert.Equal(t, 0.0, result.Duration)
		assert.Equal(t, 0.0, result.Distance)
	})
}

func TestNewWorkoutExerciseOut(t *testing.T) {
	t.Run("WorkoutExercise with sets", func(t *testing.T) {
		id1 := ksuid.New()
		id2 := ksuid.New()

		we := &WorkoutExercise{
			ExerciseID: "Push Up",
			Sets: []Set{
				{
					ModifiableModel: ModifiableModel{
						Model: Model{ID: id1},
					},
					Completed: true,
					Weight:    80.0,
					Reps:      10,
				},
				{
					ModifiableModel: ModifiableModel{
						Model: Model{ID: id2},
					},
					Completed: false,
					Weight:    85.0,
					Reps:      8,
				},
			},
		}

		result := NewWorkoutExerciseOut(we)

		require.NotNil(t, result.Exercise)
		assert.Equal(t, "Push Up", *result.Exercise)
		assert.Len(t, result.Sets, 2)

		assert.Equal(t, id1.String(), result.Sets[0].ID)
		assert.True(t, result.Sets[0].Completed)
		assert.Equal(t, 80.0, result.Sets[0].Weight)
		assert.Equal(t, 10, result.Sets[0].Reps)

		assert.Equal(t, id2.String(), result.Sets[1].ID)
		assert.False(t, result.Sets[1].Completed)
		assert.Equal(t, 85.0, result.Sets[1].Weight)
		assert.Equal(t, 8, result.Sets[1].Reps)
	})

	t.Run("WorkoutExercise with no sets", func(t *testing.T) {
		we := &WorkoutExercise{
			ExerciseID: "Plank",
			Sets:       []Set{},
		}

		result := NewWorkoutExerciseOut(we)

		require.NotNil(t, result.Exercise)
		assert.Equal(t, "Plank", *result.Exercise)
		assert.Len(t, result.Sets, 0)
	})
}

func TestNewWorkout(t *testing.T) {
	t.Run("Valid workout with exercises and sets", func(t *testing.T) {
		id := ksuid.New()
		start := time.Now()
		end := start.Add(time.Hour)

		workoutIn := &WorkoutIn{
			ID:    id.String(),
			Name:  "Test Workout",
			Start: start,
			End:   end,
			Exercises: []WorkoutExerciseIn{
				{
					Exercise: "Push Up",
					Order:    1,
					Sets: []SetIn{
						{
							Completed: true,
							Weight:    80.0,
							Reps:      10,
							Duration:  30.0,
							Distance:  0.0,
						},
						{
							Completed: false,
							Weight:    85.0,
							Reps:      8,
							Duration:  25.0,
							Distance:  0.0,
						},
					},
				},
				{
					Exercise: "Squat",
					Order:    2,
					Sets: []SetIn{
						{
							Completed: true,
							Weight:    100.0,
							Reps:      12,
							Duration:  40.0,
							Distance:  0.0,
						},
					},
				},
			},
		}

		result := NewWorkout(workoutIn, "user123")

		assert.Equal(t, id, result.ID)
		assert.Equal(t, "Test Workout", result.Name)
		assert.Equal(t, start, result.Start)
		assert.Equal(t, end, result.End)
		assert.Equal(t, "user123", result.UserID)
		assert.Len(t, result.Exercises, 2)

		// Check first exercise
		assert.Equal(t, "Push Up", result.Exercises[0].ExerciseID)
		assert.Equal(t, 1, result.Exercises[0].ExerciseOrder)
		assert.Len(t, result.Exercises[0].Sets, 2)

		assert.True(t, result.Exercises[0].Sets[0].Completed)
		assert.Equal(t, 80.0, result.Exercises[0].Sets[0].Weight)
		assert.Equal(t, 10, result.Exercises[0].Sets[0].Reps)
		assert.Equal(t, 30.0, result.Exercises[0].Sets[0].Duration)

		assert.False(t, result.Exercises[0].Sets[1].Completed)
		assert.Equal(t, 85.0, result.Exercises[0].Sets[1].Weight)
		assert.Equal(t, 8, result.Exercises[0].Sets[1].Reps)
		assert.Equal(t, 25.0, result.Exercises[0].Sets[1].Duration)

		// Check second exercise
		assert.Equal(t, "Squat", result.Exercises[1].ExerciseID)
		assert.Equal(t, 2, result.Exercises[1].ExerciseOrder)
		assert.Len(t, result.Exercises[1].Sets, 1)

		assert.True(t, result.Exercises[1].Sets[0].Completed)
		assert.Equal(t, 100.0, result.Exercises[1].Sets[0].Weight)
		assert.Equal(t, 12, result.Exercises[1].Sets[0].Reps)
		assert.Equal(t, 40.0, result.Exercises[1].Sets[0].Duration)
	})

	t.Run("Workout with no exercises", func(t *testing.T) {
		id := ksuid.New()
		start := time.Now()

		workoutIn := &WorkoutIn{
			ID:        id.String(),
			Name:      "Empty Workout",
			Start:     start,
			End:       start.Add(time.Hour),
			Exercises: []WorkoutExerciseIn{},
		}

		result := NewWorkout(workoutIn, "user456")

		assert.Equal(t, id, result.ID)
		assert.Equal(t, "Empty Workout", result.Name)
		assert.Equal(t, "user456", result.UserID)
		assert.Len(t, result.Exercises, 0)
	})
}

func TestNewWorkoutOut(t *testing.T) {
	t.Run("Workout with exercises and sets", func(t *testing.T) {
		id := ksuid.New()
		setId1 := ksuid.New()
		setId2 := ksuid.New()
		start := time.Now()
		end := start.Add(time.Hour)

		workout := &Workout{
			SoftDeleteModel: SoftDeleteModel{
				ModifiableModel: ModifiableModel{
					Model: Model{ID: id},
				},
			},
			Name:  "Test Workout",
			Start: start,
			End:   end,
			Exercises: []WorkoutExercise{
				{
					ExerciseID:    "Push Up",
					ExerciseOrder: 1,
					Sets: []Set{
						{
							ModifiableModel: ModifiableModel{
								Model: Model{ID: setId1},
							},
							Completed: true,
							Weight:    80.0,
							Reps:      10,
						},
						{
							ModifiableModel: ModifiableModel{
								Model: Model{ID: setId2},
							},
							Completed: false,
							Weight:    85.0,
							Reps:      8,
						},
					},
				},
			},
		}

		result := NewWorkoutOut(workout)

		assert.Equal(t, id.String(), result.ID)
		assert.Equal(t, "Test Workout", result.Name)
		assert.Equal(t, start, result.Start)
		assert.Equal(t, end, result.End)
		assert.Len(t, result.Exercises, 1)

		assert.Equal(t, "Push Up", *result.Exercises[0].Exercise)
		assert.Len(t, result.Exercises[0].Sets, 2)

		assert.Equal(t, setId1.String(), result.Exercises[0].Sets[0].ID)
		assert.True(t, result.Exercises[0].Sets[0].Completed)
		assert.Equal(t, 80.0, result.Exercises[0].Sets[0].Weight)
		assert.Equal(t, 10, result.Exercises[0].Sets[0].Reps)

		assert.Equal(t, setId2.String(), result.Exercises[0].Sets[1].ID)
		assert.False(t, result.Exercises[0].Sets[1].Completed)
		assert.Equal(t, 85.0, result.Exercises[0].Sets[1].Weight)
		assert.Equal(t, 8, result.Exercises[0].Sets[1].Reps)
	})

	t.Run("Workout with no exercises", func(t *testing.T) {
		id := ksuid.New()
		start := time.Now()

		workout := &Workout{
			SoftDeleteModel: SoftDeleteModel{
				ModifiableModel: ModifiableModel{
					Model: Model{ID: id},
				},
			},
			Name:      "Empty Workout",
			Start:     start,
			End:       start.Add(time.Hour),
			Exercises: []WorkoutExercise{},
		}

		result := NewWorkoutOut(workout)

		assert.Equal(t, id.String(), result.ID)
		assert.Equal(t, "Empty Workout", result.Name)
		assert.Len(t, result.Exercises, 0)
	})
}

func TestNewWorkoutsArray(t *testing.T) {
	t.Run("Array of workouts", func(t *testing.T) {
		id1 := ksuid.New()
		id2 := ksuid.New()
		start := time.Now()

		workouts := []Workout{
			{
				SoftDeleteModel: SoftDeleteModel{
					ModifiableModel: ModifiableModel{
						Model: Model{ID: id1},
					},
				},
				Name:      "Workout 1",
				Start:     start,
				End:       start.Add(time.Hour),
				Exercises: []WorkoutExercise{},
			},
			{
				SoftDeleteModel: SoftDeleteModel{
					ModifiableModel: ModifiableModel{
						Model: Model{ID: id2},
					},
				},
				Name:      "Workout 2",
				Start:     start.Add(time.Hour),
				End:       start.Add(2 * time.Hour),
				Exercises: []WorkoutExercise{},
			},
		}

		result := NewWorkoutsArray(workouts)

		assert.Len(t, result, 2)
		assert.Equal(t, id1.String(), result[0].ID)
		assert.Equal(t, "Workout 1", result[0].Name)
		assert.Equal(t, id2.String(), result[1].ID)
		assert.Equal(t, "Workout 2", result[1].Name)
	})

	t.Run("Empty array", func(t *testing.T) {
		workouts := []Workout{}
		result := NewWorkoutsArray(workouts)
		assert.Len(t, result, 0)
	})
}

func TestSetIn_Validation(t *testing.T) {
	t.Run("Valid SetIn", func(t *testing.T) {
		setIn := SetIn{
			Completed: true,
			Weight:    100.0,
			Reps:      12,
			Duration:  30.0,
			Distance:  2.5,
		}

		assert.True(t, setIn.Completed)
		assert.Equal(t, 100.0, setIn.Weight)
		assert.Equal(t, 12, setIn.Reps)
		assert.Equal(t, 30.0, setIn.Duration)
		assert.Equal(t, 2.5, setIn.Distance)
	})

	t.Run("SetIn with zero values", func(t *testing.T) {
		setIn := SetIn{
			Completed: false,
		}

		assert.False(t, setIn.Completed)
		assert.Equal(t, 0.0, setIn.Weight)
		assert.Equal(t, 0, setIn.Reps)
		assert.Equal(t, 0.0, setIn.Duration)
		assert.Equal(t, 0.0, setIn.Distance)
	})
}

func TestWorkoutExerciseIn_Validation(t *testing.T) {
	t.Run("Valid WorkoutExerciseIn", func(t *testing.T) {
		we := WorkoutExerciseIn{
			Exercise: "Push Up",
			Order:    1,
			Sets: []SetIn{
				{
					Completed: true,
					Weight:    80.0,
					Reps:      10,
				},
			},
		}

		assert.Equal(t, "Push Up", we.Exercise)
		assert.Equal(t, 1, we.Order)
		assert.Len(t, we.Sets, 1)
		assert.True(t, we.Sets[0].Completed)
	})

	t.Run("WorkoutExerciseIn with empty sets", func(t *testing.T) {
		we := WorkoutExerciseIn{
			Exercise: "Plank",
			Order:    2,
			Sets:     []SetIn{},
		}

		assert.Equal(t, "Plank", we.Exercise)
		assert.Equal(t, 2, we.Order)
		assert.Len(t, we.Sets, 0)
	})
}

func TestWorkoutIn_Validation(t *testing.T) {
	t.Run("Valid WorkoutIn", func(t *testing.T) {
		id := ksuid.New()
		start := time.Now()
		end := start.Add(time.Hour)

		workoutIn := WorkoutIn{
			ID:    id.String(),
			Name:  "Test Workout",
			Start: start,
			End:   end,
			Exercises: []WorkoutExerciseIn{
				{
					Exercise: "Push Up",
					Order:    1,
					Sets: []SetIn{
						{
							Completed: true,
							Weight:    80.0,
							Reps:      10,
						},
					},
				},
			},
		}

		assert.Equal(t, id.String(), workoutIn.ID)
		assert.Equal(t, "Test Workout", workoutIn.Name)
		assert.Equal(t, start, workoutIn.Start)
		assert.Equal(t, end, workoutIn.End)
		assert.Len(t, workoutIn.Exercises, 1)
	})
}

func TestSetStructFields(t *testing.T) {
	t.Run("Set with all fields", func(t *testing.T) {
		id := ksuid.New()
		set := Set{
			ModifiableModel: ModifiableModel{
				Model: Model{ID: id},
			},
			WorkoutExerciseID: "workout_exercise_123",
			Completed:         true,
			Weight:            100.5,
			Reps:              12,
			Duration:          45.0,
			Distance:          2.5,
		}

		assert.Equal(t, id, set.ID)
		assert.Equal(t, "workout_exercise_123", set.WorkoutExerciseID)
		assert.True(t, set.Completed)
		assert.Equal(t, 100.5, set.Weight)
		assert.Equal(t, 12, set.Reps)
		assert.Equal(t, 45.0, set.Duration)
		assert.Equal(t, 2.5, set.Distance)
	})
}

func TestResponseStructs(t *testing.T) {
	t.Run("ExercisesResponse", func(t *testing.T) {
		exercises := []ExerciseOut{
			{Name: "Push Up", Category: "Body weight", Target: "Chest"},
			{Name: "Squat", Category: "Body weight", Target: "Legs"},
		}

		response := ExercisesResponse{
			Exercises: exercises,
		}

		assert.Len(t, response.Exercises, 2)
		assert.Equal(t, "Push Up", response.Exercises[0].Name)
		assert.Equal(t, "Squat", response.Exercises[1].Name)
	})

	t.Run("WorkoutResponse", func(t *testing.T) {
		workouts := []WorkoutOut{
			{ID: "1", Name: "Workout 1"},
			{ID: "2", Name: "Workout 2"},
		}

		response := WorkoutResponse{
			Workouts: workouts,
			Cursor:   "cursor123",
		}

		assert.Len(t, response.Workouts, 2)
		assert.Equal(t, "cursor123", response.Cursor)
	})
}
