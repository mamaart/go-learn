package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mamaart/go-learn/pkg/inside"
)

func main() {
	if len(os.Args) > 3 {
		log.Fatal("username and password not provided")
	}

	username := os.Args[1]
	password := os.Args[2]

	i, err := inside.New(username, password)
	if err != nil {
		log.Fatal(err)
	}

	grades, err := i.GetGrades()
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range grades {
		fmt.Println(e)
	}

}
