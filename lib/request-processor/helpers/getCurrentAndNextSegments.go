package helpers

func GetCurrentAndNextSegments[T any](array []T) []map[string]T {
	var segments []map[string]T
	for i := 0; i < len(array)-1; i++ {
		segment := map[string]T{
			"currentSegment": array[i],
			"nextSegment":    array[i+1],
		}
		segments = append(segments, segment)
	}
	return segments
}
