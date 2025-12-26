package utils

func Map[T any, R any](items []T, f func(T) R) []R {
    result := make([]R, len(items))
    for i, item := range items {
        result[i] = f(item)
    }
    return result
}

func Map1Par[T any, R any, S any](items []T, par S, f func(T, S) R) []R {
    result := make([]R, len(items))
    for i, item := range items {
        result[i] = f(item, par)
    }
    return result
}
