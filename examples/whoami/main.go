package main

import (
	"fmt"
	"log"
	"os"

	golearn "github.com/mamaart/go-learn"
)

func main() {
	if len(os.Args) > 3 {
		log.Fatal("username and password not provided")
	}

	username := os.Args[1]
	password := os.Args[2]

	l, err := golearn.New(username, password)
	if err != nil {
		log.Fatal(err)
	}

	r, err := l.Whoami()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(r))
}
