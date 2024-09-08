package mux

import (
	"context"
	"fmt"
	"my-discord-bot/db"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *Mux) Bal(ds *discordgo.Session, i *discordgo.InteractionCreate) {
	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user db.User

	if err := db.UsersCollection.FindOne(dbCtx, bson.M{"Id": i.Member.User.ID}).Decode(&user); err != nil {
		ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "User not found, try running the profit command to make a user account",
			},
		})
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Balance",
		Color:       0,
		Description: fmt.Sprintf("%s has %d units!", user.Name, user.Balance),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://www.pngitem.com/pimgs/m/101-1016890_icon-bank-logo-png-transparent-png.png",
		},
	}

	ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
