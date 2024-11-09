package np

func OrDef[T comparable](value T, defaultValue T) T {
	var zero T
	if value != zero {
		return value
	}
	return defaultValue
}
