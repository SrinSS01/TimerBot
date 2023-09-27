package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"regexp"
	"time"
)

type TimerCommand struct {
	Command *discordgo.ApplicationCommand
}

var Timer = TimerCommand{
	Command: &discordgo.ApplicationCommand{
		Name:        "timer",
		Description: "start a timer",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "user",
				Description: "user",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "start-date",
				Description: "start date",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "duration",
				Description: "duration in days",
				Required:    true,
			},
		},
	},
}

var regex = regexp.MustCompile("([1-9]|1[0-2])(?P<sep1>[/\\- ])([1-9]|[12][0-9]|3[01])(?P<sep2>[/\\- ])(20[2-9][3-9]|2[1-9][0-9][0-9])")

func (t *TimerCommand) Execute(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	data := interaction.ApplicationCommandData()
	user := data.Options[0].UserValue(session)
	startDateVal := data.Options[1].StringValue()
	durationDays := data.Options[2].IntValue()
	durationHours := time.Duration(durationDays) * time.Hour * 24
	matcher := regex.FindStringSubmatch(startDateVal)
	if len(matcher) == 0 {
		_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please enter a valid date format",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
	sep1 := matcher[regex.SubexpIndex("sep1")]
	sep2 := matcher[regex.SubexpIndex("sep2")]
	if sep1 != sep2 {
		if len(matcher) == 0 {
			_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Please enter a valid date format",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}
	}
	layout := fmt.Sprintf("1%s2%s2006", sep1, sep2)
	startDate, err := time.Parse(layout, startDateVal)
	if err != nil {
		_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Please enter a valid date\n```\n%s\n```", err.Error()),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}
	endDate := startDate.Add(durationHours)
	now := time.Now()
	startDiff := startDate.Sub(now)
	time.Sleep(startDiff - (time.Minute * 5))
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			_, _ = session.ChannelMessageSend(interaction.ChannelID, "@everyone")
		}
	}()
	time.Sleep(time.Minute * 5)
	ticker.Stop()
	_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:  "Username",
							Value: fmt.Sprintf("%s [ `%d` days ]", user.Mention(), durationDays),
						},
						{
							Name:  "Time left",
							Value: fmt.Sprintf("<t:%d:R>", endDate.Unix()),
						},
						{
							Name:  "Date Started",
							Value: fmt.Sprintf("<t:%d:d>", startDate.Unix()),
						},
						{
							Name:  "Date Finished",
							Value: fmt.Sprintf("<t:%d:d>", endDate.Unix()),
						},
					},
				},
			},
		},
	})
	//time.Sleep(durationHours)
	//ticker.Stop()
}
