package main

import (
	"fmt"

	tele "gopkg.in/tucnak/telebot.v3"
)

func genCaption(user *tele.User) string {
	desc := "Select all the emoji you see in the picture." +
		"\n\n Max failure: 2 mistake \n Duration: 2 minutes" +
		"\n\n Please leave group immediately if you are not ready with the bot"

	mention := fmt.Sprintf(`[%v](tg://user?id=%v)`, user.FirstName, user.ID)
	caption := fmt.Sprintf("%v, %v", mention, desc)
	return caption
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
