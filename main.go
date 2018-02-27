package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/github"
	"github.com/nlopes/slack"
)

func main() {
	api := slack.New("xoxb-316351224772-qity21lbvYeXEHWun94VjoCO")
	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			fmt.Print(" Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Print("Conn Counter: ", ev.ConnectionCount)
			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", ev)
				info := rtm.GetInfo()
				prefix := fmt.Sprintf("<@%s> ", info.User.ID)

				if ev.User != info.User.ID && strings.HasPrefix(ev.Text, prefix) {
					//rtm.SendMessage(rtm.NewOutgoingMessage("What's up buddy!?!?", ev.Channel))
					respond(rtm, ev, prefix)
				}
			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())
			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials\n")
				break Loop
			default:
				// Ignore other events

			}

		}
	}
}

func respond(rtm *slack.RTM, msg *slack.MessageEvent, prefix string) {
	var response string
	text := msg.Text
	text = strings.TrimPrefix(text, prefix)
	text = strings.TrimSpace(text)
	text = strings.ToLower(text)

	if text == "tip" {
		response = getJsTip()
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
	} else if text == "tip" {
		response = getJsNews()
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
	} else {
		response = "Hey man! If you want a tip type \" tip \" else type \" news \" and i will send you a link to my latest find!"
		rtm.SendMessage(rtm.NewOutgoingMessage(response, msg.Channel))
	}
}

func getJsTip() string {
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.RepositoryContentGetOptions{Ref: "master"}
	var tips []string

	_, src1, _, _ := client.Repositories.GetContents(ctx, "screenlantern", "jstips", "_posts/en/javascript", opt)
	_, src2, _, _ := client.Repositories.GetContents(ctx, "screenlantern", "30-seconds-of-code", "/snippets", opt)

	for _, val1 := range src1 {
		tips = append(tips, *val1.DownloadURL)
	}

	for _, val2 := range src2 {
		tips = append(tips, *val2.DownloadURL)
	}

	fmt.Print(tips)
	return tips[len(tips)-1]
}

func getJsNews() string {
	return "news"
}

func downloadTip(val, owner, repo string) string {
	client := github.NewClient(nil)
	ctx := context.Background()
	opt := &github.RepositoryContentGetOptions{Ref: "master"}
	resp, err := client.Repositories.DownloadContents(ctx, owner, repo, val, opt)
	if err != nil {
		log.Fatal(err)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp)
	return buf.String()
}
