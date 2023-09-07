package d2l

import (
	"fmt"
	"log"

	"github.com/mamaart/go-learn/internal/auth"
	"github.com/mamaart/go-learn/pkg/d2l"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "d2l",
		Aliases:   []string{"l", "new", "learn"},
		Short:     "call actions on the d2l api",
		ValidArgs: []string{"whoami"},
		Args:      cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(login())
	cmd.AddCommand(wmoami())
	return cmd
}

func login() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "login",
		Args: cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]
			password := args[1]
			l, err := d2l.New(d2l.Options{
				Credentials: &auth.Credentials{
					Username: username,
					Password: password,
				},
			})
			if err != nil {
				log.Fatal(err)
			}
			r, err := l.Whoami()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(r))
		},
	}
	return cmd
}

func wmoami() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "shows the basic info of your user",
		Run: func(cmd *cobra.Command, args []string) {
			l, err := d2l.New(d2l.DefaultOptions())
			if err != nil {
				log.Fatal(err)
			}

			r, err := l.Whoami()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(r))
		},
	}
	return cmd
}
