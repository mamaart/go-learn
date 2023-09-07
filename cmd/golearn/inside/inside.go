package inside

import (
	"fmt"
	"log"

	"github.com/mamaart/go-learn/internal/auth"
	"github.com/mamaart/go-learn/pkg/inside"
	"github.com/spf13/cobra"
)

func Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:       "inside",
		Aliases:   []string{"i", "old"},
		Short:     "call actions on the dtu inside api",
		ValidArgs: []string{"grades"},
		Args:      cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	cmd.AddCommand(login())
	cmd.AddCommand(whoami())
	cmd.AddCommand(grades())
	cmd.AddCommand(gpa())
	return cmd
}

func login() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "login to inside",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]
			password := args[1]
			_, err := inside.New(inside.Options{
				Credentials: &auth.Credentials{
					Username: username,
					Password: password,
				},
			})
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("logged in successfully!")
		},
	}
	return cmd
}

func whoami() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "whoami",
		Short: "get the base info about your user",
		Run: func(cmd *cobra.Command, args []string) {
			i, err := inside.New(inside.DefaultOptions())
			if err != nil {
				log.Fatal(err)
			}
			resp, err := i.Whoami()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(resp))
		},
	}
	return cmd
}

func grades() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grades",
		Short: "get a list of grades",
		Run: func(cmd *cobra.Command, args []string) {
			i, err := inside.New(inside.DefaultOptions())
			if err != nil {
				log.Fatal(err)
			}
			grades, err := i.Grades()
			if err != nil {
				log.Fatal(err)
			}
			for _, e := range grades {
				fmt.Println(e)
			}
		},
	}
	return cmd
}

func gpa() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gpa",
		Short: "get the gpa",
		Run: func(cmd *cobra.Command, args []string) {
			i, err := inside.New(inside.DefaultOptions())
			if err != nil {
				log.Fatal(err)
			}
			gpa, err := i.GPA()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("GPA: %d\n", gpa)
		},
	}
	return cmd
}
