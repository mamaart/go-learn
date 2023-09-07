package main

import (
	"fmt"

	"github.com/mamaart/go-learn/pkg/functools"
	"github.com/mamaart/go-learn/pkg/inside"
)

func main() {
	i := functools.MustV(inside.New(inside.DefaultOptions()))
	grades(i)
}

func grades(i *inside.Inside) {
	fmt.Println(functools.MustV(
		i.Grades(),
	))
}

func whoami(i *inside.Inside) {
	fmt.Println(functools.MustV(
		i.Whoami(),
	))
}
