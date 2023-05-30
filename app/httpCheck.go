package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewMessage(chatID int64, text string) MessageConfig {
	return MessageConfig{
		BaseChat: BaseChat{
			ChatID:           chatID,
			ReplyToMessageID: 0,
		},
		Text:                  text,
		DisableWebPagePreview: true,
	}
}

// Checking the availability of the site via the HTTP protocol
func httpCheck(update uint16, bot *tgbotapi.BotAPI, group int64, site struct {
	Url      string
	Elements []string
}, timeout uint8, repeat uint8, delay float64) {
	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	for {
		errorHTML := 0
		deface := false
		for i := 0; i < int(repeat); i++ {

			start := time.Now()

			msg := tgbotapi.NewMessage(group, "Start check Site "+site.Url)
			msg.DisableWebPagePreview = true
			bot.Send(msg)

			req, err := http.NewRequest("GET", site.Url, nil)
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.3")
			resp, err := client.Do(req)

			elapsed := time.Since(start).Seconds()

			if err != nil {
				errorHTML++
			} else {
				if resp.StatusCode != 200 {
					msg := tgbotapi.NewMessage(group, "Site "+site.Url+" HTTP error. Code "+strconv.Itoa(resp.StatusCode))
					msg.DisableWebPagePreview = true
					bot.Send(msg)
					break
				}

				if elapsed >= delay {
					msg := tgbotapi.NewMessage(group, "Site "+site.Url+" HTTP delay "+strconv.FormatFloat(elapsed, 'f', 3, 32)+" sec.")
					msg.DisableWebPagePreview = true
					bot.Send(msg)
					break
				}

				bodyBytes, err := io.ReadAll(resp.Body)
				if err != nil {
					errorHTML++
				} else {
					body := string(bodyBytes)

					for _, element := range site.Elements {
						if !strings.Contains(body, element) {
							msg := tgbotapi.NewMessage(group, "Site "+site.Url+" defaced. Element '"+element+"' not found.")
							msg.DisableWebPagePreview = true
							bot.Send(msg)
							deface = true
							break
						}
					}
					if deface {
						break
					}
				}
			}
		}
		if errorHTML >= int(repeat) {
			msg := tgbotapi.NewMessage(group, "Site "+site.Url+" HTTP get error")
			msg.DisableWebPagePreview = true
			bot.Send(msg)
			errorHTML = 0
		}
		time.Sleep(time.Duration(update) * time.Second)
	}
}
