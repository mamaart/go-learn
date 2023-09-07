package functools

func Fst[A, B any](a A, _ B) A {
	return a
}

func Snd[A, B any](_ A, b B) B {
	return b
}
