package functools

import (
	"fmt"
)

func Retry[T any](retries int, action func() (T, bool), err error) (v T, _ error) {
	if retries == 0 {
		return v, err
	}

	if v, ok := action(); ok {
		return v, nil
	}
	fmt.Println("retrying..")
	return Retry(retries-1, action, err)
}
