package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/nlopes/slack"
)

var (
	apiReady    chan struct{}
	apiInterval int
)

func fatal(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(1)
}

func initApiThrottling() {
	apiReady = make(chan struct{})
	go func() {
		for {
			apiReady <- struct{}{}
			<-time.After(time.Duration(apiInterval) * time.Second)
		}
	}()
}

func countChannel(api *slack.Client, ch *slack.Channel) (count int) {
	h := &slack.History{HasMore: true}
	param := slack.NewHistoryParameters()
	for h.HasMore {
		var err error
		<-apiReady
		h, err = api.GetChannelHistory(ch.ID, param)
		if err != nil {
			fatal("F: GetChannelsHistory failed: %v", err)
		}
		if len(h.Messages) > 0 {
			count += len(h.Messages)
			param.Latest = h.Messages[len(h.Messages)-1].Timestamp
		}
		if h.HasMore {
			fmt.Fprintf(os.Stderr, ".")
		}
	}
	return
}

type channel struct {
	Name  string
	Count int
}

type channelSlice []channel

func (ch channelSlice) Len() int {
	return len(ch)
}

func (ch channelSlice) Less(i, j int) bool {
	return ch[i].Count > ch[j].Count
}

func (ch channelSlice) Swap(i, j int) {
	ch[i], ch[j] = ch[j], ch[i]
}

func init() {
	flag.IntVar(&apiInterval, "i", 1, "interval (sec) for api call")
	initApiThrottling()
}

func main() {
	flag.Parse()

	token := os.Getenv("SLACK_API_TOKEN")
	if token == "" {
		fatal("Environment value SLACK_API_TOKEN is not set")
	}

	api := slack.New(token)
	<-apiReady
	atr, err := api.AuthTest()
	if err != nil {
		fatal("AuthTest() failed: %v", err)
	}
	fmt.Fprintf(os.Stderr, "Connected to %s as %s\n", atr.Team, atr.User)

	<-apiReady
	channels, err := api.GetChannels(true)
	if err != nil {
		fatal("F: GetChannels(true) failed: %v", err)
	}
	var total int
	cs := make([]channel, len(channels))

	fmt.Fprintf(os.Stderr, "There are %d channels\n", len(channels))
	for i, ch := range channels {
		cs[i].Name = ch.Name
		fmt.Fprintf(os.Stderr, "(%d/%d) %s: ", i+1, len(channels), ch.Name)
		cs[i].Count = countChannel(api, &ch)
		fmt.Fprintf(os.Stderr, " %d\n", cs[i].Count)
		total += cs[i].Count
	}
	fmt.Fprintf(os.Stderr, "\n")

	sort.Stable(channelSlice(cs))

	fmt.Printf("%6d TOTAL\n", total)
	for _, c := range cs {
		fmt.Printf("%6d #%s\n", c.Count, c.Name)
	}
}
