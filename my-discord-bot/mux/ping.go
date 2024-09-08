package mux

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (m *Mux) Ping(ds *discordgo.Session, i *discordgo.InteractionCreate) {
	// Получение времени создания взаимодействия
	timestamp, err := discordgo.SnowflakeTimestamp(i.ID)
	if err != nil {
		fmt.Println("Ошибка при получении времени создания:", err)
		return
	}

	// Вычисление времени задержки
	response := "Pong! " + time.Since(timestamp).String()

	// Ответ на слэш-команду
	err = ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
	if err != nil {
		fmt.Println("Ошибка при отправке ответа:", err)
	}
}
