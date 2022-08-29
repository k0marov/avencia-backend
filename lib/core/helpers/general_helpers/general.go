package general_helpers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// DecodeFloat is needed because if an integer value is stored in firestore it cannot be decoded as float64 easily
func DecodeFloat(value any) (float64, error) {
	asFloat, ok := value.(float64)
	if ok {
		return asFloat, nil
	}
	asInt, ok := value.(int64)
	if ok {
		return float64(asInt), nil
	}
	return 0, fmt.Errorf("failed to decode %+v as float64", value)
}

func DecodeTime(value any) (time.Time, error) {
	asInt, ok := value.(int64) 
	if !ok {
    return time.Time{}, fmt.Errorf("failed to convert %v to int64", value)
	}
	return time.Unix(asInt, 0), nil
}


// FindInSlice returns -1 if element does not exist
func FindInSlice[T any](slice []T, match func(T) bool) (index int) {
	for i := range slice {
		if match(slice[i]) {
			return i
		}
	}
	return -1
}

func RandomId() string {
	uuid, _ := uuid.NewUUID() 
	return uuid.String()
}
