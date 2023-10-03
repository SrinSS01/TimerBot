package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"regexp"
	"strings"
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
	durationHours := time.Duration(durationDays) * 24 * time.Hour
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
	_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   discordgo.MessageFlagsEphemeral,
			Content: "⏱️ Scheduled",
		},
	})
	now := time.Now()
	endDate := startDate.Add(durationHours)
	startDiff := startDate.Sub(now)
	time.AfterFunc(startDiff, func() {
		embed := discordgo.MessageEmbed{
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Username",
					Value: fmt.Sprintf("%s [ `%d` days ]", user.Mention(), durationDays),
				},
				{
					Name:  "Time Passed",
					Value: fmt.Sprintf("%s", time.Now().Sub(startDate)),
				},
				{
					Name:  "Time Left",
					Value: fmt.Sprintf("%s", endDate.Sub(time.Now())),
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
		}
		msg, _ := session.ChannelMessageSendEmbed(interaction.ChannelID, &embed)
		ticker := time.NewTicker(1 * time.Second)
		time.AfterFunc(durationHours, func() {
			ticker.Stop()
		})
		go func() {
			for range ticker.C {
				passed := time.Now().Sub(startDate)
				passedHrs := int(passed.Hours())
				passedMins := int(passed.Minutes())
				passedSecs := int(passed.Seconds())
				builder := strings.Builder{}
				if passedHrs > 0 {
					builder.WriteString(fmt.Sprintf("`%d hrs`", passedHrs))
				}
				if passedMins > 0 {
					builder.WriteString(fmt.Sprintf("`%d mins`", passedMins))
				}
				builder.WriteString(fmt.Sprintf("`%d secs`", passedSecs))
				embed.Fields[1] = &discordgo.MessageEmbedField{
					Name:  "Time Passed",
					Value: builder.String(),
				}
				leftBuilder := strings.Builder{}
				left := endDate.Sub(time.Now())
				leftHrs := int(left.Hours())
				leftMins := int(left.Minutes())
				leftSecs := int(left.Seconds())
				if leftHrs > 0 {
					leftBuilder.WriteString(fmt.Sprintf("`%d hrs`", leftHrs))
				}
				if leftMins > 0 {
					leftBuilder.WriteString(fmt.Sprintf("`%d mins`", leftMins))
				}
				leftBuilder.WriteString(fmt.Sprintf("`%d secs`", leftSecs))
				embed.Fields[2] = &discordgo.MessageEmbedField{
					Name:  "Time Left",
					Value: builder.String(),
				}
				_, _ = session.ChannelMessageEditEmbed(msg.ChannelID, msg.ID, &embed)
			}
		}()
	})
	go func() {
		announceTime := endDate.Sub(now) - (5 * time.Minute)
		time.AfterFunc(announceTime, func() {
			ticker := time.NewTicker(1 * time.Second)
			go func() {
				for range ticker.C {
					_, _ = session.ChannelMessageSend(interaction.ChannelID, "@everyone")
				}
			}()
			time.AfterFunc(5*time.Second, func() {
				ticker.Stop()
				time.Sleep(1 * time.Second)
				_, _ = session.ChannelMessageSendEmbed(interaction.ChannelID, &discordgo.MessageEmbed{
					Description: fmt.Sprintf("%s got released", user.Mention()),
				})
			})
		})
	}()
}
