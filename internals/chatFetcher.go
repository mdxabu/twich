package chat

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/rand"
	"os/exec"
	"time"

	"github.com/fatih/color"
	"github.com/gempir/go-twitch-irc/v4"
)

var (
	MessageColor = color.New(color.FgHiWhite)
	timeColor    = color.New(color.FgYellow)
	userColorMap = make(map[string]*color.Color)
	userColors   = []color.Attribute{
		color.FgRed, color.FgCyan, color.FgGreen, color.FgMagenta,
		color.FgHiCyan, color.FgHiGreen, color.FgHiRed, color.FgHiMagenta,
	}
	ytUserColor = color.New(color.FgRed)
)

func getColoredUser(displayName string) string {
	userColor, exists := userColorMap[displayName]
	if !exists {
		randomAttr := userColors[rand.Intn(len(userColors))]
		userColor = color.New(randomAttr)
		userColorMap[displayName] = userColor
	}
	return userColor.Sprint("@" + displayName)
}

func getColoredYTUser(displayName string) string {
	return ytUserColor.Sprintd(displayName)
}

func printMessage(author, message string) {
	timestamp := timeColor.Sprint(time.Now().Format("15:04:05"))
	user := getColoredUser(author)
	msg := MessageColor.Sprint(message)
	fmt.Printf("%s [%s]: %s\n", timestamp, user, msg)
}

// FetchCombinedChat runs both Twitch and YouTube in parallel
func FetchCombinedChat(twitchUser string, ytVideoID string) {
	rand.Seed(time.Now().UnixNano())

	// 1. Start Twitch in a goroutine
	go func() {
		client := twitch.NewAnonymousClient()
		client.OnPrivateMessage(func(message twitch.PrivateMessage) {
			printMessage(message.User.DisplayName, message.Message)
		})
		client.Join(twitchUser)
		if err := client.Connect(); err != nil {
			fmt.Printf("Twitch Error: %v\n", err)
		}
	}()

	// 2. Start YouTube via Python Proxy
	go func() {
		if ytVideoID == "" {
			return
		}

		cmd := exec.Command("python3", "internals/yt_proxy.py", ytVideoID)

		// Capture standard error so we can see Python crashes
		stderr, _ := cmd.StderrPipe()
		stdout, _ := cmd.StdoutPipe()

		if err := cmd.Start(); err != nil {
			fmt.Printf("Failed to start Python: %v\n", err)
			return
		}

		// Print Python errors in a separate goroutine
		go func() {
			errScanner := bufio.NewScanner(stderr)
			for errScanner.Scan() {
				fmt.Printf("YouTube Proxy Error: %s\n", errScanner.Text())
			}
		}()

		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			var data struct {
				Author  string `json:"author"`
				Message string `json:"message"`
			}
			if err := json.Unmarshal(scanner.Bytes(), &data); err == nil {
				// Use YouTube color and @ for YouTube users
				timestamp := timeColor.Sprint(time.Now().Format("15:04:05"))
				user := getColoredYTUser(data.Author)
				msg := MessageColor.Sprint(data.Message)
				fmt.Printf("%s [%s]: %s\n", timestamp, user, msg)
			}
		}

		if err := cmd.Wait(); err != nil {
			fmt.Printf("Python script exited with error: %v\n", err)
		}
	}()

	// Keep the main thread alive
	select {}
}