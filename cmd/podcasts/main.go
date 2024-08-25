package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/mp3"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/mramsden/podcasts/internal/rss"
)

func main() {
	cmd := os.Args[1]

	switch cmd {
	case "play":
		PlayChannel()
	}
}

func PlayChannel() {
	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s play <url>\n", os.Args[0])
		os.Exit(1)
	}
	feedUrl := os.Args[2]

	channels, err := downloadFeed(feedUrl)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	if len(channels) == 0 {
		return
	}

	if len(channels[0].Items) == 0 {
		return
	}

	item := channels[0].Items[0]

	err = play(item.Enclosure.URL)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
}

func downloadFeed(url string) ([]rss.Channel, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	channels, err := rss.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return *channels, nil
}

func play(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	streamer, format, err := mp3.Decode(resp.Body)
	if err != nil {
		return err
	}
	defer streamer.Close()

	err = speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	if err != nil {
		return err
	}

	done := make(chan bool)
	speaker.Play(beep.Seq(streamer, beep.Callback(func() {
		done <- true
	})))
	<-done

	return nil
}
