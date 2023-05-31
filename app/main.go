package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		Update uint16
	}
	Telegram struct {
		Token string
		Group int64
	}
	Http struct {
		Repeat  uint8
		Timeout uint8
		Delay   float64
		Sites   []struct {
			Url      string
			Elements []string
		}
	}
}

func main() {
	// Read config from yaml
	config := Config{}
	filename, _ := filepath.Abs("./conf/config.yaml")
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	// Parse yaml
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	// Telegram bot
	bot, err := tgbotapi.NewBotAPI(config.Telegram.Token)
	if err != nil {
		panic(err)
	}

	// Debug Telegram bot
	bot.Debug = false

	// Running HTTP checker
	siteIndex := 0
	siteTotal := len(config.Http.Sites)

	for _, site := range config.Http.Sites {
		siteIndex++
		go httpCheck(config.App.Update, bot, config.Telegram.Group, site, config.Http.Timeout, config.Http.Repeat, config.Http.Delay, siteIndex, siteTotal)
	}

	botUpdate(bot, config.Http.Sites, config)
}

// Telegram bot for listening to incoming commands
func botUpdate(bot *tgbotapi.BotAPI, sites []struct{
	Url      string
	Elements []string
}, config Config) {

	// Create string for HTTP(s) monitoring sites
	sitesString := ""
	for _, site := range sites {
		sitesString += site.Url + "\n"
	}

	// Telegram bot listener
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 300
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		msg.DisableWebPagePreview = true

		switch update.Message.Command() {
		case "start":
			msg.Text = "Hi, I am a monitoring bot! Your (group) ID = " + strconv.FormatInt(update.Message.Chat.ID, 10)
		case "list":
			msg.Text = "Listed HTTP(s) monitoring sites:\n" + sitesString
		case "check":
			msg.Text = "Domain check:"
			for _, site := range config.Http.Sites {
				go httpCheckOnline(config.App.Update, bot, config.Telegram.Group, site, config.Http.Timeout, config.Http.Repeat, config.Http.Delay)
			}
		default:
			msg.Text = "I don't know that command"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}
