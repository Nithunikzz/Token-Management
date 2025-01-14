# Token Manager

A simple Go program for managing a set of tokens. It selects the least-used token and tracks its usage. The usage counts are reset every 24 hours.

## Problem Statement

The goal is to implement a token manager system that:
- Manages 1000 tokens.
- Tracks the usage count of each token.
- Selects tokens with the least usage for each operation.
- Resets the usage count for all tokens after 24 hours.

## Approach

- **Token Structure:** Each token has an ID and a `Usage` count.
- **Token Manager:** The `TokenManager` holds the list of tokens and ensures thread safety with a mutex.
- **Token Selection:** The `SelectToken` function selects a token with the least usage. If multiple tokens have the same least usage, one of them is selected randomly.
- **Usage Reset:** Every 24 hours, the usage count of all tokens is reset to zero using a timer.
- **Simulation:** The program allows users to simulate operations, where each operation selects a token and increments its usage.

