package useragent

import (
	_ "embed"
	"math/rand"
	"strings"
	"time"
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

var (
	//go:embed user-agents.txt
	userAgents     string
	userAgentsList []string
)

func init() {
	userAgentsList = strings.Split(userAgents, "\n")
}

func GetRandomUserAgent() string {
	return userAgentsList[random.Intn(len(userAgentsList))]
}
