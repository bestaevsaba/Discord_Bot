package mux

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (m *Mux) Embed(ds *discordgo.Session, i *discordgo.InteractionCreate) {
	// Упоминание пользователя в обычном сообщении
	mention := fmt.Sprintf("<@%s>", i.Member.User.ID)

	// Создание и отправка embed сообщения
	embed := &discordgo.MessageEmbed{
		Title:       "This is an Embed",
		Color:       1752220,
		Timestamp:   time.Now().Format(time.RFC3339),
		Description: "This is Description " + mention, // Упоминание в Description
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQx35EZEMKjJsxFAyh8ZUi4eDRmHmyZbJqKGw&s",
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "@by Maxaraja", // Footer не поддерживает упоминания напрямую
			IconURL: "https://i.imgur.com/5KWSvGO.jpeg",
		},
	}

	err := ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
	if err != nil {
		fmt.Println("Error sending embed message:", err)
	}
}
