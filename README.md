# Network Calculator Telegram Bot

## Requirements
- A Go compiler
- A Telegram bot token
---
## Installation
Install the Go-Telegram-Bot-Api library

```
go get -u github.com/go-telegram-bot-api/telegram-bot-api
```

Clone the bot repository

```
git clone https://github.com/Knocks83/go-Telegram-NetworkCalculator-Bot.git
```
---
## Configuration
To start using the bot you have to configure it by editing the `config/config.go` file. In that file you have to set:

- `Token` to your bot token.
- `RolesFile` (optional) if you want to change the path of the JSON file that'll contain the admins and banned people.
- `LogChat` (optional) to the ChatID of the chat you'll use as log.

You also have to add your UserID to the `roles.json` file, so you'll be able to use admin-only commands and add other people to the admin list directly from Telegram.
