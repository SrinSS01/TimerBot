package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
)

type TimerCommand struct {
	Command *discordgo.ApplicationCommand
}

var (
	DayMinimum    = 1.0
	MinuteMinimum = 0.0
	HourMinimum   = 0.0
	YearMinimum   = float64(time.Now().Year())
)
var Timer = TimerCommand{
	Command: &discordgo.ApplicationCommand{
		Name:        "timer",
		Description: "start a timer",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "user",
				Description: "user",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "start-day",
				Description: "start day",
				Required:    true,
				MinValue:    &DayMinimum,
				MaxValue:    31.0,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "start-month",
				Description: "start month",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "January",
						Value: "Jan",
					},
					{
						Name:  "February",
						Value: "Feb",
					},
					{
						Name:  "March",
						Value: "Mar",
					},
					{
						Name:  "April",
						Value: "Apr",
					},
					{
						Name:  "May",
						Value: "May",
					},
					{
						Name:  "June",
						Value: "Jun",
					},
					{
						Name:  "July",
						Value: "Jul",
					},
					{
						Name:  "August",
						Value: "Aug",
					},
					{
						Name:  "September",
						Value: "Sep",
					},
					{
						Name:  "October",
						Value: "Oct",
					},
					{
						Name:  "November",
						Value: "Nov",
					},
					{
						Name:  "December",
						Value: "Dec",
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "start-year",
				Description: "start year",
				Required:    true,
				MinValue:    &YearMinimum,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "start-hour",
				Description: "start hour",
				Required:    true,
				MinValue:    &HourMinimum,
				MaxValue:    12,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "start-minute",
				Description: "start minute",
				Required:    true,
				MinValue:    &MinuteMinimum,
				MaxValue:    59,
			},
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "time",
				Description: "AM/PM",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "AM",
						Value: "AM",
					},
					{
						Name:  "PM",
						Value: "PM",
					},
				},
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

func (t *TimerCommand) Execute(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	data := interaction.ApplicationCommandData()
	user := data.Options[0].StringValue()
	day := data.Options[1].IntValue()
	month := data.Options[2].StringValue()
	year := data.Options[3].IntValue()
	hours := data.Options[4].IntValue()
	minutes := data.Options[5].IntValue()
	AmPm := data.Options[6].StringValue()
	durationDays := data.Options[7].IntValue()
	durationHours := time.Duration(durationDays) * 24 * time.Hour

	/*matcher := regex.FindStringSubmatch(startDateVal)
	if len(matcher) == 0 {
		_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Please enter a valid date format",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}*/
	/*sep1 := matcher[regex.SubexpIndex("sep1")]
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
	}*/
	minStr := ""
	if minutes < 10 {
		minStr = fmt.Sprintf("0%d", minutes)
	} else {
		minStr = fmt.Sprintf("%d", minutes)
	}
	startDateVal := fmt.Sprintf("%s %d, %d at %d:%s%s (AST)", month, day, year, hours, minStr, AmPm)
	const layout = "Jan 2, 2006 at 3:04PM (MST)"
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
	now := time.Now()
	if startDate.Compare(now) < 0 {
		_ = session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("The time provided must be greater than the current time <t:%s>.", now.Unix()),
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
	endDate := startDate.Add(durationHours)
	startDiff := startDate.Sub(now)
	time.AfterFunc(startDiff, func() {
		now := time.Now()
		d, h, m, s := getTime(now.Sub(startDate))
		tld, tlh, tlm, tls := getTime(endDate.Sub(now))
		embed := discordgo.MessageEmbed{
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:  "Username",
					Value: fmt.Sprintf("%s [ `%d` days ]", user, durationDays),
				},
				{
					Name:  "Time Passed",
					Value: fmt.Sprintf("`%d days` `%d hrs` `%d mins` `%d secs`", d, h, m, s),
				},
				{
					Name:  "Time Left",
					Value: fmt.Sprintf("`%d days` `%d hrs` `%d mins` `%d secs`", tld, tlh, tlm, tls),
				},
				{
					Name:  "Start Day",
					Value: fmt.Sprintf("<t:%d>", startDate.Unix()),
				},
				{
					Name:  "End Day",
					Value: fmt.Sprintf("<t:%d>", endDate.Unix()),
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
				now := time.Now()
				passed := now.Sub(startDate)
				pd, ph, pm, ps := getTime(passed)
				builder := strings.Builder{}
				if pd > 0 {
					builder.WriteString(fmt.Sprintf("`%d days` ", pd))
				}
				if ph > 0 {
					builder.WriteString(fmt.Sprintf("`%d hrs` ", ph))
				}
				if pm > 0 {
					builder.WriteString(fmt.Sprintf("`%d mins` ", pm))
				}
				builder.WriteString(fmt.Sprintf("`%d secs`", ps))
				embed.Fields[1] = &discordgo.MessageEmbedField{
					Name:  "Time Passed",
					Value: builder.String(),
				}
				leftBuilder := strings.Builder{}
				left := endDate.Sub(now)
				ld, lh, lm, ls := getTime(left)
				if ld > 0 {
					leftBuilder.WriteString(fmt.Sprintf("`%d days` ", ld))
				}
				if lh > 0 {
					leftBuilder.WriteString(fmt.Sprintf("`%d hrs` ", lh))
				}
				if lm > 0 {
					leftBuilder.WriteString(fmt.Sprintf("`%d mins` ", lm))
				}
				leftBuilder.WriteString(fmt.Sprintf("`%d secs`", ls))
				embed.Fields[2] = &discordgo.MessageEmbedField{
					Name:  "Time Left",
					Value: leftBuilder.String(),
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
					Description: fmt.Sprintf("%s got released", user),
				})
			})
		})
	}()
}

func getTime(duration time.Duration) (int, int, int, int) {
	return int(duration.Hours() / 24), int(duration.Hours()) % 24, int(duration.Minutes()) % 60, int(duration.Seconds()) % 60
}
