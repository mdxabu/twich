/*
Copyright © 2026 @mdxabu
*/
package cmd

import (
	// "fmt"

	chat "github.com/mdxabu/twich/internals"
	"github.com/spf13/cobra"
)

var username string

// chatCmd represents the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {

		chat.FetchChat(username)
	},
}

func init() {
	rootCmd.AddCommand(chatCmd)

	chatCmd.Flags().StringVarP(&username, "username", "u", "", "username to send stream chats")

}
