package chat

import (
	"fmt"
	"math/rand"
	"time"
	"github.com/fatih/color"
	"github.com/gempir/go-twitch-irc/v4"
)

	var MessageColor = color.New(color.FgHiWhite)
		
	var userColors = []color.Attribute{
		color.FgRed,
		color.FgCyan,
		color.FgGreen,
		color.FgMagenta,
		color.FgHiCyan,
		color.FgHiGreen,
		color.FgHiRed,
		color.FgHiMagenta,
	}
	
	// Map to store color for per user
	var userColorMap = make(map[string]*color.Color)
		
		

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
		
		displayName := message.User.DisplayName
		
		userColor, exists := userColorMap[displayName]
		
		if !exists{
			randomAttr := userColors[rand.Intn(len(userColors))]
			userColor = color.New(randomAttr)
			userColorMap[displayName] = userColor
		}
		
		user := userColor.Sprintf(displayName)
		
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
