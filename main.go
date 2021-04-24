package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type numberSlice []int

type session interface {
	ChannelMessageSend(channelID string, content string) (*discordgo.Message, error)
	AddHandler(handler interface{}) func()
	Open() error
	Close() error
}

var messages []string = []string{"You've brought up Reddit %v times. Don't you think that's a bit much?", "Users have been an intellectual %v times", "%v SubReddits mentioned.", "This server's love of Reddit is at least %v", "r/%vsubreddits", "Talk all you want, because SubReddit #%v will be your last!", "%v SubReddits. You went too far, you fool.", "%v scenarios saved.", "The princess is now a permanent guest at one of my %v Koopa hotels.", "Lamp oil? Rope? Bombs? You want it, it's yours, my friend, as long as you have %v rubies.", "I found %v pecha berries. Me too! I found a healing potion, but it won't be enough!", "Don't let a single one get away! We'll each need to take down about %v. Stow your fear. It's now or never!", "Jacob now has %v less knees.", "Back in my heyday, 5 hits would have been enough to knock that guy out. But today it took more than %v mighty blows.", "IQ's %v, but you can leave it all to me!", "C%v18", "%v reasons why I'm here", "Heh heh heh heh ha ha ha. You fool. I have %v alternative accounts!", "Ha ha, ha. %v!", "Fallout %v Battle Royale", "Error %v, could not generate new error message", "OwO, you have UwU'd exactwy %v times todaw my fwendo", "1 + 1 = %v", "Space, the final frontier. These are the voyages of the startship Enterprise. Its continuing mission: to explore strange new worlds, to seek out new life and new civilizations, to boldy go where %v people have gone before.", "It has been %v femtoseconds since the last bot related accident.", "Hey guys, have you prepared for Area %v Raid yet?", "The %v Precepts of Zote", "I am %v parallel dimensions ahead of you!", "There's %v days of summer vacation...", "Sir, the possibility of successfully navigating an asteroid field is approximately %v to 1.", "You have been given %v patridges in a pear tree.", "Wait, are you just doing this to see if there are any more messages? Because that is one of the worst reasons to call this bot %v times.", "You have procrastinated on going to your grandfathers funeral by clicking %v objects.", "Welp, that's %v more titles that Ethan has.", "Error. %v is less than infinity.", "Still not convinced? Okay, okay, I'll cut you a deal... The game's available for %v dollars, and that's a great price.", "There are %v lights!", "Somehow, you have rolled a %v ... what was your modifier again?", "Self destruct sequence initiated. You have %v seconds until detonation."}
var numbers numberSlice = make([]int, 0)
var count int = 0
var index int = 0

func main() {
	bs, err := ioutil.ReadFile(".secret")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	dg, err := discordgo.New("Bot " + strings.Trim(string(bs), "\n"))
	run(dg, err, false)
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	sendMessage(s, m, s.State.User.ID)
}

func sendMessage(s session, m *discordgo.MessageCreate, botID string) {
	if m.Author.ID == botID {
		return
	}
	regexp1 := regexp.MustCompile(`(^[^a-zA-Z0-9_]*)(r​\/[a-zA-Z0-9_]+)|([^a-zA-Z0-9_]+)(r​\/[a-zA-Z0-9_]+)\w+`)
	regexp2 := regexp.MustCompile(`(^[^a-zA-Z0-9_]*)(r\/[a-zA-Z0-9_]+)|([^a-zA-Z0-9_]+)(r\/[a-zA-Z0-9_]+)\w+`)

	if regexp1.Find([]byte(m.Content)) != nil || regexp2.Find([]byte(m.Content)) != nil {
		count++
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(messages[numbers[index]], count))
		if numbers[index] == 38 {
			for i := count - 1; i > count-6; i-- {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprint(i))
			}
			s.ChannelMessageSend(m.ChannelID, "...")
			s.ChannelMessageSend(m.ChannelID, "Taylor defused the bomb!")
		}
		index++
		if index >= len(numbers) {
			index = 0
			numbers.shuffle()
		}
	}
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	t := 0
	a := discordgo.Activity{
		Name: "for Subreddits",
		Type: discordgo.ActivityTypeListening,
		URL: "",
	}
	as := make([]*discordgo.Activity, 1)
	as[0] = &a
	afk := false
	stat := "online"
	usd := discordgo.UpdateStatusData{
		IdleSince: &t,
		Activities:      as,
		AFK:       afk,
		Status:    stat,
	}
	s.UpdateStatusComplex(usd)
}

func run(dg session, err error, debug bool) {
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}
	count = readCount("data.txt")
	numbers.setup()
	numbers.shuffle()
	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening connection: ", err)
		return
	}
	fmt.Println("Bot is now running, Press Ctrl-C to exit.")
	if !debug {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
		writeCount("data.txt", count)
	}
	dg.Close()
}

func (ns *numberSlice) setup() {
	for i := range messages {
		*ns = append(*ns, i)
	}
}

func (ns numberSlice) shuffle() {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	for i := range ns {
		j := r.Intn(len(ns) - 1)
		ns[i], ns[j] = ns[j], ns[i]
	}
}

func readCount(p string) int {
	bs, err := ioutil.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			writeCount(p, 0)
			return 0
		}
		fmt.Println("Error: ", err)
		return 0
	}
	c, err := strconv.Atoi(string(bs))
	if err != nil {
		fmt.Println("Error: ", err)
	}
	return c
}

func writeCount(p string, c int) error {
	return ioutil.WriteFile(p, []byte(fmt.Sprint(c)), 0666)
}
