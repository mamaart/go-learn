package functools

import "log"

func MustV[T any](v T, err error) T {
	if err != nil {
		log.Fatal(err)
	}
	return v
}

func Must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
