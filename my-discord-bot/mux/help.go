package mux

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (m *Mux) Help(ds *discordgo.Session, i *discordgo.InteractionCreate) {
	var helpMessage strings.Builder
	helpMessage.WriteString("**Available Commands:**\n")

	for _, route := range m.Routes {
		helpMessage.WriteString(fmt.Sprintf("**%s**: %s\n", route.Pattern, route.Description))
	}

	ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: helpMessage.String(),
		},
	})
}
