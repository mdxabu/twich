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
	ytVideoID  string
)

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Fetch live chat from Twitch and YouTube simultaneously",
	Long: `This command connects to a Twitch channel's IRC chat and 
scrapes a YouTube Live stream chat using a Python bridge.
Example: twich chat -t helpytv -y JFfPyuo67E8`,

	Run: func(cmd *cobra.Command, args []string) {
		if twitchUser == "" && ytVideoID == "" {
			log.Fatal("Please provide at least a Twitch username (-t) or a YouTube Video ID (-y)")
		}

		// Calling the new combined function from your internals package
		chat.FetchCombinedChat(twitchUser, ytVideoID)
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)

	// Defining separate flags for Twitch and YouTube
	chatCmd.Flags().StringVarP(&twitchUser, "twitch", "t", "", "Twitch username to fetch chat from")
	chatCmd.Flags().StringVarP(&ytVideoID, "youtube", "y", "", "YouTube Video ID to fetch chat from")
}