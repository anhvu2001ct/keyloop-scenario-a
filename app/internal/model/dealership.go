package model

import (
	"fmt"
	"scenario-a/pkg/common"
	"time"
)

type Dealership struct {
	CommonModel

	Name          string
	OpenTime      string
	CloseTime     string
	Timezone      string // e.g Asia/Singapore
	IsWeekendOpen bool
}

func (*Dealership) TableName() string {
	return "dealerships"
}

func (d *Dealership) GetTimezone() (*time.Location, error) {
	loc, err := time.LoadLocation(d.Timezone)
	if err != nil {
		err = fmt.Errorf("cannot load dealership timezone (%s): %w", d.Timezone, err)
	}
	return loc, err
}

func (d *Dealership) GetOpenCloseTime() (time.Time, time.Time, error) {
	openTime, err := time.Parse(time.TimeOnly, d.OpenTime)
	if err != nil {
		return time.Time{}, time.Time{}, &common.SafeError{
			Short:   "DEALERSHIP_INVALID_OPENING_TIME",
			UserMsg: "Invalid opening time format",
			Cause:   err,
			Metadata: map[string]any{
				"dealership_id": d.ID,
				"opening_time":  d.OpenTime,
			},
		}
	}

	closeTime, err := time.Parse(time.TimeOnly, d.CloseTime)
	if err != nil {
		return time.Time{}, time.Time{}, &common.SafeError{
			Short:   "DEALERSHIP_INVALID_CLOSING_TIME",
			UserMsg: "Invalid closing time format",
			Cause:   err,
			Metadata: map[string]any{
				"dealership_id": d.ID,
				"closing_time":  d.CloseTime,
			},
		}
	}

	if closeTime.Before(openTime) {
		return time.Time{}, time.Time{}, &common.SafeError{
			Short:   "DEALERSHIP_INVALID_TIME_RANGE",
			UserMsg: "Invalid dealership time, closeTime cannot less than openTime",
			Metadata: map[string]any{
				"dealership_id": d.ID,
				"open_time":     d.OpenTime,
				"close_time":    d.CloseTime,
			},
		}
	}

	return openTime, closeTime, nil
}

func (d *Dealership) IsOpen(reqStart, reqEnd time.Time) error {
	if reqEnd.Before(reqStart) {
		return &common.SafeError{
			Short:   "INVALID_BOOKING_TIME_RANGE",
			UserMsg: "The request end time must be after start time",
			Metadata: map[string]any{
				"req_start_at": reqStart.Format(time.RFC3339),
				"req_end_at":   reqEnd.Format(time.RFC3339),
			},
		}
	}

	loc, err := d.GetTimezone()
	if err != nil {
		return &common.SafeError{
			Short:   "DEALERSHIP_INVALID_TIMEZONE",
			UserMsg: "Invalid dealership timezone",
			Cause:   err,
			Metadata: map[string]any{
				"dealership_id": d.ID,
				"timezone":      d.Timezone,
			},
		}
	}

	// Convert to dealership's local time before any comparison
	reqStart, reqEnd = reqStart.In(loc), reqEnd.In(loc)

	isWeekend := func(t time.Time) bool {
		return t.Weekday() == time.Saturday || t.Weekday() == time.Sunday
	}
	if !d.IsWeekendOpen && (isWeekend(reqStart) || isWeekend(reqEnd)) {
		return &common.SafeError{
			Short:   "INVALID_BOOKING_WEEKEND",
			UserMsg: "The request time is in weekend, but dealership is closed on weekends",
			Metadata: map[string]any{
				"dealership_id": d.ID,
				"req_start_at":  reqStart.Format(time.RFC3339),
				"req_end_at":    reqEnd.Format(time.RFC3339),
			},
		}
	}

	openTime, closeTime, err := d.GetOpenCloseTime()
	if err != nil {
		return err
	}

	// Normalize to a dummy date for time-of-day comparison only.
	// time.Parse returns UTC, so we use UTC consistently across all four values.
	normalize := func(t time.Time) time.Time {
		h, m, _ := t.Clock()
		return time.Date(0, 1, 1, h, m, 0, 0, time.UTC)
	}

	if normalize(reqStart).Before(normalize(openTime)) || normalize(reqEnd).After(normalize(closeTime)) {
		return &common.SafeError{
			Short:   "INVALID_BOOKING_TIME_RANGE",
			UserMsg: "The request time is outside operating hours",
			Metadata: map[string]any{
				"dealership_id": d.ID,
				"req_start_at":  reqStart.Format(time.RFC3339),
				"req_end_at":    reqEnd.Format(time.RFC3339),
			},
		}
	}

	return nil
}
