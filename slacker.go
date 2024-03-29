package main

import (
	"flag"
	"fmt"
	"github.com/james-bowman/slack"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	slackToken := getToken()

	text, err := ioutil.ReadFile("message.txt")
	if err != nil {
		log.Panic(fmt.Sprintf("Error opening message.txt for canned response: %s", err))
	}

	conn, err := slack.Connect(slackToken)

	if err != nil {
		log.Fatal(err)
	}

	slack.EventProcessor(conn,
		func(message *slack.Message) {
			message.Respond(string(text))
		},
		nil)
}

// getToken for authenticating with Slack.  Ordered lookup process trying first the command line,
// then environment variables and finally a config file
func getToken() string {
	var slackToken string

	log.Println("Looking for Slack auth token as command line argument")
	flag.StringVar(&slackToken, "slacktoken", "", "Slack authentication token - if not specified, will look for an environment variable or config file")
	flag.Parse()

	if slackToken == "" {
		log.Println("Slack auth token not found - looking for environment variable")
		slackToken = os.Getenv("SLACKER_SLACK_TOKEN")
	}

	if slackToken == "" {
		log.Println("Slack auth token not found - looking for config file")
		slackTokenFileName := "slack.token"

		slackTokenFile, err := ioutil.ReadFile(slackTokenFileName)
		if err != nil {
			log.Panic(fmt.Sprintf("Error opening slack authentication token file %s: %s", slackTokenFileName, err))
		}

		slackToken = string(slackTokenFile)
	}

	if slackToken != "" {
		log.Println("Slack auth token found")
	}

	return slackToken
}
