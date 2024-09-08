package mux

import (
	"math/rand/v2"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func (m *Mux) Hilo(ds *discordgo.Session, i *discordgo.InteractionCreate) {
	number := rand.IntN(101)
	ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Guess a number between 0 and 100",
		},
	})

	guess := getGuess(ds, i)

	for guess != number {
		if guess < number {
			ds.ChannelMessageSend(i.ChannelID, "Too low!")
		}
		if guess > number {
			ds.ChannelMessageSend(i.ChannelID, "Too high!")
		}
		guess = getGuess(ds, i)
	}
	ds.ChannelMessageSend(i.ChannelID, "Correct! The number was "+strconv.Itoa(number))
}

func getGuess(ds *discordgo.Session, i *discordgo.InteractionCreate) int {
	returnGuess, err := strconv.Atoi(GetUserMsg())
	for err != nil {
		ds.ChannelMessageSend(i.ChannelID, "That is not a number, guess again")
		returnGuess, err = strconv.Atoi(GetUserMsg())
	}
	return returnGuess
}
