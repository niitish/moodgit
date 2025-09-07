package cmd

import (
	"fmt"
	"moodgit/internal"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add a mood entry",
	Long: `add a new mood entry to your moodgit repository.

this command allows you to record your current mood with various details:
- intensity: A scale from 0-10 indicating how strong the mood is
- mood: The type of mood (happy, sad, angry, anxious, excited, calm, stressed, tired, neutral)
- message: An optional description of your current state or what triggered the mood
- tags: Comma-separated tags to categorize or group related entries
- amend: Modify the last mood entry instead of creating a new one

examples:
  moodgit add -i 8 -o happy -m "got a promotion at work!" -t work,achievement
  moodgit add -i 3 -o sad -m "feeling down today"
  moodgit add -i 7 -o excited -t weekend,vacation
  moodgit add -a -i 9 -m "actually feeling even better!"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		intensity, _ := cmd.Flags().GetInt8("intensity")
		mood, _ := cmd.Flags().GetString("mood")
		message, _ := cmd.Flags().GetString("message")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		amend, _ := cmd.Flags().GetBool("amend")

		entry := internal.Entry{
			Intensity: intensity,
			Mood:      mood,
			Message:   message,
			Tags:      tags,
		}

		if amend {
			if err := internal.AmendLastEntry(entry); err != nil {
				return fmt.Errorf("failed to amend last mood entry: %w", err)
			}
		} else {
			if err := internal.AddEntry(entry); err != nil {
				return fmt.Errorf("failed to add mood entry: %w", err)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().BoolP("amend", "a", false, "amend the last mood entry")
	addCmd.Flags().Int8P("intensity", "i", 0, "mood intensity (0-10)")
	addCmd.Flags().StringP("mood", "o", "", "select your mood")
	addCmd.Flags().StringP("message", "m", "", "describe your mood")
	addCmd.Flags().StringSliceP("tags", "t", []string{}, "add tags to your entry (comma separated)")

	addCmd.MarkFlagRequired("intensity")
	addCmd.MarkFlagRequired("mood")
}
