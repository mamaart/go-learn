package main

import (
	"fmt"

	"github.com/mamaart/go-learn/pkg/d2l"
	"github.com/mamaart/go-learn/pkg/functools"
)

func main() {
	l := functools.MustV(d2l.New(d2l.DefaultOptions()))
	whoami(l)
}

func whoami(l *d2l.D2L) {
	fmt.Println(string(functools.MustV(l.Whoami())))
}
