package main

import tele "gopkg.in/tucnak/telebot.v3"

type JoinStatus struct {
	UserID         int64
	UserFullName   string
	CaptchaAnswer  []string
	SolvedCaptcha  int
	FailCaptcha    int
	ChatID         int64
	CaptchaMessage tele.Message
	Buttons        []tele.InlineButton
}
