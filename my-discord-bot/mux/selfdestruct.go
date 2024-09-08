package mux

import (
	"context"
	"fmt"
	"my-discord-bot/db"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *Mux) SelfDestruct(ds *discordgo.Session, i *discordgo.InteractionCreate) {
	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user db.User

	if err := db.UsersCollection.FindOneAndDelete(dbCtx, bson.M{"Id": i.Member.User.ID}).Decode(&user); err != nil {
		ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "User not found, try running the profit command to make a user account",
			},
		})
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Terminated",
		Color:       16711680,
		Description: fmt.Sprintf("%s's account has been terminated with %d units!", user.Name, user.Balance),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://cdn.pixabay.com/photo/2014/04/03/11/54/headstone-312540_960_720.png",
		},
	}

	ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
