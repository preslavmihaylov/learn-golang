# Exercise #11: Blackjack AI

[![exercise status: released](https://img.shields.io/badge/exercise%20status-released-green.svg?style=for-the-badge)](https://gophercises.com/exercises/blackjack_ai)

## Exercise details

When we were completing the [blackjack](https://gophercises.com/exercises/blackjack) exercise we were focused on creating a command line game to play blackjack. As a result, we made a lot of decisions that made sense in that context but might not necessarily make sense in other contexts. For instance, we had a lot of exported fields, types, and functions that in a normal package wouldn't be exported, but for a `main` package it didn't really matter. 

In this exercise we are going to explore how we might refactor some code - like the code used to create our blackjack game - into a package that might be used by other developers. In doing this, we will need to ask ourselves which functions, types, and other data types should be exported and documented vs which should be unexported and kept as simply implementation details. 

While we are refactoring our code we are also going to adjust how our game is played; rather than always expecting the player to be a person typing input into the console, we are instead going to expose an interface that can be implemented in order to play the blackjack game. Something like this:

```go
type AI interface {
  // Define functions needed for an AI to play your blackjack game.
}
```

*NOTE: I'm not going to tell you what functions to add to the `AI` type just yet, but I'm going to try to create a single video that covers these and add it to the start of the series so you can check it out if you are looking for ideas without wanting to spoil the entire exercise.*

Explaining what all this entails is tricky without spoiling the exercise, but in the end you will have two packages - `main` and `blackjack` - where `main` contains the AI implementation and will start the game, while `blackjack` contains the `AI` interface as well as all the (mostly) hidden logic to actually run the blackjack game based on decisions made by the AI. In short, you want a `main` package that looks kinda like this:

```go
package main

import (
  "fmt"

  "github.com/gophercises/blackjack_ai/blackjack"
  // along with other imports
)

// This type should implement the blackjack.AI interface
type AI struct {}

func main() {
  var ai AI
  // setup ai if you need to...

  opts := blackjack.Options{
    Hands: 100,
    Decks: 3,
  }
  game := blackjack.New(opts)
  winnings := game.Play(ai)
  fmt.Println("Our AI won/lost:", winnings)
}
```

Your implementation details may vary. You may name things differently, call your functions by different names, or whatever else you see fit. The primary point is that in the end our `blackjack` package won't export more than it needs to, it will power playing the actual game, and we will be writing an AI to play the game.

While refactoring your code and writing the `blackjack` package, you will also need to add support for a few more parts of the blackjack game. For starters, you will definitely need betting so that we can keep track of how our AI works. You can start off by just making all hands worth 1 and keeping track of wins/losses, but having a way to bet would be ideal for the final version.

You may also want to add in more blackjack options like splitting, higher payouts on blackjack, and doubling down on a hand. Supporting splitting can be especially tricky if you support an arbitrary number of splits, so keep that in mind as you code it up and feel free to skip it until the end if you wish.

## Hints (if you get stuck)

I'll eventually add some hints and tips here, but for now just email me - <jon@calhoun.io> - if you get stuck and let me know what is confusing you and I'll try to help out then update this section based on your questions 😀

