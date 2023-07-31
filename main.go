package main

import (
	"context"
	"math/rand"
	"time"
)

const (
	roundInterval = 10 * time.Second
	numPlayers    = 5
	RandLimit     = 10
)

type Round struct {
	RightNumber int
}

func (r *Round) Start(roundChan chan *Round) {
	r.RightNumber = rand.Intn(RandLimit)
	roundChan <- r
}

type Player struct {
	Name string
}

func (p *Player) generateAnswer(answerChan chan Answer) {
	answer := Answer{
		PlayerName:   p.Name,
		PlayerAnswer: rand.Intn(RandLimit),
	}

	answerChan <- answer
}

type Answer struct {
	PlayerName   string
	PlayerAnswer int
}

func main() {
	roundChan := make(chan *Round)
	answerChan := make(chan Answer)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, roundInterval)

}
