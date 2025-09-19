package funki

/*
Standard map (a la map-reduce) but map has a different meaning in go.
*/
func Apply[F, T any](source []F, xform func(F) T) []T {
	result := make([]T, len(source))
	for i, value := range source {
		result[i] = xform(value)
	}
	return result
}
