package utils

func Ptr[T any](v T) *T {
	return &v
}

// RemoveFromArray removes element from array by index
func RemoveFromArray[T any](arr []T, index int) []T {
	if index < 0 || index >= len(arr) {
		// Index out of bounds
		return arr
	}

	return append(arr[:index], arr[index+1:]...)
}
