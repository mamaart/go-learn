# 🎓 Go Learn

![GitHub](https://img.shields.io/github/license/mamaart/go-learn)

This is a simple webscraper to access the api of the intranet at DTU. The intranet is using the [Desire2Learn](https://docs.valence.desire2learn.com/reference.html) API so every possible endpoint is available at their website. 

> ⚠️ As a student you dont have permission to use all endpoints. 

You can use this api for automation of tasks related to your user, and if you get extra permissions from the administration you will be able to manage other users as well.

The webscraper can be used both as a library in another application or as a command line tool

There is another part of the program as well whic can be used to get info from the old dtu inside

### Examples

#### Example of getting grade list from inside

```go
package main

import (
	"log"

	"github.com/mamaart/go-learn/pkg/inside"
)

func login(username, password string) {
	i, err := inside.New(inside.Options{
		Credentials: &auth.Credentials{
			Username: username,
			Password: password,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("logged in successfully!")

	grades, err := i.GetGrades()
	if err != nil {
		log.Fatal(err)
	}

	for _, e := range grades {
		log.Println(e)
	}
}
```

#### Example of calling whoami on d2l

```go
package main

import (
	"log"

	"github.com/mamaart/go-learn/pkg/d2l"
)

func login(username, password string) {
	i, err := d2l.New(inside.Options{
		Credentials: &auth.Credentials{
			Username: username,
			Password: password,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("logged in successfully!")

	me, err := i.Whoami()
	if err != nil {
		log.Fatal(err)
	}
    log.Println(me)
}
```

## 🛠️ Todo

- [ ] Handle partial authentication problem when one or more cookies are expired. 
- [ ] Finish implementing the endpoints from python
- [ ] Make tests on endpoints
- [ ] Automize implementaion of endpoints from the reference
- [ ] Parse and return types from the functions.
- [x] Make a TUI for simple use cases
- [ ] Make tests of all endpoints and register which are not available for students


