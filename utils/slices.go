package utils

func Keys[K string | int, V any, M map[K]V](m M) []K {
	var keys []K
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
