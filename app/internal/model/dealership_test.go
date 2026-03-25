package model_test

import (
	"scenario-a/internal/model"
	"scenario-a/pkg/common"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDealership_TableName(t *testing.T) {
	d := &model.Dealership{}
	assert.Equal(t, "dealerships", d.TableName())
}

func TestDealership_GetTimezone(t *testing.T) {
	t.Run("valid timezone", func(t *testing.T) {
		d := &model.Dealership{Timezone: "Asia/Singapore"}
		loc, err := d.GetTimezone()
		require.NoError(t, err)
		assert.NotNil(t, loc)
		assert.Equal(t, "Asia/Singapore", loc.String())
	})

	t.Run("invalid timezone", func(t *testing.T) {
		d := &model.Dealership{Timezone: "Invalid/Timezone"}
		loc, err := d.GetTimezone()
		require.Error(t, err)
		assert.Nil(t, loc)
		assert.Contains(t, err.Error(), "cannot load dealership timezone")
	})
}

func TestDealership_GetOpenCloseTime(t *testing.T) {
	t.Run("valid times", func(t *testing.T) {
		d := &model.Dealership{OpenTime: "08:00:00", CloseTime: "17:00:00"}
		openTime, closeTime, err := d.GetOpenCloseTime()
		require.NoError(t, err)
		assert.Equal(t, 8, openTime.Hour())
		assert.Equal(t, 17, closeTime.Hour())
	})

	t.Run("invalid open time format", func(t *testing.T) {
		d := &model.Dealership{OpenTime: "8am", CloseTime: "17:00:00"}
		_, _, err := d.GetOpenCloseTime()
		require.Error(t, err)
		safeErr, ok := err.(*common.SafeError)
		require.True(t, ok)
		assert.Equal(t, "DEALERSHIP_INVALID_OPENING_TIME", safeErr.Short)
	})

	t.Run("invalid close time format", func(t *testing.T) {
		d := &model.Dealership{OpenTime: "08:00:00", CloseTime: "5pm"}
		_, _, err := d.GetOpenCloseTime()
		require.Error(t, err)
		safeErr, ok := err.(*common.SafeError)
		require.True(t, ok)
		assert.Equal(t, "DEALERSHIP_INVALID_CLOSING_TIME", safeErr.Short)
	})

	t.Run("close time before open time", func(t *testing.T) {
		d := &model.Dealership{OpenTime: "17:00:00", CloseTime: "08:00:00"}
		_, _, err := d.GetOpenCloseTime()
		require.Error(t, err)
		safeErr, ok := err.(*common.SafeError)
		require.True(t, ok)
		assert.Equal(t, "DEALERSHIP_INVALID_TIME_RANGE", safeErr.Short)
	})
}

func TestDealership_IsOpen(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Singapore")

	// Set up a base dealership for testing time scenarios
	getDealership := func() *model.Dealership {
		return &model.Dealership{
			Timezone:      "Asia/Singapore",
			OpenTime:      "08:00:00",
			CloseTime:     "17:00:00",
			IsWeekendOpen: false,
		}
	}

	t.Run("valid operating hours on weekday", func(t *testing.T) {
		d := getDealership()
		// A weekday (e.g. Wednesday 2026-03-25)
		reqStart := time.Date(2026, 3, 25, 9, 0, 0, 0, loc)
		reqEnd := time.Date(2026, 3, 25, 10, 0, 0, 0, loc)
		err := d.IsOpen(reqStart, reqEnd)
		assert.NoError(t, err)
	})

	t.Run("end before start", func(t *testing.T) {
		d := getDealership()
		reqStart := time.Date(2026, 3, 25, 10, 0, 0, 0, loc)
		reqEnd := time.Date(2026, 3, 25, 9, 0, 0, 0, loc)
		err := d.IsOpen(reqStart, reqEnd)
		require.Error(t, err)
		safeErr, ok := err.(*common.SafeError)
		require.True(t, ok)
		assert.Equal(t, "INVALID_BOOKING_TIME_RANGE", safeErr.Short)
	})

	t.Run("invalid timezone", func(t *testing.T) {
		d := getDealership()
		d.Timezone = "Invalid/Timezone"
		reqStart := time.Date(2026, 3, 25, 9, 0, 0, 0, loc)
		reqEnd := time.Date(2026, 3, 25, 10, 0, 0, 0, loc)
		err := d.IsOpen(reqStart, reqEnd)
		require.Error(t, err)
		safeErr, ok := err.(*common.SafeError)
		require.True(t, ok)
		assert.Equal(t, "DEALERSHIP_INVALID_TIMEZONE", safeErr.Short)
	})

	t.Run("weekend closed", func(t *testing.T) {
		d := getDealership()
		// Saturday
		reqStart := time.Date(2026, 3, 28, 9, 0, 0, 0, loc)
		reqEnd := time.Date(2026, 3, 28, 10, 0, 0, 0, loc)
		err := d.IsOpen(reqStart, reqEnd)
		require.Error(t, err)
		safeErr, ok := err.(*common.SafeError)
		require.True(t, ok)
		assert.Equal(t, "INVALID_BOOKING_WEEKEND", safeErr.Short)
	})

	t.Run("weekend opened", func(t *testing.T) {
		d := getDealership()
		d.IsWeekendOpen = true
		// Saturday
		reqStart := time.Date(2026, 3, 28, 9, 0, 0, 0, loc)
		reqEnd := time.Date(2026, 3, 28, 10, 0, 0, 0, loc)
		err := d.IsOpen(reqStart, reqEnd)
		assert.NoError(t, err)
	})

	t.Run("outside operating hours - too early", func(t *testing.T) {
		d := getDealership()
		reqStart := time.Date(2026, 3, 25, 7, 0, 0, 0, loc)
		reqEnd := time.Date(2026, 3, 25, 8, 30, 0, 0, loc)
		err := d.IsOpen(reqStart, reqEnd)
		require.Error(t, err)
		safeErr, ok := err.(*common.SafeError)
		require.True(t, ok)
		assert.Equal(t, "INVALID_BOOKING_TIME_RANGE", safeErr.Short)
	})

	t.Run("outside operating hours - too late", func(t *testing.T) {
		d := getDealership()
		reqStart := time.Date(2026, 3, 25, 16, 0, 0, 0, loc)
		reqEnd := time.Date(2026, 3, 25, 18, 0, 0, 0, loc)
		err := d.IsOpen(reqStart, reqEnd)
		require.Error(t, err)
		safeErr, ok := err.(*common.SafeError)
		require.True(t, ok)
		assert.Equal(t, "INVALID_BOOKING_TIME_RANGE", safeErr.Short)
	})

	t.Run("different timezone input correctly converts", func(t *testing.T) {
		d := getDealership()
		// Singapore is UTC+8. Open at 08:00 SGT is 00:00 UTC.
		// So reqStart 01:00 UTC is 09:00 SGT (Valid)
		// reqEnd 02:00 UTC is 10:00 SGT (Valid)
		utcLoc, _ := time.LoadLocation("UTC")
		reqStart := time.Date(2026, 3, 25, 1, 0, 0, 0, utcLoc)
		reqEnd := time.Date(2026, 3, 25, 2, 0, 0, 0, utcLoc)
		err := d.IsOpen(reqStart, reqEnd)
		assert.NoError(t, err)
	})
	
	t.Run("invalid internal opening time", func(t *testing.T) {
		d := getDealership()
		d.OpenTime = "invalid"
		reqStart := time.Date(2026, 3, 25, 9, 0, 0, 0, loc)
		reqEnd := time.Date(2026, 3, 25, 10, 0, 0, 0, loc)
		err := d.IsOpen(reqStart, reqEnd)
		require.Error(t, err)
		safeErr, ok := err.(*common.SafeError)
		require.True(t, ok)
		assert.Equal(t, "DEALERSHIP_INVALID_OPENING_TIME", safeErr.Short)
	})
}
