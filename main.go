package main

import (
	"context"
	"fmt"
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

func (p *Player) generateAnswer(answerChan chan<- Answer, round *Round) {
	answerNum := rand.Intn(RandLimit)
	fmt.Println(answerNum)
	answer := Answer{
		PlayerName:   p.Name,
		PlayerAnswer: answerNum,
		RightNumber:  round.RightNumber,
	}

	answerChan <- answer
}

type Answer struct {
	PlayerName   string
	PlayerAnswer int
	RightNumber  int
}

func main() {
	roundChan := make(chan *Round)
	answerChan := make(chan Answer)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, roundInterval)

	defer cancel()

	// Start the generator goroutine
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				round := &Round{}
				round.Start(roundChan)
				time.Sleep(roundInterval)
			}
		}
	}()

	for i := 1; i <= numPlayers; i++ {
		fmt.Println("Player run ", i)
		player := &Player{Name: fmt.Sprintf("Player-%d", i)}
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				case round := <-roundChan:
					player.generateAnswer(answerChan, round)
				}
			}
		}()
	}

	// Counter goroutine
	go func() {
		correctAnswers := 0
		for {
			select {
			case <-ctx.Done():
				return
			case answer := <-answerChan:
				if answer.PlayerAnswer == answer.RightNumber {
					fmt.Printf("\n%s made a right answer == %d", answer.PlayerName, answer.PlayerAnswer)
					correctAnswers++
				} else {
					fmt.Printf("\n%s made a wrong answer == %d", answer.PlayerName, answer.PlayerAnswer)
				}

				if correctAnswers == numPlayers {
					fmt.Println("\nAll players answered correctly!")
					cancel()
				}
			}
		}
	}()

	// Wait for the game to end
	<-ctx.Done()
	fmt.Println("\nGame Over!")
}
