package slice

func Convert[T, K any](slice []T, fn func(T) (K, error)) ([]K, error) {
	newSlice := make([]K, 0, len(slice))
	for _, el := range slice {
		newEl, err := fn(el)
		if err != nil {
			return nil, err
		}
		newSlice = append(newSlice, newEl)
	}
	return newSlice, nil
}

func MustConvert[T, K any](slice []T, fn func(T) K) []K {
	newSlice := make([]K, 0, len(slice))
	for _, el := range slice {
		newSlice = append(newSlice, fn(el))
	}
	return newSlice
}
