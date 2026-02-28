package chat

import (
	"bufio"
	_ "embed" // Required for go:embed
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fatih/color"
	"github.com/gempir/go-twitch-irc/v4"
)

//go:embed yt_proxy.py
var ytProxySource string

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
	// Fixed: Sprintd -> Sprint
	return ytUserColor.Sprint(displayName)
}

func printFormattedMessage(author, message, platform string) {
	timestamp := timeColor.Sprint(time.Now().Format("15:04:05"))
	var user string
	if platform == "yt" {
		user = getColoredYTUser(author)
	} else {
		user = getColoredUser(author)
	}
	msg := MessageColor.Sprint(message)
	fmt.Printf("%s [%s]: %s\n", timestamp, user, msg)
}

func FetchCombinedChat(twitchUser string, ytVideoID string) {
	rand.Seed(time.Now().UnixNano())

	// 1. Start Twitch
	if twitchUser != "" {
		go func() {
			client := twitch.NewAnonymousClient()
			client.OnPrivateMessage(func(message twitch.PrivateMessage) {
				printFormattedMessage(message.User.DisplayName, message.Message, "twitch")
			})
			client.Join(twitchUser)
			if err := client.Connect(); err != nil {
				fmt.Printf("Twitch Error: %v\n", err)
			}
		}()
	}

	// 2. Start YouTube via Embedded Python Proxy
	if ytVideoID != "" {
		go func() {
			// Write the embedded script to a temporary file
			tmpFile := filepath.Join(os.TempDir(), "twich_yt_proxy.py")
			err := os.WriteFile(tmpFile, []byte(ytProxySource), 0644)
			if err != nil {
				fmt.Printf("Failed to setup YT proxy: %v\n", err)
				return
			}
			defer os.Remove(tmpFile)

			cmd := exec.Command("python3", tmpFile, ytVideoID)
			stderr, _ := cmd.StderrPipe()
			stdout, _ := cmd.StdoutPipe()

			if err := cmd.Start(); err != nil {
				fmt.Printf("Failed to start Python: %v\n", err)
				return
			}

			// Capture Python errors
			go func() {
				errScanner := bufio.NewScanner(stderr)
				for errScanner.Scan() {
					fmt.Printf("YouTube Error: %s\n", errScanner.Text())
				}
			}()

			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				var data struct {
					Author  string `json:"author"`
					Message string `json:"message"`
				}
				if err := json.Unmarshal(scanner.Bytes(), &data); err == nil {
					printFormattedMessage(data.Author, data.Message, "yt")
				}
			}
			cmd.Wait()
		}()
	}

	select {}
}