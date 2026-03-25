package model_test

import (
	"scenario-a/internal/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBaseModel_GetID(t *testing.T) {
	b := &model.BaseModel{ID: 10}
	assert.Equal(t, int64(10), b.GetID())
}

func TestCommonModel_GetUUID(t *testing.T) {
	c := &model.CommonModel{UUID: "test-uuid"}
	assert.Equal(t, "test-uuid", c.GetUUID())
}

func TestWithSoftDelete_GetDeletedAt(t *testing.T) {
	now := time.Now()
	w := &model.WithSoftDelete{DeletedAt: &now}
	assert.Equal(t, &now, w.GetDeletedAt())
}
