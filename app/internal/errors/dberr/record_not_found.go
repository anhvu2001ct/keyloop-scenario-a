package dberr

import (
	"fmt"
	"strings"

	"github.com/samber/lo"
)

type RecordNotFound struct {
	TableName  string
	Conditions map[string]any
}

func (e *RecordNotFound) Error() string {
	conditions := lo.MapToSlice(e.Conditions, func(k string, v any) string {
		if str, ok := v.(string); ok {
			return fmt.Sprintf("%s=%q", k, str)
		}
		return fmt.Sprintf("%s=%v", k, v)
	})
	conditionStr := strings.Join(conditions, ", ")
	return fmt.Sprintf("record not found for %s(%s)", e.TableName, conditionStr)
}
