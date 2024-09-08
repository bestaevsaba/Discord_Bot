package main

import (
	"fmt"
	"my-discord-bot/commands"
	"my-discord-bot/config"
	"my-discord-bot/db"
	"my-discord-bot/mux"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	bot *discordgo.Session
)

func main() {
	fmt.Println("Бот запущен!")

	var wg sync.WaitGroup
	wg.Add(1) // Подключение к БД

	// Подключение к базе данных
	go func() {
		defer wg.Done()
		db.Connect()
	}()

	commands.Router = &mux.Mux{} // Инициализация Router

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	// Открываем сессию бота
	err := bot.Open()
	if err != nil {
		fmt.Println("Ошибка при запуске бота:", err)
		return
	}
	defer bot.Close()

	// Удаление всех существующих команд (если нужно)
	// err = commands.DeleteAllCommands(bot)
	// if err != nil {
	// 	fmt.Println("Ошибка при удалении команд:", err)
	// 	return
	// }

	// Регистрируем команды при необходимости
	err = commands.RegisterCommands(bot)
	if err != nil {
		fmt.Println("Ошибка при регистрации команд:", err)
	}

	<-sc

	db.Disconnect() // Закрытие соединения с БД
}

func init() {
	err := config.ReadConfig()
	if err != nil {
		fmt.Println("Ошибка чтения конфигурации:", err)
		return
	}

	bot, err = discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("Ошибка при инициализации бота:", err)
		return
	}

	bot.AddHandler(commands.OnInteractionCreate)
}
