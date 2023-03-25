package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	_ "net/http"
	"regexp"
	"strings"
	"time"
)

type MessageResp struct {
	Messages []string `json:"messages"`
	MinSeq   int      `json:"min_seq"`
}

func main() {
	// Create a new WebSocket dialer
	dialer := websocket.DefaultDialer

	// Connect to WSS URL
	//wss://tch7578.tch.quora.com/up/chan50-8888/updates?min_seq=3428346961&channel=poe-chan50-8888-kgkbixjbmqrvykhqexay&hash=897287112690331892
	//

	conn, _, err := dialer.Dial("wss://tch505872.tch.quora.com/up/chan50-8888/updates?min_seq=3442885839&channel=poe-chan50-8888-kgkbixjbmqrvykhqexay&hash=897287112690331892", nil)
	if err != nil {
		log.Fatal("WebSocket dial error:", err)
	}

	// Continuously read incoming messages
	var messageResp MessageResp
	re := regexp.MustCompile(`"text":"(.*?)"`)
	for {
		err := conn.ReadJSON(&messageResp)

		if err != nil {
			log.Println("WebSocket read error:", err)
			return
		}
		//dispaly消息

		messages := messageResp.Messages
		for _, message := range messages {
			if strings.Contains(message, "displayName") {
				continue
			}
			if strings.Contains(message, "\"state\":\"complete\"") && strings.Contains(message, "\"suggestedReplies\":[]") {
				submatch := re.FindStringSubmatch(message)
				//fmt.Println(message)
				for _, c := range submatch[1] {
					fmt.Print(string(c))
					time.Sleep(20 * time.Millisecond)
				}
				fmt.Print("\n")
			} else {
				continue
			}
		}

		fmt.Printf("for next message,late min_seq = %d\n", messageResp.MinSeq)
	}
}
