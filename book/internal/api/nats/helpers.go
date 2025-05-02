package nats

type num interface {
	~int | ~int32 | ~int64
}

func fromTo[I num, O num](is []I) []O {
	os := make([]O, len(is))
	for j := range is {
		os[j] = O(is[j])
	}

	return os
}
