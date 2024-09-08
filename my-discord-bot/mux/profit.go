package mux

import (
	"context"
	"fmt"
	"my-discord-bot/db"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *Mux) Profit(ds *discordgo.Session, i *discordgo.InteractionCreate) {
	var user db.User
	dbCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.UsersCollection.FindOne(dbCtx, bson.M{"Id": i.Member.User.ID}).Decode(&user); err != nil {
		userResult, err := db.UsersCollection.InsertOne(dbCtx, bson.D{
			{Key: "Id", Value: i.Member.User.ID},
			{Key: "Name", Value: i.Member.User.Username},
			{Key: "Balance", Value: 0},
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(userResult.InsertedID)
	}

	err := db.UsersCollection.FindOneAndUpdate(
		dbCtx,
		bson.M{"Id": i.Member.User.ID},
		bson.D{
			{Key: "$set", Value: bson.D{{Key: "Balance", Value: user.Balance + 100}}},
		},
	).Decode(&user)
	if err != nil {
		fmt.Println(err)
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Profit",
		Color:       42622,
		Description: fmt.Sprintf("%s earned 100 units!", user.Name),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://e7.pngegg.com/pngimages/450/717/png-clipart-dollar-dollar.png",
		},
	}

	ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}
