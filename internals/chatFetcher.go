package chat

import (
	"fmt"

	"github.com/gempir/go-twitch-irc/v4"
)

// First gonna try with anonymous mode
func FetchChat(username string) {
	client := twitch.NewAnonymousClient()

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		fmt.Println(message.Time.Local().Format("15:04:05")+" [" + message.User.DisplayName + "]: " + message.Message)
	})

	client.Join(username)

	err := client.Connect()
	if err != nil {
		panic(err)
	}

}
