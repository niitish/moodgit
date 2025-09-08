package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gookit/color"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize a moodgit repository in your home directory",
	Long: `initialize a new moodgit repository in your home directory.

this command sets up the necessary directory structure and database file
for moodgit to store your mood entries. it creates:
- a .moodgit directory in your home folder
- a SQLite database file (moodgit.db) to store your mood entries
- required database schema for tracking moods

the repository will be created at ~/.moodgit/ and is persistent across
all your moodgit sessions. you only need to run this command once.

examples:
  moodgit init                 # initialize a new repository
  moodgit init --force         # reinitialize and reset all data`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// show ascii art
		fmt.Println()
		color.C256(201).Println("█▀▄▀█ ████▄ ████▄ ██▄     ▄▀  ▄█    ▄▄▄▄▀")
		color.C256(165).Println("█ █ █ █   █ █   █ █  █  ▄▀    ██ ▀▀▀ █   ")
		color.C256(129).Println("█ ▄ █ █   █ █   █ █   █ █ ▀▄  ██     █   ")
		color.C256(93).Println("█   █ ▀████ ▀████ █  █  █   █ ▐█    █    ")
		color.C256(57).Println("   █              ███▀   ███   ▐    ▀    ")
		color.C256(21).Println("  ▀                                      ")
		fmt.Println()

		force, _ := cmd.Flags().GetBool("force")

		// get home dir, and create a .moodgit directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}

		repoPath := filepath.Join(homeDir, ".moodgit")
		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			if err := os.Mkdir(repoPath, 0755); err != nil {
				return fmt.Errorf("failed to create .moodgit directory: %w", err)
			}
		}

		// create moodgit.db file
		dbPath := filepath.Join(repoPath, "moodgit.db")
		if _, err := os.Stat(dbPath); os.IsNotExist(err) || force {
			file, err := os.Create(dbPath)
			if err != nil {
				return fmt.Errorf("failed to create moodgit.db file: %w", err)
			}
			defer file.Close()
		} else if !force {
			return fmt.Errorf("moodgit.db file already exists, use --force to overwrite")
		}

		fmt.Printf("initialized empty moodgit repository in %s\n", repoPath)
		fmt.Printf("you can now start using moodgit!\n")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolP("force", "f", false, "force initialization of moodgit repository")

	initCmd.Aliases = []string{"initialize", "create"}
}
