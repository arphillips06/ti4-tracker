# Twilight Imperium Stats Tracker - Backend

This is the backend API service for the Twilight Imperium Stats Tracker. It manages player data, objectives, relics, scores, and game sessions for TI4 games.

## Features
- REST API to create games, assign players, score objectives
- Support for secret and public objective decks
- Relic logic (Shard of the Throne, Crown of Emphidia, The Obsidian)
- Agenda resolution and scoring
- Round progression and victory detection

## Tech Stack
- **Go** (Golang)
- **Gin** for routing and HTTP handling
- **GORM** for ORM/database layer
- **SQLite** (or optional PostgreSQL/MySQL support)

## Requirements
- Go >= 1.21

## Getting Started

```bash
# Clone the repo
cd backend
go run main.go
```
The server starts on localhost:8080

---
# Project Structure
```bash
├───controllers
├───database
│   ├───factions
│   └───objectives
├───models
└───services
```

---

## Contributing
Issues will be raised to keep track of work, please assign to yourself and then make a brach to work on it. Please PR back to main once done.
