package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mamaart/go-learn/pkg/d2l"
)

func main() {
	if len(os.Args) > 3 {
		log.Fatal("username and password not provided")
	}

	username := os.Args[1]
	password := os.Args[2]

	l, err := d2l.New(username, password)
	if err != nil {
		log.Fatal(err)
	}

	r, err := l.Whoami()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(r))
}
