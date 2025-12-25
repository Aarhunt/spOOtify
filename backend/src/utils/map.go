package utils

func Map[T any, R any](items []T, f func(T) R) []R {
    result := make([]R, len(items))
    for i, item := range items {
        result[i] = f(item)
    }
    return result
}
