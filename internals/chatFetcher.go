package chat

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
	"github.com/gempir/go-twitch-irc/v4"
)

	var MessageColor = color.New(color.FgHiWhite)
		
	var userColor = []color.Attribute{
		color.FgRed,
		// color.FgBlue,
		color.FgCyan,
		color.FgGreen,
		color.FgMagenta,
		// color.FgHiBlue,
		color.FgHiCyan,
		color.FgHiGreen,
		color.FgHiRed,
		color.FgHiMagenta,
	}
		
		

// First gonna try with anonymous mode
func FetchChat(username string) {
	client := twitch.NewAnonymousClient()
	
	rand.Seed(time.Now().UnixNano())

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		
		// timestamp color in yellow
		timeColor := color.New(color.FgYellow)
		timestamp := timeColor.Sprintf(
			message.Time.Local().Format("15:04:05"),
		)
		
		// randcom color for the usernames/users
		randomColor := userColor[rand.Intn(len(userColor))]
		userColor := color.New(randomColor)
		
		user := userColor.Sprintf(message.User.DisplayName)
		
		// message color
		msg := MessageColor.Sprintf(message.Message)
		
		fmt.Printf("%s [%s]: %s\n",timestamp,user,msg)
		
	})

	client.Join(username)

	err := client.Connect()
	if err != nil {
		panic(err)
	}

}
