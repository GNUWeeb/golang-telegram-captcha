<img width="100" height="100" align="right" src="https://storage.googleapis.com/gopherizeme.appspot.com/gophers/7ae53733dbc7a071917d1f78dfd36cbe2f033c00.png">

# Go Telegram Captcha Bot

Fancy, fully-featured, easy to use, scalable Telegram CAPTCHA bot written in Go based on [tucnak's telebot](https://github.com/tucnak/telebot)

<img height="520" src="https://github.com/GNUWeeb/golang-telegram-captcha/raw/master/screenshot.png">

_this robot only has one job, and he does it well_

## Usage

1. Make sure your account got "add new admin" privilege
2. Invite *@SatpamEmojiBot* (or deploy your own) to your group
3. Make the bot as admin in the group
4. Enjoy...

## Deployment

1. Put your bot token to `BOT_TOKEN` env var
2. Run the program

## TODO

- [ ] Refactoring
- [ ] Persistent storage
- [ ] Ability to stay reliable and not spammy even on join traffic-spike
- [ ] Trust-based privilege escalation on normal user (trust from admin) 
- [ ] Integration with Telegram spam server
- [ ] Container deployment
- [ ] Write tests
- [ ] Assets licenses
