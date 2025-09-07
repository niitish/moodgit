package cmd

import (
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "moodgit",
	Short: "log and track your mood via cli",
	Long: `moodgit - a simple CLI tool to log and track your mood.

moodgit helps you maintain a personal mood journal through the command line.
track your emotional state with intensity levels, descriptive messages,
and tags to build insights into your mood patterns over time.

example workflow:
  moodgit init
  moodgit add -i 8 -o happy -m "great day at work!" -t work
  moodgit log`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println()
		color.C256(201).Println("█▀▄▀█ ████▄ ████▄ ██▄     ▄▀  ▄█    ▄▄▄▄▀")
		color.C256(165).Println("█ █ █ █   █ █   █ █  █  ▄▀    ██ ▀▀▀ █   ")
		color.C256(129).Println("█ ▄ █ █   █ █   █ █   █ █ ▀▄  ██     █   ")
		color.C256(93).Println("█   █ ▀████ ▀████ █  █  █   █ ▐█    █    ")
		color.C256(57).Println("   █              ███▀   ███   ▐    ▀    ")
		color.C256(21).Println("  ▀                                      ")
		fmt.Println()
		fmt.Println("a simple CLI tool to log and track your mood.")
		fmt.Println("run 'moodgit --help' to see available commands.")

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
