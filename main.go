package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/stefankopieczek/gossip/log"
)

var (
	// Caller parameters
	caller = &EndPoint{
		DisplayName: "Alice Phone",
		UserName:    "alice",
		Host:        "127.0.0.1",
		Port:        5062,
		Transport:   "UDP",
	}

	// Callee parameters
	callee = &EndPoint{
		DisplayName: "Bob Phone",
		UserName:    "bob",
		Host:        "127.0.0.1",
		Port:        5060,
		Transport:   "UDP",
	}
)

func main() {
	log.SetDefaultLogLevel(log.WARN)
	err := caller.Start()
	if err != nil {
		log.Warn("Failed to start caller: %v", err)
		return
	}
	err = callee.Start()
	if err != nil {
		log.Warn("Failed to start caller: %v", err)
		return
	}

	go forever()
	select {}
}

func forever() {
	for {
		time.Sleep(25 * time.Millisecond)
		go callee.ServeInvite()
		go caller.Invite(callee)
		go callee.Bye(caller)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

const branchMagicCookie = "z9hG4bK"

// generateRandom generates random strings for tags/call-ids/branches
func generateRandom(charsToDouble int) string {
	var buf bytes.Buffer
	for i := 0; i < charsToDouble; i++ {
		buf.WriteByte(byte(randInt(65, 90)))
		buf.WriteByte(byte(randInt(97, 122)))
	}

	return string(buf.Bytes())
}

// Helper for generating random number between two given numbers
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

// GenerateBranch generates branch for via header
func GenerateBranch() string {
	randomPart := generateRandom(4)
	return branchMagicCookie + randomPart
}

// GenerateTag generates tags for To/From headers
func GenerateTag() string {
	return generateRandom(10)
}

// GenerateCallID generates call-id
func GenerateCallID() string {
	return generateRandom(20)
}

// CheckConnError is used to check for errors encountered during connection establishment
func CheckConnError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
