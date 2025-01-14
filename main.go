package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Token struct {
	ID    string
	Usage int
}

type TokenManager struct {
	Tokens     []Token
	Mutex      sync.Mutex
	ResetTimer *time.Timer
}

// Initialize TokenManager with 1000 tokens
func NewTokenManager() *TokenManager {
	tokens := make([]Token, 1000)
	for i := 0; i < 1000; i++ {
		tokens[i] = Token{
			ID:    fmt.Sprintf("Token %d", i+1),
			Usage: 0,
		}
	}

	tm := &TokenManager{
		Tokens: tokens,
	}

	tm.ResetTimer = time.AfterFunc(24*time.Hour, tm.ResetUsage)
	return tm
}

// Select a token with the least usage count
func (tm *TokenManager) SelectToken() *Token {
	tm.Mutex.Lock()
	defer tm.Mutex.Unlock()

	// Find the token(s) with the least usage count
	minUsage := tm.Tokens[0].Usage
	indices := []int{0}

	for i, token := range tm.Tokens {
		if token.Usage < minUsage {
			minUsage = token.Usage
			indices = []int{i}
		} else if token.Usage == minUsage {
			indices = append(indices, i)
		}
	}

	// Randomly select one token from the least-used tokens
	selectedIndex := indices[rand.Intn(len(indices))]
	tm.Tokens[selectedIndex].Usage++

	// Print the selected token and its usage
	fmt.Printf("Selected Token: %s, Current Usage: %d\n", tm.Tokens[selectedIndex].ID, tm.Tokens[selectedIndex].Usage)
	return &tm.Tokens[selectedIndex]
}

// Reset the usage counts of all tokens
func (tm *TokenManager) ResetUsage() {
	tm.Mutex.Lock()
	defer tm.Mutex.Unlock()

	for i := range tm.Tokens {
		tm.Tokens[i].Usage = 0
	}

	tm.ResetTimer.Reset(24 * time.Hour)

	fmt.Println("Token usage counts have been reset to zero.")
}

// Simulate token usage
func (tm *TokenManager) Simulate(operations int) {
	for i := 0; i < operations; i++ {
		tm.SelectToken()
	}
}

func (tm *TokenManager) DisplayResults() {
	leastUsage := tm.Tokens[0].Usage
	leastUsedTokens := []Token{}

	for _, token := range tm.Tokens {
		fmt.Printf("%s: %d uses\n", token.ID, token.Usage)
		// Track the least used token(s)
		if token.Usage < leastUsage {
			leastUsage = token.Usage
			leastUsedTokens = []Token{token} // Reset the list with the new least-used token
		} else if token.Usage == leastUsage {
			leastUsedTokens = append(leastUsedTokens, token)
		}
	}

	// Display the least used token(s)
	fmt.Println("\nLeast Used Token(s):")
	for _, token := range leastUsedTokens {
		fmt.Printf("%s (%d uses)\n", token.ID, token.Usage)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var operations int
	fmt.Print("Enter simulation time (operations): ")
	fmt.Scan(&operations)

	// Initialize the token manager
	tm := NewTokenManager()

	// Simulate the token usage for the specified number of operations
	tm.Simulate(operations)

	// Display the results after the simulation
	fmt.Printf("\nSimulation Time: %d operations\n\n", operations)
	tm.DisplayResults()

	//tm.ResetUsage()
}
