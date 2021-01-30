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

package main

import (
	"fmt"
	"go-Telegram-Network-bot/config"
	"go-Telegram-Network-bot/roles"
	"go-Telegram-Network-bot/telegram"
	"runtime"
)

func main() {
	rolesDb, err := roles.NewRoles(config.RolesFile)

	if err != nil {
		fmt.Println(err)
		panic("Unable to start roles.")
	}

	// Configure all parameters and run goroutines
	networkBot, err := telegram.NewTelegramBot(config.Token, rolesDb)

	if err != nil {
		rolesDb.Close()
		fmt.Println(err)
		panic("Unable to configure Telegram bot from token.")
	}

	if err = networkBot.ManageUpdates(networkBot.HandleUpdate); err != nil {
		rolesDb.Close()
		fmt.Println(err)
		panic("Unable to start Telegram polling routine.")
	}

	// Terminate main goroutine but keep running the others
	runtime.Goexit()
}
