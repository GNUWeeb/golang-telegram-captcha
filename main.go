package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/codenoid/minikv"
	tele "gopkg.in/tucnak/telebot.v3"
)

var (
	bot *tele.Bot

	db = minikv.New(2*time.Minute, 5*time.Second)
)

func main() {

	botToken := os.Getenv("BOT_TOKEN")

	// listen for janitor expiration removal ( 5*time.Second )
	db.OnEvicted(onEvicted)

	b, err := tele.NewBot(tele.Settings{
		Token:  botToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	bot = b

	b.Handle("/testcaptcha", onJoin)
	b.Handle(tele.OnUserJoined, onJoin)
	b.Handle(tele.OnCallback, handleAnswer)
	b.Handle(tele.OnUserLeft, func(c tele.Context) error {
		c.Delete()
		db.Delete(fmt.Sprint(c.Sender().ID))
		return nil
	})

	b.Start()
}
