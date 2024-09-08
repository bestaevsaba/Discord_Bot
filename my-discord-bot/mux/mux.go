package mux

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	

	"github.com/bwmarrin/discordgo"
)

type Route struct {
	Pattern     string      
	Description string     
	Help        string      
	Run         HandlerFunc 
}

type Context struct {
	Fields          []string
	Content         string
	IsDirected      bool
	IsPrivate       bool
	HasPrefix       bool
	HasMention      bool
	HasMentionFirst bool
}

type HandlerFunc func(*discordgo.Session, *discordgo.Message, *Context)

type Mux struct {
	Routes  []*Route
	Default *Route
	Prefix  string
}

func New() *Mux {
	m := &Mux{}
	m.Prefix = "-dg "
	return m
}

func (m *Mux) Route(pattern, desc string, cb HandlerFunc) (*Route, error) {

	r := Route{}
	r.Pattern = pattern
	r.Description = desc
	r.Run = cb
	m.Routes = append(m.Routes, &r)

	return &r, nil
}

func (m *Mux) FuzzyMatch(msg string) (*Route, []string) {

	fields := strings.Fields(msg)

	if len(fields) == 0 {
		return nil, nil
	}

	var r *Route
	var rank int

	var fk int
	for fk, fv := range fields {

		for _, rv := range m.Routes {

			if rv.Pattern == fv {
				return rv, fields[fk:]
			}

			if strings.HasPrefix(rv.Pattern, fv) {
				if len(fv) > rank {
					r = rv
					rank = len(fv)
				}
			}
		}
	}
	return r, fields[fk:]
}

var msg string
func (m *Mux) OnMessageCreate(ds *discordgo.Session, mc *discordgo.MessageCreate) {

	msg = ""

	var err error

	if mc.Author.ID == ds.State.User.ID {
		return
	}

	ctx := &Context{
		Content: strings.TrimSpace(mc.Content),
	}

	var c *discordgo.Channel
	c, err = ds.State.Channel(mc.ChannelID)
	if err != nil {
		c, err = ds.Channel(mc.ChannelID)
		if err != nil {
			log.Printf("unable to fetch Channel for Message, %s", err)
		} else {
			err = ds.State.ChannelAdd(c)
			if err != nil {
				log.Printf("error updating State with Channel, %s", err)
			}
		}
	}
	if c != nil {
		if c.Type == discordgo.ChannelTypeDM {
			ctx.IsPrivate, ctx.IsDirected = true, true
		}
	}

	if !ctx.IsDirected {
		for _, v := range mc.Mentions {

			if v.ID == ds.State.User.ID {

				ctx.IsDirected, ctx.HasMention = true, true

				reg := regexp.MustCompile(fmt.Sprintf("<@!?(%s)>", ds.State.User.ID))
				if reg.FindStringIndex(ctx.Content)[0] == 0 {
					ctx.HasMentionFirst = true
				}
				ctx.Content = reg.ReplaceAllString(ctx.Content, "")

				break
			}
		}
	}

	if !ctx.IsDirected && len(m.Prefix) > 0 {
		if strings.HasPrefix(ctx.Content, m.Prefix) {
			ctx.IsDirected, ctx.HasPrefix, ctx.HasMentionFirst = true, true, true
			ctx.Content = strings.TrimPrefix(ctx.Content, m.Prefix)
		}
	}

	if !ctx.IsDirected {
		msg = ctx.Content
		return
	}

	r, fl := m.FuzzyMatch(ctx.Content)
	if r != nil {
		ctx.Fields = fl
		r.Run(ds, mc.Message, ctx)
		return
	}

	if m.Default != nil && (ctx.HasMentionFirst) {
		m.Default.Run(ds, mc.Message, ctx)
	}
}

func GetUserMsg() string{
	for msg == ""{
		continue
	}
	return msg
}
