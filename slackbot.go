package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"log"
	"math/rand"
	"strings"
	"time"
)

const (
	MYTOKEN = "xoxb-71342955047-6W175TjjXECUsVL0tpEFUBMe"
)

type lexElement []string

func clean_symbols(s string) (result string) {
	for _, v := range strings.ToLower(s) {
		if rune(v) >= rune('a') && rune(v) <= rune('z') {
			result += string(v)
		}
	}
	return
}

func (l lexElement) in(s string) (result bool, pos int) {
	s_words := strings.Fields(s)
	for i, v := range s_words {
		s_words[i] = clean_symbols(v)
	}
	for _, v_alternative := range l {
		words_to_find := strings.Fields(v_alternative)
	S_WORDS_LOOP:
		for i_in_s_words, _ := range s_words {
			if i_in_s_words > len(s_words)-len(words_to_find) {
				break
			}
			for i := 0; i < len(words_to_find); i++ {
				if words_to_find[i] != s_words[i_in_s_words+i] {
					continue S_WORDS_LOOP
				}
			}
			result = true
			pos = i_in_s_words
			return
		}
	}
	return
}

func get_alternative(l lexElement) (result string) {
	r := rand.New(rand.NewSource(int64(time.Now().Second())))
	index := int(r.Uint32() % uint32(len(l)))
	return l[index]
}

func containsInOrder(where string, what ...lexElement) (result bool) {
	last_index := 0
	for i, v := range what {
		found, pos := v.in(where)
		if !found {
			return
		}
		if i > 0 && pos < last_index {
			return
		}
	}
	result = true
	return
}

var LE_greering = lexElement{
	"hello",
	"hi",
	"howdy",
	"how is it going",
	"how are you doing",
	"what's up",
	"what up",
	"how are you",
	"what's the deal",
	"whats the deal",
}

var LE_bot = lexElement{
	"mybot",
	"bot",
	"robot",
}

var LE_doing_good = lexElement{
	"I'm cool",
	"I'm doing great",
	"I'm fine",
	"I'm ok",
	"It's a great day",
	"I'm dandy",
}

func main() {

	slackapi := slack.New(MYTOKEN)
	// If you set debugging, it will log all requests to the console
	// Useful when encountering issues
	// api.SetDebug(true)

	channels, err := slackapi.GetChannels(true)
	if err != nil {
		log.Fatal("Error getting slack channels:", err)
		return
	}
	channelsByName := make(map[string]slack.Channel)
	channelsByID := make(map[string]slack.Channel)
	for _, c := range channels {
		channelsByName[c.Name] = c
		channelsByID[c.ID] = c
	}

	users, err := slackapi.GetUsers()
	if err != nil {
		log.Fatal("Error getting slack users:", err)
		return
	}
	usersByName := make(map[string]slack.User)
	usersByID := make(map[string]slack.User)
	for _, u := range users {
		usersByName[u.Name] = u
		usersByID[u.ID] = u
	}

	/*
	   var general *slack.Channel
	   for _, c := range channels {
	       if c.Name == "general" {
	           general = &c
	           break
	       }
	   }
	   if general == nil {
	       fmt.Println("Did not find channel general")
	       return
	   }
	*/

	/*
	   fmt.Println("Now posting")
	   params := slack.PostMessageParameters{}
	   channelID, timestamp, err := api.PostMessage(general.ID, "Yo there", params)
	   if err != nil {
	       log.Fatalf("%s\n", err)
	   }
	   log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
	   parameters := slack.NewHistoryParameters()
	   parameters.Unreads = true
	   parameters.Inclusive = true
	   history, err := api.GetChannelHistory(general.ID, parameters)
	   if err != nil {
	       fmt.Println(err)
	       return
	   }
	   for i, m := range history.Messages {
	       fmt.Printf("Message %v from user %v (%v) is: %v\n", i, m.User, m.Username, m.Text)
	   }
	*/

	rtm := slackapi.NewRTM()
	go rtm.ManageConnection()

	/*
		params := slack.PostMessageParameters{AsUser: true}
		_, _, err = rtm.PostMessage(general.ID, "Berf!", params)
		if err != nil {
			log.Fatal(err)
		}
	*/
Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.InvalidAuthEvent:
				log.Printf("Invalid credentials")
				break Loop
			case *slack.MessageEvent:
				log.Printf("Message from %v in channel %v: %v\n", usersByID[ev.User].Name, channelsByID[ev.Channel].Name, ev.Text)
				if containsInOrder(ev.Text, LE_greering, LE_bot) {
					log.Println("Detected a greeting meesage")
					params := slack.PostMessageParameters{AsUser: true}
					_, _, err = rtm.PostMessage(channelsByID[ev.Channel].Name, fmt.Sprintf("%v %v?", get_alternative(LE_greering), usersByID[ev.User].Name), params)
				}
			case *slack.LatencyReport, *slack.ReconnectUrlEvent:
				// Ignore
				// log.Printf("Current latency: %v\n", ev.Value)
				// log.Printf("Reonnect_url event, urs is", ev.URL)
			default:
				log.Printf("Slack event, type %v, data %v\n", msg.Type, ev)
				/*
					case *slack.HelloEvent:
						// Ignore hello
					case *slack.ConnectedEvent:
						log.Println("Infos:", ev.Info)
						log.Println("Connection counter:", ev.ConnectionCount)
						// Replace #general with your Channel ID
						rtm.SendMessage(rtm.NewOutgoingMessage("Hello world", "#general"))

					case *slack.PresenceChangeEvent:
						log.Printf("Presence Change: %v\n", ev)

					case *slack.RTMError:
						log.Printf("Error: %s\n", ev.Error())
				*/
			}
		}
	}

}
