package chat

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/gempir/go-twitch-irc/v4"
)

//go:embed yt_proxy.py
var ytProxySource string

var (
	MessageColor = color.New(color.FgHiWhite)
	timeColor    = color.New(color.FgYellow)

	// Twitch user color palette.
	twUserColorMap = make(map[string]*color.Color)
	twUserColors   = []color.Attribute{
		color.FgRed, color.FgCyan, color.FgGreen, color.FgMagenta,
		color.FgHiCyan, color.FgHiGreen, color.FgHiRed, color.FgHiMagenta,
	}

	// YouTube user color palette — intentionally different from Twitch
	// so you can visually distinguish the two at a glance.
	ytUserColorMap = make(map[string]*color.Color)
	ytUserColors   = []color.Attribute{
		color.FgHiYellow, color.FgHiBlue, color.FgYellow, color.FgBlue,
		color.FgHiWhite, color.FgWhite, color.FgHiCyan, color.FgCyan,
	}

	// printMu prevents interleaved output when Twitch and YouTube
	// goroutines write at the same time.
	printMu sync.Mutex
)

func getColoredTwUser(displayName string) string {
	c, exists := twUserColorMap[displayName]
	if !exists {
		c = color.New(twUserColors[rand.Intn(len(twUserColors))])
		twUserColorMap[displayName] = c
	}
	return c.Sprint("@" + displayName)
}

func getColoredYTUser(displayName string) string {
	c, exists := ytUserColorMap[displayName]
	if !exists {
		c = color.New(ytUserColors[rand.Intn(len(ytUserColors))])
		ytUserColorMap[displayName] = c
	}
	return c.Sprint("@" + displayName)
}

// printMessage formats and immediately prints a chat message to stdout.
// A mutex ensures lines from concurrent goroutines never interleave.
func printMessage(author, message, platform string) {
	timestamp := timeColor.Sprint(time.Now().Format("15:04:05"))

	var user string
	if platform == "yt" {
		user = getColoredYTUser(author)
	} else {
		user = getColoredTwUser(author)
	}

	msg := MessageColor.Sprint(message)

	var platformTag string
	if platform == "yt" {
		platformTag = color.New(color.FgRed, color.Bold).Sprint("[YT]")
	} else {
		platformTag = color.New(color.FgHiMagenta, color.Bold).Sprint("[TW]")
	}

	line := fmt.Sprintf("%s %s [%s]: %s", timestamp, platformTag, user, msg)

	printMu.Lock()
	fmt.Println(line)
	printMu.Unlock()
}

// FetchCombinedChat connects to Twitch and/or YouTube live chat.
// ytChannel accepts a YouTube @handle (e.g. "@mkbhd" or "mkbhd"),
// a full channel URL, or a raw video ID. The Python proxy resolves
// the handle to the active live stream video ID automatically.
func FetchCombinedChat(twitchUser string, ytChannel string) {
	rand.Seed(time.Now().UnixNano())

	// header := color.New(color.FgHiYellow, color.Bold).Sprintf("═══ Live Chat ═══")
	// fmt.Println(header)
	// fmt.Println()

	// 1. Start Twitch
	if twitchUser != "" {
		go func() {
			client := twitch.NewAnonymousClient()
			client.OnPrivateMessage(func(message twitch.PrivateMessage) {
				printMessage(message.User.DisplayName, message.Message, "twitch")
			})
			client.Join(twitchUser)
			if err := client.Connect(); err != nil {
				// Silently log to stderr; don't pollute the chat view.
				fmt.Fprintf(os.Stderr, "Twitch Error: %v\n", err)
			}
		}()
	}

	// 2. Start YouTube via Embedded Python Proxy
	// The proxy accepts a @handle, channel URL, or raw video ID and
	// resolves it to the currently active live stream automatically.
	if ytChannel != "" {
		go func() {
			// Write the embedded script to a temporary file.
			tmpFile := filepath.Join(os.TempDir(), "twich_yt_proxy.py")
			err := os.WriteFile(tmpFile, []byte(ytProxySource), 0644)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to setup YT proxy: %v\n", err)
				return
			}
			defer os.Remove(tmpFile)

			cmd := exec.Command("python3", tmpFile, ytChannel)
			stderr, _ := cmd.StderrPipe()
			stdout, _ := cmd.StdoutPipe()

			if err := cmd.Start(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to start Python: %v\n", err)
				return
			}

			// Silently consume Python stderr (resolution logs, errors)
			// so they never appear in the chat output.
			go func() {
				errScanner := bufio.NewScanner(stderr)
				for errScanner.Scan() {
					// Discard — these are internal resolution/debug messages.
					_ = errScanner.Text()
				}
			}()

			// Read chat messages from stdout.
			scanner := bufio.NewScanner(stdout)
			for scanner.Scan() {
				var data struct {
					Author  string `json:"author"`
					Message string `json:"message"`
				}
				if err := json.Unmarshal(scanner.Bytes(), &data); err == nil {
					printMessage(data.Author, data.Message, "yt")
				}
			}
			cmd.Wait()
		}()
	}

	select {}
}
