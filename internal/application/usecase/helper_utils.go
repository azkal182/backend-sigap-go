package usecase

import (
	"time"

	"github.com/google/uuid"
)

func uuidPtrToString(id *uuid.UUID) *string {
	if id == nil {
		return nil
	}
	val := id.String()
	return &val
}

func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	val := t.Format(time.RFC3339)
	return &val
}

func timePtrToString(t *time.Time) *string {
	if t == nil {
		return nil
	}
	val := t.Format("2006-01-02")
	return &val
}
