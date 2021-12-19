package main

import (
	"fmt"
	"log"
	"time"

	"cgeba-shift-opening-alerter/internal/bot"
	"cgeba-shift-opening-alerter/internal/dependencies"
	"cgeba-shift-opening-alerter/internal/shifttable"

	tb "gopkg.in/tucnak/telebot.v2"
)

type sendMsg func(msg interface{}) error

func main() {
	container := dependencies.BuildContainer()
	bot, err := container.GetBot()
	if err != nil {
		log.Fatal(err)
	}
	c := make(chan error, 1)
	go listen(bot, c)
	go func() {
		for {
			msg := checkUpdates()
			if msg != nil {
				err := bot.SendToSubscribers(msg)
				if err != nil {
					c <- err
				}
			}
			msg = checkStatus(bot)
			err = bot.SendToSubscribers(msg)
			if err != nil {
				c <- err
				return
			}
			now := time.Now()
			t := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
			if now.After(t) {
				t = t.Add(24 * time.Hour)
			}
			time.Sleep(time.Until(t))
		}
	}()
	err = <-c
	log.Fatal(err.Error())
}

func listen(bot *bot.Bot, c chan error) {
	bot.Handle("/start", func(m *tb.Message) {
		s := m.Sender
		msg := `/subscribe Subscribe to the news
/status Check status
/updates Check for updates`
		_, err := bot.Send(s, msg)
		if err != nil {
			c <- err
		}
	})
	bot.Handle("/subscribe", func(m *tb.Message) {
		s := m.Sender
		msg := subscribe(s, bot)
		_, err := bot.Send(s, msg)
		if err != nil {
			c <- err
		}
	})
	bot.Handle("/status", func(m *tb.Message) {
		msg := checkStatus(bot)
		_, err := bot.Send(m.Sender, msg)
		if err != nil {
			c <- err
		}
	})
	bot.Handle("/updates", func(m *tb.Message) {
		msg := checkUpdates()
		var err error
		if msg != nil {
			_, err = bot.Send(m.Sender, msg)
		} else {
			_, err = bot.Send(m.Sender, "There aren't updates")
		}
		if err != nil {
			c <- err
		}
	})
	bot.Start()
}

func subscribe(s *tb.User, bot *bot.Bot) string {
	err := bot.AddSubscriber(s)
	if err != nil {
		return fmt.Sprintf("‼️ERROR: %q", err.Error())
	}
	return "Te he suscrito para las nuevas fechas de matrimonios"
}

func checkUpdates() interface{} {
	table, err := shifttable.GetTable()
	if err != nil {
		return fmt.Sprintf("‼️ Error getting the updates: %q", err.Error())
	}
	if date := table["Expedición de pasaportesrenovación y primera vez"].NextOpening; date != "fecha a confirmar" {
		return "🎉 Los turnos para pasaportes se abriran el " + date
	}
	return nil
}

func checkStatus(bot *bot.Bot) string {
	err := bot.Ping()
	if err != nil {
		return fmt.Sprintf("❌ Error: %q", err.Error())
	}
	return "✔️ The bot is working"
}
