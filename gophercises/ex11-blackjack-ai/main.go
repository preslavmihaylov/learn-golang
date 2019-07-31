package main

import (
	"fmt"

	"github.com/preslavmihaylov/learn-golang/gophercises/ex11-blackjack-ai/bot"
)

func main() {
	stats := bot.Simulate(100000, 6, 25, 4)
	totalHands := stats.HandsWon + stats.HandsTied + stats.HandsLost
	fmt.Printf("Final Balance: %d\n", stats.Balance)
	fmt.Printf("Hands Won: %d (%.2f%%)\n", stats.HandsWon, percentFromTotal(stats.HandsWon, totalHands))
	fmt.Printf("Hands Tied: %d (%.2f%%)\n", stats.HandsTied, percentFromTotal(stats.HandsTied, totalHands))
	fmt.Printf("Hands Lost: %d (%.2f%%)\n", stats.HandsLost, percentFromTotal(stats.HandsLost, totalHands))
	fmt.Printf("Win-Lose Ratio: %.2f%%-%.2f%%\n",
		percentFromTotal(stats.HandsWon, totalHands-stats.HandsTied),
		percentFromTotal(stats.HandsLost, totalHands-stats.HandsTied))
}

func percentFromTotal(hands, total int) float64 {
	return float64(hands) / float64(total) * 100
}
