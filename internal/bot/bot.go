package bot

import (
	"log"
	"time"

	"cgeba-shift-opening-alerter/internal/bot/tgusers"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Bot struct {
	*tb.Bot
	usersRepository *tgusers.Repository
}

func NewBot(token string, usersRepo *tgusers.Repository) (*Bot, error) {
	bot, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	return &Bot{
		Bot:             bot,
		usersRepository: usersRepo,
	}, err
}

func (b *Bot) Ping() error {
	return b.usersRepository.Ping()
}

func (b *Bot) AddSubscriber(user *tb.User) error {
	return b.usersRepository.Save(user)
}

func (b *Bot) SendToSubscribers(msg interface{}) error {
	users, err := b.usersRepository.FindAll()
	if err != nil {
		return err
	}
	for _, u := range users {
		_, err := b.Send(u, msg)
		if err != nil {
			log.Print(err)
		}
	}
	return nil
}
