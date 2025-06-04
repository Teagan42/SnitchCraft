package utils

func Map[T any, R any](in []T, f func(T) R) []R {
	out := make([]R, len(in))
	for i, v := range in {
		out[i] = f(v)
	}
	return out
}

func Do[T any](in []T, f func(T)) []T {
	for _, v := range in {
		f(v)
	}
	return in
}
