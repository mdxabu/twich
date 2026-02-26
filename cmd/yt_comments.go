package cmd

import (
	"fmt"
	"os"

	"github.com/mdxabu/twich/internals"
	"github.com/spf13/cobra"
)

// ytCommentsCmd represents the yt-comments command
var ytCommentsCmd = &cobra.Command{
	Use:   "yt-comments [video_url]",
	Short: "List YouTube live chat comments for a given video URL",
	Long: `Fetch and display YouTube live chat comments for a specified video.
Example usage:
  twich yt-comments https://www.youtube.com/watch?v=EB916GNkAsQ
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		url := args[0]
		err := internals.ListYTComments(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error fetching YouTube comments: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(ytCommentsCmd)
}
