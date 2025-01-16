package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Token struct {
	ID         string
	UsageCount int
}

type TokenPool struct {
	tokens  []Token
	mutex   sync.Mutex
	randGen *rand.Rand
}

func NewTokenPool(numTokens int) *TokenPool {
	tokens := make([]Token, numTokens)
	for i := 0; i < numTokens; i++ {
		tokens[i] = Token{ID: fmt.Sprintf("Token %d", i+1), UsageCount: 0}
	}
	return &TokenPool{
		tokens:  tokens,
		randGen: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (tp *TokenPool) SelectLeastUsedToken() *Token {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()

	// Find the minimum usage count
	leastUsage := tp.tokens[0].UsageCount
	for _, token := range tp.tokens {
		if token.UsageCount < leastUsage {
			leastUsage = token.UsageCount
		}
	}

	// Collect all tokens with the least usage count
	leastUsedTokens := []*Token{}
	for i := range tp.tokens {
		if tp.tokens[i].UsageCount == leastUsage {
			leastUsedTokens = append(leastUsedTokens, &tp.tokens[i])
		}
	}

	randomIndex := tp.randGen.Intn(len(leastUsedTokens))
	return leastUsedTokens[randomIndex]
}

func (tp *TokenPool) IncrementUsage(token *Token) {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()
	token.UsageCount++
}

func (tp *TokenPool) ResetUsage() {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()
	for i := range tp.tokens {
		tp.tokens[i].UsageCount = 0
	}
}

func (tp *TokenPool) DisplayStats() {
	tp.mutex.Lock()
	defer tp.mutex.Unlock()

	leastUsage := tp.tokens[0].UsageCount
	leastUsedTokens := []Token{}

	fmt.Println("\nToken Usage Stats:")
	for _, token := range tp.tokens {
		fmt.Printf("%s: %d uses\n", token.ID, token.UsageCount)
		if token.UsageCount < leastUsage {
			leastUsage = token.UsageCount
			leastUsedTokens = []Token{token}
		} else if token.UsageCount == leastUsage {
			leastUsedTokens = append(leastUsedTokens, token)
		}
	}

	fmt.Printf("\nLeast Used Token(s):\n")
	for _, token := range leastUsedTokens {
		fmt.Printf("%s (%d uses)\n", token.ID, token.UsageCount)
	}
}

func simulateUserOperations(tp *TokenPool, userID int, numOperations int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < numOperations; i++ {
		token := tp.SelectLeastUsedToken()
		tp.IncrementUsage(token)
		fmt.Printf("User %d used %s\n", userID, token.ID)
		time.Sleep(10 * time.Millisecond)
	}
}

func main() {
	var numTokens, numUsers, numOperations int
	fmt.Print("Enter total number of tokens: ")
	fmt.Scanln(&numTokens)

	fmt.Print("Enter total number of users: ")
	fmt.Scanln(&numUsers)

	fmt.Print("Enter total number of operations per user: ")
	fmt.Scanln(&numOperations)

	// Initialize the token pool
	tokenPool := NewTokenPool(numTokens)

	var wg sync.WaitGroup
	for userID := 1; userID <= numUsers; userID++ {
		wg.Add(1)
		go simulateUserOperations(tokenPool, userID, numOperations, &wg)
	}

	// Wait for all users to complete their operations
	wg.Wait()

	// Simulate a 24-hour reset if needed
	fmt.Print("Simulate 24-hour reset? (yes/no): ")
	var resetInput string
	fmt.Scanln(&resetInput)
	if resetInput == "yes" {
		tokenPool.ResetUsage()
		fmt.Println("All token usages have been reset.")
	}

	// Output: Display stats
	fmt.Printf("\nSimulation Complete. Total Operations: %d\n", numUsers*numOperations)
	tokenPool.DisplayStats()
}
