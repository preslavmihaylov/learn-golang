# Exercise #10: Blackjack Game

[![exercise status: released](https://img.shields.io/badge/exercise%20status-released-green.svg?style=for-the-badge)](https://gophercises.com/exercises/blackjack)

## Exercise details

**Note:** This exercise will use the completed code from the [Deck of Cards Exercise](https://gophercises.com/exercises/blackjack), so if you haven't already completed it you will either need to grab the finished source code (currently this is in the [p8 branch](https://github.com/gophercises/deck/tree/p8)), or you will need to complete the exercise on your own.

Now that we have a deck of cards to work with, we want to build something using the cards. Blackjack is a relatively simple card game with a simple AI, so that will be the focus of this exercise. We won't be building out all of the rules, but we will implement the major ones and you are welcome to expand upon the game once you have it working.

Specifically, we are going to create a blackjack game that supports the following rules:

#### 1. Every player is dealt 2 cards

For simplicity, our first version will only support two players - the dealer and the human player.

The dealing starts at the first player, and continues around the table until the dealer is dealt and then repeats, starting again with the first player, until all players have two cards.

The dealer only has one visible card. The other is "face down" and isn't visible to players. All player cards are visible.

#### 2. The player's turn

In our limited version of blackjack, the player will have two options: Hit or Stand.

If a player chooses to hit, they are dealt a new card and will then be allowed to choose between the hit and stand options again.

If a player chooses to stand their turn ends and the next player is up.

#### 3. The dealer's turn

In the first iteration our dealer won't do anything, and will just display their hand. After that the game will end.

In our second iteration the dealer will play with typical dealer rules - if they have a score of 16 or less, or a soft 17, they will hit. This means we will need to implement scoring, and will be able to determine which player has won the game.

#### 4. Determining the winner

The winner is the player who has the highest score without going over 21. Cards 2-10 are worth their face value in points. Aces are worth 1 or 11 points (whichever is best without busting), and face cards (J, Q, K) are all worth 10 points.

A "soft 17" is a score of 17 in which 11 of the points come from an Ace card.

If the player busts during their turn, the dealer automatically wins. 

Blackjack occurs when a player has an Ace with a face card (J, Q, K), or a 10 card. In traditional blackjack there are special rules for this, but for our simple game we won't be adding that.


## Bonus

This exercise is ripe with bonus options. You can:

1. Expand upon the game, adding in new rules like allowing the player or dealer to win immediately if they have a natural blackjack.
2. Add a way for players to bet, and to keep track of their wins and losses.
3. Add in additional rules like splitting and doubling down.
4. Add support for additional human players, or for a single player to play multiple hands against the computer.
5. Add AI for the "human" players with varying degrees of intelligence. Eg the dealer has a specific strategy regardless of what players have visible, but a player's optimal strategy will vary based on the cards the dealer has visible. Try different strategies out or look up some and implement them.
6. Lastly, you could go a step further adding in some support for card counting and varying bets based on the state of the deck. See if you can come up with a strategy that reliably wins money over time given a limited cash pool.
