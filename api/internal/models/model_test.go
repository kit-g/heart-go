package models

import (
	"github.com/segmentio/ksuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

// Test structures to use with our models
type TestModel struct {
	Model
	Name string
}

type TestModifiableModel struct {
	ModifiableModel
	Name string
}

type TestSoftDeleteModel struct {
	SoftDeleteModel
	Name string
}

func setupTestDBForModel(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&TestModel{}, &TestModifiableModel{}, &TestSoftDeleteModel{})
	require.NoError(t, err)

	return db
}

func TestModelF(t *testing.T) {
	db := setupTestDBForModel(t)

	t.Run("BeforeCreate with nil UUID", func(t *testing.T) {
		model := &TestModel{
			Name: "Test Model",
		}

		result := db.Create(model)
		assert.NoError(t, result.Error)
		assert.NotEqual(t, ksuid.Nil, model.ID)
		assert.False(t, model.CreatedAt.IsZero())
	})

	t.Run("BeforeCreate with preset UUID", func(t *testing.T) {
		presetUUID := ksuid.New()
		model := &TestModel{
			Model: Model{
				ID: presetUUID,
			},
			Name: "Test Model with Preset UUID",
		}

		result := db.Create(model)
		assert.NoError(t, result.Error)
		assert.Equal(t, presetUUID, model.ID)
		assert.False(t, model.CreatedAt.IsZero())
	})

	t.Run("Retrieve model", func(t *testing.T) {
		model := &TestModel{
			Name: "Test Model for Retrieval",
		}

		// Create
		result := db.Create(model)
		assert.NoError(t, result.Error)

		// Retrieve
		var retrieved TestModel
		result = db.First(&retrieved, "id = ?", model.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, model.ID, retrieved.ID)
		assert.Equal(t, model.Name, retrieved.Name)
		assert.Equal(t, model.CreatedAt.UTC(), retrieved.CreatedAt.UTC())
	})
}

func TestModifiableModelF(t *testing.T) {
	db := setupTestDBForModel(t)

	t.Run("Create and update", func(t *testing.T) {
		model := &TestModifiableModel{
			Name: "Test Modifiable Model",
		}

		// Create
		result := db.Create(model)
		assert.NoError(t, result.Error)
		assert.NotEqual(t, ksuid.Nil, model.ID)
		assert.False(t, model.CreatedAt.IsZero())
		assert.False(t, model.UpdatedAt.IsZero())
		initialUpdatedAt := model.UpdatedAt

		// Wait a moment to ensure time difference
		time.Sleep(time.Millisecond)

		// Update
		model.Name = "Updated Name"
		result = db.Save(model)
		assert.NoError(t, result.Error)
		assert.NotEqual(t, initialUpdatedAt, model.UpdatedAt)
	})

	t.Run("Retrieve with timestamps", func(t *testing.T) {
		model := &TestModifiableModel{
			Name: "Test Modifiable Model for Retrieval",
		}

		// Create
		result := db.Create(model)
		assert.NoError(t, result.Error)

		// Retrieve
		var retrieved TestModifiableModel
		result = db.First(&retrieved, "id = ?", model.ID)
		assert.NoError(t, result.Error)
		assert.Equal(t, model.ID, retrieved.ID)
		assert.Equal(t, model.Name, retrieved.Name)
		assert.Equal(t, model.CreatedAt.UTC(), retrieved.CreatedAt.UTC())
		assert.Equal(t, model.UpdatedAt.UTC(), retrieved.UpdatedAt.UTC())
	})
}

