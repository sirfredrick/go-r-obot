package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
)

var sent bool = false
var content string = ""
var messageCount int = 0

type mockSession struct {
	err error
}

func (mockSession) ChannelMessageSend(id string, c string) (*discordgo.Message, error) {
	m := &discordgo.Message{}
	var err error = nil
	sent = true
	content = c
	messageCount++
	return m, err
}

func (mockSession) AddHandler(handler interface{}) func() {
	return test
}

func test() {
}

func (s mockSession) Open() error {
	var err error = s.err
	return err
}

func (mockSession) Close() error {
	var err error = nil
	return err
}

func TestRun(t *testing.T) {
	ms := mockSession{err: nil}
	var err error = os.ErrNotExist
	run(ms, err, true)

	ms = mockSession{}
	err = nil
	run(ms, err, true)

	ms = mockSession{err: os.ErrNotExist}
	err = nil
	run(ms, err, true)
}

func TestSendMessage(t *testing.T) {
	sent = false
	u := discordgo.User{
		ID: "1",
	}
	m := discordgo.Message{
		Author:  &u,
		Content: "r/obot",
	}
	mc := discordgo.MessageCreate{
		Message: &m,
	}
	s := mockSession{}
	sendMessage(s, &mc, "1")
	if sent {
		t.Errorf("Expected bot to not send a message as a response to itself, instead it sent a message.")
	}
	sent = false
	numbers.setup()
	numbers.shuffle()
	sendMessage(s, &mc, "2")
	if !sent {
		t.Errorf("Expected bot to send message in response to an 'r/' but no message was sent.")
	}
	for i := range numbers {
		numbers[i] = 38
	}
	content = ""
	messageCount = 0
	sendMessage(s, &mc, "2")
	if !strings.Contains(content, "Taylor defused the bomb!") {
		t.Errorf("Expected last message on bomb to be 'Taylor defused the bomb!' but got, %v", content)
	}
	if messageCount != 8 {
		t.Errorf("Expected 8 messages to be sent on bomb, but got %v", messageCount)
	}
}

func TestSetup(t *testing.T) {
	var ns numberSlice = make([]int, 0)
	ns.setup()
	if len(ns) != 39 {
		t.Errorf("Expected Number Slice length of 39, but got %v", len(ns))
	}
	if ns[0] != 0 {
		t.Errorf("Expected first number in Number Slice to be 0, but got %v", ns[0])
	}
	if ns[len(ns)-1] != 38 {
		t.Errorf("Expected last number in Number Slice to be 38, but got %v", ns[len(ns)-1])
	}
}

func TestShuffle(t *testing.T) {
	var ns numberSlice = make([]int, 0)
	ns.setup()
	var sns numberSlice = make([]int, 0)
	sns.setup()
	sns.shuffle()
	shuffled := false
	for i, n := range sns {
		if n != ns[i] {
			shuffled = true
		}
	}
	if !shuffled {
		t.Errorf("Expected shuffled array, but got %v", sns)
	}
}

func TestReadWriteCount(t *testing.T) {
	os.Remove("_countTest.txt")
	c := 0
	for i := 0; i < 5; i++ {
		c++
	}
	err := writeCount("_countTest.txt", c)
	if err != nil {
		t.Errorf("Expected file to be written but got this Error: %v", err)
	}
	rc := readCount("_countTest.txt")
	if rc != 5 {
		t.Errorf("Expected the count read from the file to be 5, but got %v", rc)
	}
	os.Remove("_countTest.txt")
}

func TestReadCountError(t *testing.T) {
	os.Remove("_countError.txt")
	readCount("_countError.txt")
	bs, err := ioutil.ReadFile("_countError.txt")
	if errors.Is(err, os.ErrNotExist) {
		t.Errorf("Expected readCount to create file on error, but got file doesn't exist")
	}
	c, err := strconv.Atoi(string(bs))
	if err != nil {
		t.Errorf("Expected readCount to add number to file but got, Error: %v", err)
	}
	if c != 0 {
		t.Errorf("Expected readCount to write 0 to file on error, but got %v", c)
	}
	os.Remove("_countError.txt")

	ioutil.WriteFile("_countError.txt", []byte("a"), 0666)
	c = readCount("_countError.txt")
	if c != 0 {
		t.Errorf("Expected readCount to return 0 on error, but got %v", c)
	}
	os.Remove("_countError.txt")

	ioutil.WriteFile("_countError.txt", []byte(fmt.Sprint(1)), 0000)
	c = readCount("_countError.txt")
	if c != 0 {
		t.Errorf("Expected readCount to return 0 when file has no permissions, but got %v", c)
	}

	os.Remove("_countError.txt")
}
