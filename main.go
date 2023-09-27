package main

import (
	"TimerBot/commands"
	"TimerBot/config"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	cnfg    = config.Config{}
	discord *discordgo.Session
	cmds    = []*discordgo.ApplicationCommand{
		commands.Timer.Command,
	}
	commandHandlers = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
		commands.Timer.Command.Name: commands.Timer.Execute,
	}
)

func init() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Print("Enter bot token: ")
		if _, err := fmt.Scanln(&cnfg.Token); err != nil {
			log.Fatal("Error during Scanln(): ", err)
		}
		configJson()
		return
	}
	if err := json.Unmarshal(file, &cnfg); err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}
}

func configJson() {
	marshal, err := json.Marshal(&cnfg)
	if err != nil {
		log.Fatal("Error during Marshal(): ", err)
		return
	}
	if err := os.WriteFile("config.json", marshal, 0644); err != nil {
		log.Fatal("Error during WriteFile(): ", err)
	}
}

func onReady(session *discordgo.Session, _ *discordgo.Ready) {
	log.Println(session.State.User.Username + " is ready")
}

func slashCommandInteraction(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
	if interaction.Type != discordgo.InteractionApplicationCommand {
		return
	}
	commandHandlers[interaction.ApplicationCommandData().Name](session, interaction)
}

func main() {
	var err error
	discord, err = discordgo.New("Bot " + cnfg.Token)
	if err != nil {
		log.Fatal("Error creating Discord session", err)
		return
	}
	discord.AddHandler(onReady)
	discord.AddHandler(slashCommandInteraction)
	if err := discord.Open(); err != nil {
		log.Fatal("Error opening connection", err)
		return
	}
	for _, command := range cmds {
		_, err := discord.ApplicationCommandCreate(discord.State.User.ID, "", command)
		if err != nil {
			log.Fatal("Error creating slash command", err)
			return
		}
	}
	log.Println("Bot is running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
	if err := discord.Close(); err != nil {
		log.Fatal("Error closing connection", err)
		return
	}
	log.Println("Bot is shutting down")
}
