package main

import (
	"github.com/mamaart/go-learn/cmd/golearn/d2l"
	"github.com/mamaart/go-learn/cmd/golearn/inside"
	"github.com/spf13/cobra"
)

func main() {
	cmd := cobra.Command{
		Use:   "golearn",
		Short: "automize tasks on DTU inside and learn",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.CompletionOptions.HiddenDefaultCmd = true
	cmd.AddCommand(d2l.Cmd())
	cmd.AddCommand(inside.Cmd())
	cmd.Execute()
}
