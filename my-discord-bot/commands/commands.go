package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"my-discord-bot/mux"
	"os"

	"github.com/bwmarrin/discordgo"
)

var (
	Router    *mux.Mux
	stateFile = "commands/state.json" // Файл для хранения состояния команд
)

// Сохраняем текущее состояние команд в файл
func saveState(commands []*discordgo.ApplicationCommand) error {
	data, err := json.Marshal(commands)
	if err != nil {
		return fmt.Errorf("error marshaling commands state: %v", err)
	}
	return ioutil.WriteFile(stateFile, data, 0644)
}

// Загружаем состояние команд из файла
func loadState() ([]*discordgo.ApplicationCommand, error) {
	if _, err := os.Stat(stateFile); os.IsNotExist(err) {
		return nil, nil // Файл состояния не существует
	}

	data, err := ioutil.ReadFile(stateFile)
	if err != nil {
		return nil, fmt.Errorf("error reading state file: %v", err)
	}

	var commands []*discordgo.ApplicationCommand
	if err := json.Unmarshal(data, &commands); err != nil {
		return nil, fmt.Errorf("error unmarshaling state file: %v", err)
	}

	return commands, nil
}

// Проверяем, нужно ли перерегистрировать команды
func shouldRegisterCommands(newCommands []*discordgo.ApplicationCommand) bool {
	// Загружаем текущее состояние команд
	existingCommands, err := loadState()
	if err != nil || existingCommands == nil {
		return true // Состояние не загружено или не существует, нужно регистрировать
	}

	// Сравниваем существующие команды с новыми
	if len(existingCommands) != len(newCommands) {
		return true // Количество команд изменилось
	}

	for i, cmd := range existingCommands {
		if cmd.Name != newCommands[i].Name || cmd.Description != newCommands[i].Description {
			return true // Одна из команд изменилась
		}
	}

	return false // Команды не изменились
}

func RegisterCommands(s *discordgo.Session) error {
	commands, err := loadCommandsFromFile("commands/commands.json")
	if err != nil {
		return fmt.Errorf("error loading commands: %v", err)
	}

	if !shouldRegisterCommands(commands) {
		fmt.Println("Команды не изменились, повторная регистрация не требуется.")
		return nil
	}

	fmt.Println("Загруженные команды из JSON файла:")

	// Регистрация команд последовательно с проверкой ошибок
	for _, cmd := range commands {
		fmt.Printf("Регистрация команды: %s\n", cmd.Name)
		_, err := s.ApplicationCommandCreate(s.State.User.ID, "", cmd)
		if err != nil {
			fmt.Printf("Ошибка при регистрации команды %s: %v\n", cmd.Name, err)
			return fmt.Errorf("failed to create application command %s: %v", cmd.Name, err)
		}
	}

	// Сохраняем текущее состояние команд
	if err := saveState(commands); err != nil {
		fmt.Printf("Ошибка при сохранении состояния команд: %v\n", err)
	}

	fmt.Println("Все команды успешно зарегистрированы.")
	return nil
}

// Загружаем команды из JSON файла
func loadCommandsFromFile(filename string) ([]*discordgo.ApplicationCommand, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	var commands []*discordgo.ApplicationCommand
	if err := json.Unmarshal(byteValue, &commands); err != nil {
		return nil, fmt.Errorf("error unmarshaling json: %v", err)
	}

	return commands, nil
}

// Обрабатывает взаимодействия (команды) от пользователей
func OnInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		switch i.ApplicationCommandData().Name {
		case "ping":
			Router.Ping(s, i)
		case "embed":
			Router.Embed(s, i)
		case "hilo":
			Router.Hilo(s, i)
		case "tictactoe":
			Router.TicTacToe(s, i)
		case "roll":
			Router.Roll(s, i)
		case "bal":
			Router.Bal(s, i)
		case "profit":
			Router.Profit(s, i)
		case "selfdestruct":
			Router.SelfDestruct(s, i)
		case "helpme":
			Router.Help(s, i)
		}
	}
}

func DeleteAllCommands(s *discordgo.Session) error {
	commands, err := s.ApplicationCommands(s.State.User.ID, "")
	if err != nil {
		return fmt.Errorf("error fetching commands: %v", err)
	}
	for _, cmd := range commands {
		err := s.ApplicationCommandDelete(s.State.User.ID, "", cmd.ID)
		if err != nil {
			return fmt.Errorf("error deleting command %s: %v", cmd.Name, err)
		}
		fmt.Printf("Команда %s удалена\n", cmd.Name)
	}
	return nil
}
