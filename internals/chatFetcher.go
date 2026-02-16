package chat

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/gempir/go-twitch-irc/v4"
)

// First gonna try with anonymous mode
func FetchChat(username string) {
	client := twitch.NewAnonymousClient()

	// Define colors for timestamp and username
	cyan := color.New(color.FgCyan).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		timestamp := cyan(message.Time.Local().Format("15:04:05"))
		displayName := yellow(message.User.DisplayName)
		fmt.Printf("%s [%s]: %s\n", timestamp, displayName, message.Message)
	})

	client.Join(username)

	err := client.Connect()
	if err != nil {
		panic(err)
	}

}
