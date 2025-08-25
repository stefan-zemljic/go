package jso

func filterMap[I, O any](input []I, predicate func(I) (O, bool)) []O {
	var output []O
	for _, item := range input {
		if mapped, ok := predicate(item); ok {
			output = append(output, mapped)
		}
	}
	return output
}
