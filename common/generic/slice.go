package generic

func SliceConvert[S any, T any](source []S, convert func(S) T) (target []T) {
	for _, element := range source {
		target = append(target, convert(element))
	}
	return
}
