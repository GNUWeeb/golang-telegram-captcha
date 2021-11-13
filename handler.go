package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"math/rand"
	"strings"
	"time"

	gim "github.com/codenoid/goimagemerge"
	"github.com/codenoid/minikv"
	tele "gopkg.in/tucnak/telebot.v3"
)

func onJoin(c tele.Context) error {
	if c.Chat().Type == tele.ChatPrivate {
		return nil
	}

	// delete any incoming message before challenge solved
	bot.Delete(c.Message())

	// kvID is combination of user id and chat id
	kvID := fmt.Sprintf("%v-%v", c.Sender().ID, c.Chat().ID)

	// skip captcha-generation if data still exist
	if _, found := db.Get(kvID); found {
		c.Respond(&tele.CallbackResponse{Text: "Please solve existing captcha.", ShowAlert: true})
		return nil
	}

	// Go's map iteration are not ordered, but also not guaranteed
	// to be *always* randomized, so we give 1000 iteration trial
	// then stop if the 4 selected answer already filled up
	answerMoji := map[string]string{}
	for i := 0; i < 1000; i++ {
		// store 4 catpcha answer
		if len(answerMoji) == 4 {
			break
		}
		for key, val := range emojis {
			answerMoji[key] = val
			break
		}
	}

	// generate image
	captchaGrids := make([]*gim.Grid, 0)
	i := 0
	for key := range answerMoji {
		x := 10
		if i > 0 {
			x = i * 100
		}
		captchaGrids = append(captchaGrids, &gim.Grid{
			ImageFilePath: fmt.Sprintf("./assets/image/emoji/%v.png", key),
			OffsetX:       x, OffsetY: 120,
			Rotate: float64(rand.Intn(200-0) + 0),
		})
		i++
	}

	grids := []*gim.Grid{
		{
			ImageFilePath: "./gopherbg.jpg",
			Grids:         captchaGrids,
		},
	}

	rgba, _ := gim.New(grids, 1, 1).Merge()

	var img bytes.Buffer
	jpeg.Encode(&img, rgba, &jpeg.Options{Quality: 100})

	// challenge moji
	nonAnswerMoji := map[string]string{}
	for key, val := range emojis {
		if len(nonAnswerMoji) == 6 {
			break
		}
		if _, ok := answerMoji[key]; !ok {
			nonAnswerMoji[key] = val
		}
	}
	// challengeMoji contain answer and non-answer (wrong) emoji
	challengeMoji := map[string]string{}
	for key, val := range nonAnswerMoji {
		challengeMoji[key] = val
	}
	for key, val := range answerMoji {
		challengeMoji[key] = val
	}

	// generate keyboard and send image
	menu := &tele.ReplyMarkup{ResizeKeyboard: true}
	btn1 := make([]tele.Btn, 0)
	btn2 := make([]tele.Btn, 0)
	buttons := make([]tele.InlineButton, 0)

	// Go's map iterator are no ordered (randomized?)
	i = 0
	for key, emoji := range challengeMoji {
		buttons = append(buttons, tele.InlineButton{Text: emoji, Unique: key})
		if i < 5 {
			btn1 = append(btn1, menu.Data(emoji, key))
		} else {
			btn2 = append(btn2, menu.Data(emoji, key))
		}
		i++
	}

	menu.Inline(
		menu.Row(btn1...),
		menu.Row(btn2...),
	)

	file := tele.FromReader(bytes.NewReader(img.Bytes()))
	photo := &tele.Photo{File: file, ParseMode: tele.ModeMarkdown}
	photo.Caption = genCaption(c.Sender())

	msg, err := bot.Send(c.Chat(), photo, menu, tele.ModeMarkdown)
	if err == nil {
		captchaAnswer := make([]string, 0)
		for key := range answerMoji {
			captchaAnswer = append(captchaAnswer, strings.TrimSpace(key))
		}
		status := JoinStatus{
			UserID:         c.Sender().ID,
			CaptchaAnswer:  captchaAnswer,
			ChatID:         c.Chat().ID,
			CaptchaMessage: *msg,
			Buttons:        buttons,
		}

		status.UserFullName = c.Sender().FirstName + " " + c.Sender().LastName
		status.UserFullName = sanitizeName(status.UserFullName)

		db.Set(kvID, status, minikv.DefaultExpiration)

		chatMember, _ := bot.ChatMemberOf(c.Chat(), c.Sender())
		chatMember.RestrictedUntil = time.Now().Add(2 * time.Minute).Unix()
		bot.Restrict(c.Chat(), chatMember)
	}

	return nil
}

