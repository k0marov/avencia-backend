package general_helpers

import "fmt"

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

// FindInSlice returns -1 if element does not exist
func FindInSlice[T any](slice []T, match func(T) bool) (index int) {
	for i := range slice {
		if match(slice[i]) {
			return i
		}
	}
	return -1
}
