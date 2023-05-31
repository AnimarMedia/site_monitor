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

// Checking the availability of the site via the HTTP protocol
func httpCheckOnline(update uint16, bot *tgbotapi.BotAPI, group int64, site struct {
	Url      string
	Elements []string
}, timeout uint8, repeat uint8, delay float64) {
	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	for {
		currentTime := time.Now()
		errorHTML := 0
		deface := false
		for i := 0; i < int(repeat); i++ {

			start := time.Now()

			req, err := http.NewRequest("GET", site.Url, nil)
			if err != nil {
				log.Fatalf("error: %v", err)
			}
			req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.3")
			resp, err := client.Do(req)

			// fmt.Println("siteIndex :", siteIndex)
			// fmt.Println("siteTotal :", siteTotal)

			elapsed := time.Since(start).Seconds()

			if err != nil {
				errorHTML++
			} else {
				if resp.StatusCode != 200 {
					msg := tgbotapi.NewMessage(group, "["+currentTime.Format("2006.01.02 15:04:05")+"] "+site.Url+" [ERROR]. Code "+strconv.Itoa(resp.StatusCode))
					msg.DisableWebPagePreview = true
					bot.Send(msg)
					break
				}

				if elapsed >= delay {
					msg := tgbotapi.NewMessage(group, "["+currentTime.Format("2006.01.02 15:04:05")+"] "+site.Url+" [ONLINE] "+strconv.FormatFloat(elapsed, 'f', 3, 32)+" sec.")
					msg.DisableWebPagePreview = true
					bot.Send(msg)
					break
				} else {
					msg := tgbotapi.NewMessage(group, "["+currentTime.Format("2006.01.02 15:04:05")+"] "+site.Url+" [ONLINE] "+strconv.FormatFloat(elapsed, 'f', 3, 32)+" sec.")
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
							msg := tgbotapi.NewMessage(group, "["+currentTime.Format("2006.01.02 15:04:05")+"] "+site.Url+" defaced. Element '"+element+"' not found.")
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
			msg := tgbotapi.NewMessage(group, "["+currentTime.Format("2006.01.02 15:04:05")+"] "+site.Url+" [ERROR] HTTP get error")
			msg.DisableWebPagePreview = true
			bot.Send(msg)
			errorHTML = 0
		}
		time.Sleep(time.Duration(update) * time.Second)
	}
}
