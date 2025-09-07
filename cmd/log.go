package cmd

import (
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
  moodgit log -l 20            # show last 20 entries`,
	Run: func(cmd *cobra.Command, args []string) {
		limit, _ := cmd.Flags().GetUint16("limit")
		internal.GetHistory(limit)
	},
}

func init() {
	rootCmd.AddCommand(logCmd)

	logCmd.Flags().Uint16P("limit", "l", 10, "number of entries to show")
}
