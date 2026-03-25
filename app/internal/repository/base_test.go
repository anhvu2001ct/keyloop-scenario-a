package repository

import (
	"scenario-a/internal/config"
	"scenario-a/internal/config/sqldb"
	"scenario-a/internal/model"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBaseRepo_Dealership_FindByID(t *testing.T) {
	config.MustInitForTest()
	cfg := config.Get()
	connectFunc := sqldb.MustInitForTest(cfg.Env)
	_, gormDB := connectFunc(t)
	repo := NewBaseRepo[*model.Dealership](gormDB)
	ctx := t.Context()

	id := int64(1)
	dealership, err := repo.FindByID(ctx, id)

	require.NoError(t, err)
	require.NotNil(t, dealership)
	require.Equal(t, id, dealership.ID)
}

func TestBaseRepo_Dealership_CreateOne_Fail(t *testing.T) {
	config.MustInitForTest()
	cfg := config.Get()
	connectFunc := sqldb.MustInitForTest(cfg.Env)
	_, gormDB := connectFunc(t)
	repo := NewBaseRepo[*model.Dealership](gormDB)
	ctx := t.Context()

	dealership := &model.Dealership{
		Name:      "test",
		OpenTime:  "08:00",
		CloseTime: "18:00",
	}
	dealership.ID = 10 // postgres prevent inserting with a specific id

	err := repo.Create(ctx, dealership)

	require.Error(t, err)
}

func TestBaseRepo_Dealership_CreateOne_Success(t *testing.T) {
	config.MustInitForTest()
	cfg := config.Get()
	connectFunc := sqldb.MustInitForTest(cfg.Env)
	_, gormDB := connectFunc(t)
	repo := NewBaseRepo[*model.Dealership](gormDB)
	ctx := t.Context()

	dealership := &model.Dealership{
		Name:      "test",
		OpenTime:  "08:00",
		CloseTime: "18:00",
	}
	dealership.UUID = "test-uuid"

	err := repo.Create(ctx, dealership)

	require.NoError(t, err)
	require.NotNil(t, dealership.ID)
}
