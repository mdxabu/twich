package internals

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"

	YtChat "github.com/abhinavxd/youtube-live-chat-downloader/v2"
)

func main() {
	customCookies := []*http.Cookie{
		{Name: "PREF",
			Value:  "tz=Europe.Rome",
			MaxAge: 300},
		{Name: "CONSENT",
			Value:  fmt.Sprintf("YES+yt.432048971.it+FX+%d", 100+rand.Intn(999-100+1)),
			MaxAge: 300},
	}
	YtChat.AddCookies(customCookies)

	continuation, cfg, error := YtChat.ParseInitialData("https://www.youtube.com/watch?v=EB916GNkAsQ")
	if error != nil {
		log.Fatal(error)
	}
	for {
		chat, newContinuation, error := YtChat.FetchContinuationChat(continuation, cfg)
		if error == YtChat.ErrLiveStreamOver {
			log.Fatal("Live stream over")
		}
		if error != nil {
			log.Print(error)
			continue
		}
		continuation = newContinuation

		for _, msg := range chat {
			fmt.Print(msg.Timestamp, " | ")
			fmt.Println(msg.AuthorName, ": ", msg.Message)
		}
	}
}

// ListYTComments fetches and prints YouTube comments for a given video URL.
func ListYTComments(url string) error {
	customCookies := []*http.Cookie{
		{Name: "PREF",
			Value:  "tz=Europe.Rome",
			MaxAge: 300},
		{Name: "CONSENT",
			Value:  fmt.Sprintf("YES+yt.432048971.it+FX+%d", 100+rand.Intn(999-100+1)),
			MaxAge: 300},
	}
	YtChat.AddCookies(customCookies)

	continuation, cfg, err := YtChat.ParseInitialData(url)
	if err != nil {
		return err
	}
	for {
		chat, newContinuation, err := YtChat.FetchContinuationChat(continuation, cfg)
		if err == YtChat.ErrLiveStreamOver {
			log.Print("Live stream over")
			break
		}
		if err != nil {
			log.Print(err)
			break
		}
		continuation = newContinuation

		for _, msg := range chat {
			fmt.Print(msg.Timestamp, " | ")
			fmt.Println(msg.AuthorName, ": ", msg.Message)
		}
	}
	return nil
}