func handleAnswer(c tele.Context) error {
	if c.Chat().Type == tele.ChatPrivate {
		return nil
	}

	// kvID is combination of user id and chat id
	kvID := fmt.Sprintf("%v-%v", c.Callback().Sender.ID, c.Chat().ID)

	messageID := c.Callback().Message.ID
	answer := strings.TrimSpace(c.Callback().Data)
	answer = strings.Split(answer, "|")[0]

	status := JoinStatus{}
	if data, found := db.Get(kvID); !found {
		c.Respond(&tele.CallbackResponse{Text: "This challenge is not for you."})
		return nil
	} else {
		status = data.(JoinStatus)
	}

	if messageID != status.CaptchaMessage.ID {
		c.Respond(&tele.CallbackResponse{Text: "This challenge is not for you."})
		return nil
	}

	correct := false
	if stringInSlice(answer, status.CaptchaAnswer) {
		status.SolvedCaptcha++
		correct = true
		db.Update(kvID, status)
	} else {
		status.FailCaptcha++
	}

	newButtons := make([]tele.InlineButton, 0)
	for _, button := range status.Buttons {
		if button.Unique == answer {
			if correct {
				button.Text = "✅"
			} else {
				button.Text = "❌"
			}
		}
		newButtons = append(newButtons, button)
	}
	status.Buttons = newButtons

	db.Update(kvID, status)

	updateBtn := &tele.ReplyMarkup{
		Selective:      true,
		InlineKeyboard: [][]tele.InlineButton{},
	}
	updateBtn.InlineKeyboard = append(updateBtn.InlineKeyboard, newButtons[0:5], newButtons[5:10])
	bot.Edit(c.Callback(), updateBtn)

	if status.SolvedCaptcha >= 4 {
		db.Delete(kvID)
		c.Respond(&tele.CallbackResponse{Text: "Successfully joined.", ShowAlert: true})
		bot.Delete(&status.CaptchaMessage)

		chatMember, _ := bot.ChatMemberOf(c.Chat(), c.Sender())
		chatMember.Rights = tele.NoRestrictions()
		bot.Restrict(c.Chat(), chatMember)

		return nil
	} else if status.FailCaptcha >= 2 {
		db.Delete(kvID)
		c.Respond(&tele.CallbackResponse{Text: "Captcha failed, you have been banned, please contact admin with your another account.", ShowAlert: true})
		bot.Delete(&status.CaptchaMessage)
		bot.Ban(c.Chat(), &tele.ChatMember{User: c.Sender()}, false)

		mention := fmt.Sprintf(`[%v](tg://user?id=%v)`, status.UserFullName, status.UserID)
		msg := "Captcha failed, %v has been banned, please contact administrator if %v are real human with non-automated account"
		msg += "\n\n this message will automatically removed in 15 seconds..."
		msg = fmt.Sprintf(msg, mention, mention)
		msgr, err := bot.Send(status.CaptchaMessage.Chat, msg, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
		if err == nil {
			go func(msgr *tele.Message) {
				time.Sleep(15 * time.Second)
				bot.Delete(msgr)
			}(msgr)
		}
	}

	return nil
}

func onEvicted(key string, value interface{}) {
	if val, ok := value.(JoinStatus); ok {
		mention := fmt.Sprintf(`[%v](tg://user?id=%v)`, val.UserFullName, val.UserID)
		msg := "Captcha failed, %v has been banned, please contact administrator if %v are real human with non-automated account"
		msg += "\n\n this message will automatically removed in 15 seconds..."
		msg = fmt.Sprintf(msg, mention, mention)
		msgr, err := bot.Send(val.CaptchaMessage.Chat, msg, &tele.SendOptions{ParseMode: tele.ModeMarkdown})
		if err == nil {
			go func(msgr *tele.Message) {
				time.Sleep(15 * time.Second)
				bot.Delete(msgr)
			}(msgr)
		}

		bot.Delete(&val.CaptchaMessage)
		bot.Ban(val.CaptchaMessage.Chat, &tele.ChatMember{User: &tele.User{ID: val.UserID}}, false)
	}
}
