package main

import (
	"log"

	"gopkg.in/telegram-bot-api.v4"
)

var bot *tgbotapi.BotAPI

func startBot(token string) {
	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			parseCommand(update.Message)
		} else {
			err := GotNewMessage(update.Message)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func parseCommand(msg *tgbotapi.Message) {
	var sendStr string
	switch msg.Text {
	case "/start":
		sendStr = "<b>Welcome to Project Galatea</b>"
	case "/help":
		sendStr = "<b>Help</b>"
	default:
		sendStr = "<b>Error: Unknown Commend</b>"
	}
	newMsg := tgbotapi.NewMessage(msg.Chat.ID, sendStr)
	newMsg.ParseMode = "HTML"
	bot.Send(newMsg)
}

func SendMessage(text string, chatId int64) {
	newMsg := tgbotapi.NewMessage(chatId, text)
	bot.Send(newMsg)
}
