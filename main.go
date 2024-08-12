package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// Player struct to store the details of each player
type Player struct {
	Name         string // Full name of the player
	InitialMoney int    // Initial money the player starts with
	Balance      int    // Current balance after all games
	Wins         int    // Number of games the player has won
	TotalSap     int    // Number of 'Three of a Kind' wins (Sáp)
	TenPointWins int    // Number of 10-point wins
	TotalMoney   int    // Total money won by the player
}

// Function to initialize and return a deck of cards
func initDeck() []int {
	// Define card values
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 10, 10, 10}
	// Create a deck with 52 cards
	deck := make([]int, 0, 52)
	for i := 0; i < 4; i++ {
		deck = append(deck, values...)
	}
	return deck
}

// Function to deal cards and calculate points and check for 'Sáp' (Three of a Kind)
func dealAndCalculatePoints(deck []int) (int, bool) {
	// Shuffle the deck
	rand.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	// Draw 3 cards
	hand := deck[:3]
	// Calculate points as the sum of the three cards modulo 10
	points := (hand[0] + hand[1] + hand[2]) % 10
	// Check if all three cards are the same ('Sáp')
	isSap := hand[0] == hand[1] && hand[1] == hand[2]
	return points, isSap
}

// Function to input a number with a prompt
func inputNumber(prompt string) int {
	var num int
	fmt.Println(prompt)
	fmt.Scanln(&num)
	return num
}

// Function to input the names of all players
func inputPlayer(num int) []string {
	players := make([]string, num)
	reader := bufio.NewReader(os.Stdin)
	// Loop to input names for each player
	for i := 0; i < num; i++ {
		fmt.Printf("Nhập tên đầy đủ của player %d: ", i+1)
		name, _ := reader.ReadString('\n')   // Read full name including spaces
		players[i] = strings.TrimSpace(name) // Remove any leading/trailing spaces
	}
	return players
}

// Function to input the initial money for each player
func inputInitialMoney(players []string) map[string]int {
	initialMoney := make(map[string]int)
	for _, player := range players {
		initialMoney[player] = inputNumber(fmt.Sprintf("Nhập số vốn của %s: ", player))
	}
	return initialMoney
}

// Function to input the bet amount for each game
func inputBetAmount() int {
	return inputNumber("Nhập số tiền đặt cược cho mỗi ván: ")
}

// Main function to run the card game
func main() {
	rand.Seed(time.Now().UnixNano())                       // Seed the random number generator
	numPlayer := inputNumber("Nhập số lượng người chơi: ") // Input number of players
	numGames := inputNumber("Nhập số ván chơi: ")          // Input number of games
	players := inputPlayer(numPlayer)                      // Input names of players
	initialMoney := inputInitialMoney(players)             // Input initial money for each player
	betAmount := inputBetAmount()                          // Input bet amount for each game
	playerStats := make(map[string]*Player)                // Initialize player statistics

	// Initialize player stats with their name, initial money, and balance
	for _, name := range players {
		playerStats[name] = &Player{
			Name:         name,
			InitialMoney: initialMoney[name],
			Balance:      initialMoney[name],
		}
	}

	payments := make(map[string]int) // Initialize a map to track payments

	// Loop through each game
	for i := 0; i < numGames; i++ {
		deck := initDeck()             // Initialize the deck
		scores := make(map[string]int) // Initialize scores map
		isSap := make(map[string]bool) // Initialize 'Sáp' map
		maxPoints := -1                // Initialize max points
		winner := ""                   // Initialize winner name

		// Deal cards and calculate points for each player
		for _, player := range players {
			points, sap := dealAndCalculatePoints(deck)
			scores[player] = points
			isSap[player] = sap
			if points > maxPoints {
				maxPoints = points
				winner = player
			}
		}

		// Calculate the prize amount based on the results
		prize := betAmount * (len(players) - 1)

		// Double the prize for a 10-point win
		if scores[winner] == 10 {
			prize *= 2
			playerStats[winner].TenPointWins++
		}

		// Triple the prize for a 'Sáp' (Three of a Kind)
		if isSap[winner] {
			prize *= 3
			playerStats[winner].TotalSap++
		}

		// Update winner statistics
		playerStats[winner].Wins++
		playerStats[winner].TotalMoney += prize
		playerStats[winner].Balance += prize

		// Calculate payments (who pays whom)
		for _, player := range players {
			if player != winner {
				payments[player] -= prize / (len(players) - 1)
				playerStats[player].Balance -= prize / (len(players) - 1)
				payments[winner] += prize / (len(players) - 1)
			}
		}
	}

	// Print player statistics
	for _, player := range players {
		stats := playerStats[player]
		fmt.Printf("%s: Vốn ban đầu %dk, thắng %d ván, %d ván với Sáp, %d ván với 10 điểm, và hiện có %dk trong ví\n",
			stats.Name, stats.InitialMoney, stats.Wins, stats.TotalSap, stats.TenPointWins, stats.Balance)
	}

	fmt.Println("\nDanh sách thanh toán:")

	// Determine who profited and who lost
	profits := make([]string, 0)
	losses := make([]string, 0)

	// Track who needs to pay and who will receive money
	for _, player := range players {
		if payments[player] > 0 {
			profits = append(profits, player)
			fmt.Printf("%s nhận %dk\n", player, payments[player])
		} else if payments[player] < 0 {
			losses = append(losses, player)
			fmt.Printf("%s trả %dk\n", player, -payments[player])
		}
	}

	// Print detailed transactions (who pays whom)
	fmt.Println("\nChi tiết thanh toán:")

	for _, loser := range losses {
		for _, winner := range profits {
			if payments[loser] == 0 {
				break
			}

			// Determine the amount to be paid
			paymentAmount := min(-payments[loser], payments[winner])
			payments[loser] += paymentAmount
			payments[winner] -= paymentAmount

			fmt.Printf("%s trả %dk cho %s\n", loser, paymentAmount, winner)
		}
	}
}

// Helper function to find the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
