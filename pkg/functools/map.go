package functools

func Map[T, R any](vs []T, fn func(v T) R) []R {
	out := make([]R, len(vs))
	for i, e := range vs {
		out[i] = fn(e)
	}
	return out
}

func Foldr[T, R any](vs []T, initial R, fn func(T, R) R) R {
	if len(vs) == 0 {
		return initial
	}
	return fn(vs[0], Foldr(vs[1:], initial, fn))
}

func Foldl[T, R any](vs []T, initial R, fn func(T, R) R) R {
	if len(vs) == 0 {
		return initial
	}
	return Foldl(vs[1:], fn(vs[0], initial), fn)
}

func Sum[T Number](vs []T) T {
	return Foldr(vs, T(0), func(a T, b T) T {
		return a + b
	})
}

type Float interface {
	~float32 | ~float64
}

type Integer interface {
	Signed | Unsigned
}

type Number interface {
	Integer | Float
}

type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}
