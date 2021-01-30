/*
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package telegram

import (
	"go-Telegram-Network-Bot/network"
	"go-Telegram-Network-bot/config"
	"strconv"
	"strings"
	"fmt"

	//"encoding/binary"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// MARKDOWN V2
// *bold \*text*
// _italic \*text_
// __underline__
// ~strikethrough~
// *bold _italic bold ~italic bold strikethrough~ __underline italic bold___ bold*
// [inline URL](http://www.example.com/)
// [inline mention of a user](tg://user?id=123456789)
// `inline fixed-width code`
// ```
// pre-formatted fixed-width code block
// ```
// ```python
// pre-formatted fixed-width code block written in the Python programming language
// ```
//
//
// Any character with code between 1 and 126 inclusively can be escaped anywhere with a preceding '\' character, in which case it is treated as an ordinary character and not a part of the markup.
// Inside pre and code entities, all '`' and '\' characters must be escaped with a preceding '\' character.
// Inside (...) part of inline link definition, all ')' and '\' must be escaped with a preceding '\' character.
// In all other places characters '_', '*', '[', ']', '(', ')', '~', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!' must be escaped with the preceding character '\'.
// In case of ambiguity between italic and underline entities __ is always greadily treated from left to right as beginning or end of underline entity, so instead of ___italic underline___ use ___italic underline_\r__, where \r is a character with code 13, which will be ignored.

func (tg *Telegram) HandleUpdate(update tgbotapi.Update) {
	// Check for ban and inform the user only on private chat to avoid flood
	if tg.db.FindBan(int64(update.Message.From.ID)) >= 0 {
		if update.Message.Chat.Type == "private" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🚫 You have been banned from this bot!")
			msg.ReplyToMessageID = update.Message.MessageID
			tg.api.Send(msg)
		}

		return
	}

	// Skip messages from channels
	if update.Message.Chat.Type == "channel" {
		return
	}

	if update.CallbackQuery != nil {
		text := strings.Split(update.CallbackQuery.Data, " ")
		switch text[0] {
		case "first":
			break
		case "file":
			break
		}
	}

	// Skip if there isn't a real update
	if update.Message == nil && update.EditedMessage == nil {
		return
	}

	// If message is edited, set it as message to handle
	if update.EditedMessage != nil {
		update.Message = update.EditedMessage
	}

	// Commands for admins only
	if tg.db.FindAdmin(int64(update.Message.From.ID)) >= 0 {
		if update.Message.Text == "/ping" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🏓. Next time don't break my balls!")
			msg.ReplyToMessageID = update.Message.MessageID
			_, _ = tg.api.Send(msg)
			return
		}

		if len(update.Message.Text) >= 5 && strings.ToLower(update.Message.Text[1:5]) == "help" && update.Message.Chat.Type == "private" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "*HELP*\n\n*/ping* \\- send a test message\\.\n*/admin* \\- When you reply to a person message, he will become an admin\\.\n*/unadmin* \\- remove an admin\\.\n*/ban* \\- ban a person from the bot\\.\n*/unban* \\- unban a person from the bot\\.\n*/comestero <known key\\> <known key sector \\(0\\-15\\)\\> <known key type \\(A/B\\)\\>* \\- generate keys for a comestero vending key\\.")
			msg.ParseMode = tgbotapi.ModeMarkdownV2
			_, _ = tg.api.Send(msg)
			return
		}

		if update.Message.Text == "/admin" {
			// Check if the user is replying to a message
			if update.Message.ReplyToMessage != nil {
				err := tg.db.AddAdmin(int64(update.Message.ReplyToMessage.From.ID))
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR: "+err.Error()+".")
					msg.ReplyToMessageID = update.Message.MessageID
					_, _ = tg.api.Send(msg)
					return
				}

				// Success message
				adminSuccess := tgbotapi.NewMessage(update.Message.Chat.ID, "User [`"+update.Message.ReplyToMessage.From.FirstName+"`](tg://user?id="+strconv.Itoa(update.Message.ReplyToMessage.From.ID)+") is now an admin\\!")
				adminSuccess.ParseMode = tgbotapi.ModeMarkdownV2
				_, err = tg.api.Send(adminSuccess)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR sending the success admin message: "+err.Error()+".")
					msg.ReplyToMessageID = update.Message.MessageID
					tg.api.Send(msg)
				}

				// Log action
				adminSuccess.Text = adminSuccess.Text + "\n\nCommand issued by [" + update.Message.From.FirstName + "](tg://user?id=" + strconv.Itoa(update.Message.From.ID) + ")"
				adminSuccess.ChatID = config.LogChat
				_, _ = tg.api.Send(adminSuccess)
			} else {
				_, _ = tg.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You are using this command in the wrong way! Reply to a message to make him admin!"))
			}

			return
		}

		if update.Message.Text == "/unadmin" {
			// Check if the user is replying to a message
			if update.Message.ReplyToMessage != nil {
				err := tg.db.RemoveAdmin(int64(update.Message.ReplyToMessage.From.ID))

				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR: "+err.Error()+".")
					msg.ReplyToMessageID = update.Message.MessageID
					_, _ = tg.api.Send(msg)
					return
				}

				// Success message
				adminRemove := tgbotapi.NewMessage(update.Message.Chat.ID, "User ["+update.Message.ReplyToMessage.From.FirstName+"](tg://user?id="+strconv.FormatInt(int64(update.Message.ReplyToMessage.From.ID), 10)+") is no longer an admin\\!")
				adminRemove.ParseMode = tgbotapi.ModeMarkdownV2
				_, _ = tg.api.Send(adminRemove)

				// Log action
				adminRemove.Text += "\n\nCommand issued by [" + update.Message.From.FirstName + "](tg://user?id=" + strconv.Itoa(update.Message.From.ID) + ")"
				adminRemove.ChatID = config.LogChat
				_, _ = tg.api.Send(adminRemove)
			} else {
				_, _ = tg.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You are using this command in the wrong way! Reply to a message of an admin to remove him from the role!"))
			}

			return
		}

		if update.Message.Text == "/ban" {
			// Check if the user is replying to a message
			if update.Message.ReplyToMessage != nil {
				err := tg.db.AddBan(int64(update.Message.ReplyToMessage.From.ID))
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR banning the user: "+err.Error()+".")
					msg.ReplyToMessageID = update.Message.MessageID
					_, _ = tg.api.Send(msg)
					return
				}

				// Success message
				banSuccess := tgbotapi.NewMessage(update.Message.Chat.ID, "User ["+update.Message.ReplyToMessage.From.FirstName+"](tg://user?id="+strconv.FormatInt(int64(update.Message.ReplyToMessage.From.ID), 10)+") is now banned\\!")
				banSuccess.ParseMode = tgbotapi.ModeMarkdownV2
				_, err = tg.api.Send(banSuccess)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR sending the success ban message: "+err.Error()+".")
					msg.ReplyToMessageID = update.Message.MessageID
					tg.api.Send(msg)
				}

				// Log action
				banSuccess.Text += "\n\nCommand issued by [" + update.Message.From.FirstName + "](tg://user?id=" + strconv.Itoa(update.Message.From.ID) + ")"
				banSuccess.ChatID = config.LogChat
				_, _ = tg.api.Send(banSuccess)
			} else {
				_, _ = tg.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You are using this command in the wrong way! Reply to a message to ban him from this bot!"))
			}

			return
		}

		if update.Message.Text == "/unban" {
			// Check if the user is replying to a message
			if update.Message.ReplyToMessage != nil {
				err := tg.db.RemoveAdmin(int64(update.Message.ReplyToMessage.From.ID))
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "ERROR: "+err.Error()+".")
					msg.ReplyToMessageID = update.Message.MessageID
					_, _ = tg.api.Send(msg)
					return
				}

				// Success message
				banRemove := tgbotapi.NewMessage(update.Message.Chat.ID, "User ["+update.Message.ReplyToMessage.From.FirstName+"](tg://user?id="+strconv.FormatInt(int64(update.Message.ReplyToMessage.From.ID), 10)+") is no longer banned\\!")
				banRemove.ParseMode = tgbotapi.ModeMarkdownV2
				_, _ = tg.api.Send(banRemove)

				// Log action
				banRemove.Text += "\n\nCommand issued by [" + update.Message.From.FirstName + "](tg://user?id=" + strconv.Itoa(update.Message.From.ID) + ")"
				banRemove.ChatID = config.LogChat
				_, _ = tg.api.Send(banRemove)
			} else {
				_, _ = tg.api.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "You are using this command in the wrong way! Reply to a message of a banned person to remove him from the ban list!"))
			}

			return
		}
	}

	// Commands for all
	if len(update.Message.Text) >= 5 && strings.ToLower(update.Message.Text[0:5]) == "/help" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "*HELP*\n\n*/help* \\- to use this command, view available commands\\.\n*/mizip <uid\\>* \\- Generate keys of a mizip from the UID\\.\n*/comestero <known key\\> <known key sector \\(0\\-15\\)\\> <known key type \\(A/B\\)\\>* \\- generate keys for a comestero vending key\\.\n\n*This bot has been created by [@GNUUnicorn](t.me/GNUUnicorn) and [@LilZ73](t.me/LilZ73)*\\, join @mikaiapp Telegram group\\.\n\nThanks to Golang for existing\\! 🦝")
		msg.ParseMode = tgbotapi.ModeMarkdownV2
		_, _ = tg.api.Send(msg)
		return
	}

	// Calculate the network infos
	if len(update.Message.Text) >= 5 && strings.ToLower(update.Message.Text[0:5]) == "/calc" {
		var args []string = strings.Split(update.Message.Text, " ")
		if len(args) >= 3 {
			netInfo := network.CalculateNetwork(args[1], args[2])

			netmask := network.ByteArrToStr(netInfo.Netmask.Dotted)
			wildcard := network.ByteArrToStr(netInfo.Wildcard)
			networkAddr := network.ByteArrToStr(netInfo.Network)
			broadcast := network.ByteArrToStr(netInfo.Broadcast)
			hostMinAddress := network.ByteArrToStr(netInfo.HostMinAddress)
			hostMaxAddress := network.ByteArrToStr(netInfo.HostMaxAddress)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID,fmt.Sprintf("Address: %s\nNetmask: %s\nWildcard: %s\nNetwork: %s\nBroadcast: %s\nHost Min Address: %s\nHost Max Address: %s\nHosts quantity: %d",args[1] + "/" + fmt.Sprint(netInfo.Netmask.Decimal), netmask, wildcard, networkAddr, broadcast, hostMinAddress, hostMaxAddress, netInfo.HostsQuantity))
			_, _ = tg.api.Send(msg)
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Invalid input uwu")
			msg.ParseMode = tgbotapi.ModeMarkdownV2
			_, _ = tg.api.Send(msg)
		}
		return
	}

	// Calculate the network infos and send a prettified output (might have a bad visualization for small devices)
	if len(update.Message.Text) >= 6 && strings.ToLower(update.Message.Text[0:6]) == "/pcalc" {

	}
}
