package main

import (
	"github.com/rl404/naka/internal/utils"
	"github.com/spf13/cobra"
)

func main() {
	cmd := cobra.Command{
		Use:   "naka",
		Short: "Naka Discord Bot",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "bot",
		Short: "Run bot",
		RunE: func(*cobra.Command, []string) error {
			return bot()
		},
	})

	if err := cmd.Execute(); err != nil {
		utils.Fatal(err.Error())
	}
}
