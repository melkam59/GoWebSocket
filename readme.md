# Simple Go WebSocket Implementation

A basic real-time chat application built to demonstrate the mental model and mechanics of WebSockets in Go.

## Features
- Real-time two-way communication using **Gorilla WebSockets**.
- Live "Who's Online" tracking of connected users (connects & disconnects).
- Message broadcasting to all concurrent clients using Go channels and Goroutines.
- Clean frontend rendered with the **Jet Template Engine**.

## How to Run
1. Run the application:
   ```bash
   go run ./cmd/web
   ```
2. Open a browser and visit `http://localhost:8080/`.
3. Open a second tab to the same URL, type a username in each, and start chatting to see the broadcast in action!

## Tech Stack
- **Go** (Standard HTTP library + Goroutines/Channels)
- **Gorilla WebSocket** (`github.com/gorilla/websocket`)
- **Jet Templates** (`github.com/CloudyKit/jet/v6`)
