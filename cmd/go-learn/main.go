package main

import (
	"fmt"
	"log"

	"github.com/mamaart/go-learn/pkg/d2l"
	"github.com/mamaart/go-learn/pkg/inside"
	"github.com/spf13/cobra"
)

func main() {
	var l LearnCmd
	cmd := cobra.Command{
		Use:   "golearn",
		Short: "this is a tool to automize tasks on DTU inside and learn",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	flags := cmd.PersistentFlags()

	flags.StringVarP(&l.username, "username", "u", "", "your username")
	flags.StringVarP(&l.password, "password", "p", "", "your password")

	cmd.CompletionOptions.HiddenDefaultCmd = true
	cmd.AddCommand(l.d2lCmd())
	cmd.AddCommand(l.insideCmd())
	cmd.Execute()
}

type LearnCmd struct {
	username, password string
}

func (l *LearnCmd) d2lCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "d2l",
		Aliases:   []string{"l", "new", "learn"},
		Short:     "call actions on the d2l api",
		ValidArgs: []string{"whoami"},
		Args:      cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if l.username == "" || l.password == "" {
				log.Fatal("Failed: username and password not provided")
			}
			switch args[0] {
			case "whoami":
				l, err := d2l.New(l.username, l.password)
				if err != nil {
					log.Fatal(err)
				}

				r, err := l.Whoami()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Println(string(r))
			}
		},
	}
	return cmd
}

func (l *LearnCmd) insideCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "inside",
		Aliases:   []string{"i", "old"},
		Short:     "call actions on the dtu inside api",
		ValidArgs: []string{"grades"},
		Args:      cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if l.username == "" || l.password == "" {
				log.Fatal("Failed: username and password not provided")
			}
			switch args[0] {
			case "grades":
				i, err := inside.New(l.username, l.password)
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
		},
	}
	return cmd
}
