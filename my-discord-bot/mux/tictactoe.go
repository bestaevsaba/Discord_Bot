package mux

import (
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var grid = [3][3]string{}
var currentPlayer string

func (m *Mux) TicTacToe(ds *discordgo.Session, i *discordgo.InteractionCreate) {
	gameEmbed := &discordgo.MessageEmbed{}
	gameEmbed.Title = "Tic Tac Toe"
	gameEmbed.Color = 11845097

	// Инициализация поля с символом H
	grid = [3][3]string{
		{"⠀", "⠀", "⠀"},
		{"⠀", "⠀", "⠀"},
		{"⠀", "⠀", "⠀"}}
	tilesLeft := 9

	// Первый игрок - X
	currentPlayer = "X"

	for !gameWon("X") && !gameWon("O") && tilesLeft > 0 {
		printGame(ds, i, gameEmbed)
		ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Игрок \"" + currentPlayer + "\", ваш ход.\nВведите ваш ход в формате 'ряд колонка' (Пример: 1 1)",
			},
		})
		userMoveStr := GetUserMsg()
		userMove := strings.Split(userMoveStr, " ")

		if len(userMove) != 2 {
			ds.ChannelMessageSend(i.ChannelID, "Неверное количество аргументов. Попробуйте снова.")
			continue
		}

		row, rowErr := strconv.Atoi(userMove[0])
		column, columnErr := strconv.Atoi(userMove[1])
		if rowErr != nil || columnErr != nil {
			ds.ChannelMessageSend(i.ChannelID, "Только числа принимаются. Попробуйте снова.")
			continue
		}

		row -= 1    // Сдвиг координат для внутреннего представления
		column -= 1 // Сдвиг координат для внутреннего представления
		if row < 0 || row > 2 || column < 0 || column > 2 {
			ds.ChannelMessageSend(i.ChannelID, "Введены значения вне диапазона. Попробуйте снова.")
			continue
		}

		if grid[row][column] != "⠀" {
			ds.ChannelMessageSend(i.ChannelID, "Эта клетка уже занята, выберите другую.")
			continue
		}

		grid[row][column] = currentPlayer
		tilesLeft -= 1

		// Проверка на победу
		if gameWon(currentPlayer) {
			printGame(ds, i, gameEmbed)
			ds.ChannelMessageSend(i.ChannelID, "Игрок \""+currentPlayer+"\" победил!")
			return
		}

		// Переключение на другого игрока
		if currentPlayer == "X" {
			currentPlayer = "O"
		} else {
			currentPlayer = "X"
		}
	}

	if tilesLeft <= 0 {
		printGame(ds, i, gameEmbed)
		ds.ChannelMessageSend(i.ChannelID, "Ничья!")
		return
	}
}

func printGame(ds *discordgo.Session, i *discordgo.InteractionCreate, game *discordgo.MessageEmbed) {
	game.Description = ""
	for _, row := range grid {
		// Между символами оставляем разделители
		game.Description += "| " + strings.Join(row[:], " | ") + " |\n\n"
	}
	ds.ChannelMessageSendEmbed(i.ChannelID, game)
}

func gameWon(val string) bool {
	for i := 0; i < 3; i++ {
		if grid[i][0] == val && grid[i][1] == val && grid[i][2] == val {
			return true
		}
	}
	for j := 0; j < 3; j++ {
		if grid[0][j] == val && grid[1][j] == val && grid[2][j] == val {
			return true
		}
	}
	if grid[1][1] == val {
		if grid[0][0] == val && grid[2][2] == val {
			return true
		}
		if grid[0][2] == val && grid[2][0] == val {
			return true
		}
	}
	return false
}
