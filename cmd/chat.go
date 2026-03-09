/*
Copyright © 2026 @mdxabu
*/
package cmd

import (
	"log"

	chat "github.com/mdxabu/twich/internals"
	"github.com/spf13/cobra"
)

var (
	twitchUser string
	ytChannel  string
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Fetch live chat from Twitch and YouTube simultaneously",
	Long: `This command connects to a Twitch channel's IRC chat and
scrapes a YouTube Live stream chat using a Python bridge.

The -y flag accepts a YouTube @handle, channel username, full channel URL,
or a raw video ID. The tool will automatically resolve the channel's
currently active live stream.

Examples:
  twich chat -t helpytv -y @mkbhd
  twich chat -y mkbhd
  twich chat -y https://www.youtube.com/@mkbhd
  twich chat -y JFfPyuo67E8
  twich chat -t helpytv`,

	Run: func(cmd *cobra.Command, args []string) {
		if twitchUser == "" && ytChannel == "" {
			log.Fatal("Please provide at least a Twitch username (-t) or a YouTube channel (-y)")
		}

		// Calling the combined function from the internals package.
		// ytChannel can be a @handle, username, URL, or raw video ID —
		// the Python proxy resolves it to the active live stream automatically.
		chat.FetchCombinedChat(twitchUser, ytChannel)
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)

	// Defining separate flags for Twitch and YouTube
	chatCmd.Flags().StringVarP(&twitchUser, "twitch", "t", "", "Twitch username to fetch chat from")
	chatCmd.Flags().StringVarP(&ytChannel, "youtube", "y", "", "YouTube channel @handle, username, URL, or video ID")
}