func TestSoftDeleteModelF(t *testing.T) {
	db := setupTestDBForModel(t)

	t.Run("Soft delete and recovery", func(t *testing.T) {
		model := &TestSoftDeleteModel{
			Name: "Test Soft Delete Model",
		}

		// Create
		result := db.Create(model)
		assert.NoError(t, result.Error)
		assert.NotEqual(t, ksuid.Nil, model.ID)
		originalID := model.ID

		// Soft Delete
		result = db.Delete(model)
		assert.NoError(t, result.Error)

		// Try to retrieve normally (should fail)
		var retrieved TestSoftDeleteModel
		result = db.First(&retrieved, "id = ?", originalID)
		assert.Error(t, result.Error)
		assert.Equal(t, gorm.ErrRecordNotFound, result.Error)

		// Retrieve with unscoped (should succeed)
		result = db.Unscoped().First(&retrieved, "id = ?", originalID)
		assert.NoError(t, result.Error)
		assert.Equal(t, originalID, retrieved.ID)
		assert.NotNil(t, retrieved.DeletedAt.Time)

		// Restore
		result = db.Unscoped().Model(&retrieved).Update("deleted_at", nil)
		assert.NoError(t, result.Error)

		// Should be able to retrieve normally again
		var restored TestSoftDeleteModel
		result = db.First(&restored, "id = ?", originalID)
		assert.NoError(t, result.Error)
		assert.Equal(t, originalID, restored.ID)
	})

	t.Run("Permanent delete", func(t *testing.T) {
		model := &TestSoftDeleteModel{
			Name: "Test Permanent Delete Model",
		}

		// Create
		result := db.Create(model)
		assert.NoError(t, result.Error)
		originalID := model.ID

		// Permanent Delete
		result = db.Unscoped().Delete(model)
		assert.NoError(t, result.Error)

		// Verify permanent deletion
		var retrieved TestSoftDeleteModel
		result = db.Unscoped().First(&retrieved, "id = ?", originalID)
		assert.Error(t, result.Error)
		assert.Equal(t, gorm.ErrRecordNotFound, result.Error)
	})

	t.Run("Update timestamps", func(t *testing.T) {
		model := &TestSoftDeleteModel{
			Name: "Test Timestamps Model",
		}

		// Create
		result := db.Create(model)
		assert.NoError(t, result.Error)
		initialCreatedAt := model.CreatedAt
		initialUpdatedAt := model.UpdatedAt

		// Wait a moment to ensure time difference
		time.Sleep(time.Millisecond)

		// Update
		model.Name = "Updated Name"
		result = db.Save(model)
		assert.NoError(t, result.Error)
		assert.Equal(t, initialCreatedAt.UTC(), model.CreatedAt.UTC())
		assert.NotEqual(t, initialUpdatedAt.UTC(), model.UpdatedAt.UTC())
	})

	t.Run("Bulk soft delete", func(t *testing.T) {
		// Create multiple records
		models := []TestSoftDeleteModel{
			{Name: "Bulk Test 1"},
			{Name: "Bulk Test 2"},
			{Name: "Bulk Test 3"},
		}

		result := db.Create(&models)
		assert.NoError(t, result.Error)

		// Soft delete all
		result = db.Where("name LIKE ?", "Bulk Test%").Delete(&TestSoftDeleteModel{})
		assert.NoError(t, result.Error)

		// Verify all are soft deleted
		var count int64
		db.Model(&TestSoftDeleteModel{}).Where("name LIKE ?", "Bulk Test%").Count(&count)
		assert.Equal(t, int64(0), count)

		db.Unscoped().Model(&TestSoftDeleteModel{}).Where("name LIKE ?", "Bulk Test%").Count(&count)
		assert.Equal(t, int64(3), count)
	})
}

func TestModelInterfaces(t *testing.T) {
	t.Run("Model implements basic interface", func(t *testing.T) {
		var model interface{} = &Model{}
		_, ok := model.(interface {
			BeforeCreate(*gorm.DB) error
		})
		assert.True(t, ok, "Model should implement BeforeCreate interface")
	})

	t.Run("ModifiableModel embeds Model", func(t *testing.T) {
		model := ModifiableModel{}
		assert.NotNil(t, model.Model)
	})

	t.Run("SoftDeleteModel embeds ModifiableModel", func(t *testing.T) {
		model := SoftDeleteModel{}
		assert.NotNil(t, model.ModifiableModel)
	})
}
