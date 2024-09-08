package mux

import (
	"math/rand/v2"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func (m *Mux) Roll(ds *discordgo.Session, i *discordgo.InteractionCreate) {
	number := rand.IntN(101)
	response := strconv.Itoa(number)

	ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}
