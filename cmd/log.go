package cmd

import (
	"fmt"
	"moodgit/internal"

	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "show mood logs",
	Long: `display your mood history in chronological order.

this command shows your recent mood entries with full details including:
- mood type and intensity level (0-10)
- timestamp when the entry was created
- optional message describing your state
- any tags associated with the entry

the entries are displayed with color-coding based on mood type and
intensity for better visual representation.

examples:
  moodgit log                  # show last 10 entries
  moodgit log -l 20            # show last 20 entries
  moodgit log -i               # show interactive log with 10 entries per page
  moodgit log -i -l 25         # show interactive log with 25 entries per page`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := internal.InitDB(); err != nil {
			return fmt.Errorf("%w", err)
		}

		limit, _ := cmd.Flags().GetUint16("limit")
		interactive, _ := cmd.Flags().GetBool("interactive")

		if interactive {
			if err := internal.StartInteractiveLog(int(limit)); err != nil {
				return fmt.Errorf("error starting interactive log: %w", err)
			}
		} else {
			return internal.GetHistory(limit)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(logCmd)

	logCmd.Flags().Uint16P("limit", "l", 10, "number of entries to show (page size for interactive mode)")
	logCmd.Flags().BoolP("interactive", "i", false, "show interactive log with pagination")
}
