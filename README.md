<<<<<<< HEAD
# TI4 Tracker

This is the monorepo for 2 personal projects. 

The app is designed to track a full game of Twilight Imperium 4 and then serve the stats of all games entered. This was made as there is nothing made by the community that does this. 

## TI4 Backend
This is my Go project that facilitates the input from the frontend puts the input into a sqlite database, which is then served back. 
As a learning project I have tried to keep to DRY princables and Go best practicies. 
This was made by using the skills in Go I have learnt at my role in work and an interest in coding. 

## TI4 Frontend
This is a website using JS. Serves the game funtionality as well as stats pages. This was made primarily using AI as I know no JS.
=======
# Twilight Imperium Stats Tracker – Backend

This is the Go backend for the Twilight Imperium stats tracker app. It handles player/game management, scoring logic, and complex rules around objectives, agendas, and relics.

---

## Tech Stack

- **Language**: Go
- **Framework**: net/http
- **Database**: SQLite (via GORM)
- **Architecture**: MVC
  - `/controllers`: API endpoints
  - `/services`: Business logic
  - `/models`: Data structs
  - `/database`: Static data files (objectives, factions, etc.)

---

## Getting Started

### Prerequisites

- Go 1.21+
- SQLite installed

### Installation

```bash
git clone https://github.com/your-username/ti4-tracker.git
cd ti4-tracker/backend
go mod tidy
```

### Running Locally

```bash
go run main.go
```

Server will run at `http://localhost:8080` by default.

---

## API Overview

### Players

- `POST /game/:id/player` — Add player to a game
- `GET /game/:id` — Get game details including player scores

### Scoring

- `POST /score` — Score public/secret objective
- `POST /score/mecatol` — Score Custodians token
- `POST /score/imperial` — Score Imperial point
- `POST /score/agenda` — Score or lose points from an agenda
- `POST /score/relic` — Handle relics like Crown or Shard

### Game Management

- `POST /game` — Create a new game
- `POST /game/:id/round` — Advance to the next round

---

## Key Concepts

### Objectives

Stored in `/database/objectives.json`. Include:
- ID
- Name
- Type (`public1`, `public2`, `secret`)
- Phase (`action`, `status`)

### Scoring

Each score is saved with:
- `PlayerID`, `GameID`, `Points`, `SourceType`, `SourceID`
- Special sources: `agenda`, `mecatol`, `imperial`, `relic`

### Relics

Currently supported:
- **Crown of Emphidia** – 1 point to a selected player
- **Shard of the Throne** – Transfers point when holder changes
- **The Obsidian** – Increases secret objective limit

### Agendas

Agenda scoring allows for positive and negative points. Some agendas (e.g. **Seed of an Empire**) create new objectives. Others (e.g. **Mutiny**) just grant points.

---

## Folder Structure

```
backend/
├── controllers/
├── services/
├── models/
├── database/
│   ├── factions.json
│   ├── objectives.json
│   └── ...
├── main.go
├── helpers/
```

---

## Contributing
Issues will be raised to keep track of work, please assign to yourself and then make a brach to work on it. Please PR back to main once done.
=======
# Twilight Imperium Stats Tracker - Frontend

This is the frontend for the Twilight Imperium Stats Tracker. It provides a responsive web UI to manage, track, and visualize player progress, objectives, relics, and agendas during a game of Twilight Imperium 4th Edition.

## Features
- Game setup: players, factions, and points
- Scoring interface for objectives (public and secret)
- Agenda and relic effects (e.g., Mutiny, Shard, Crown, Obsidian)
- Dynamic scoring and real-time UI updates
- Player sidebar for secret tracking and expansion controls
- Round advancement and scoring lock toggle
- Graphs and history features (WIP)

## Tech Stack
- **React** with **Vite** 
- **CSS** for styling
- **Bootstrap** components for layout and dropdowns
- **React Icons** for iconography
- **Fetch API** for backend communication

## Requirements
- Node.js >= 18
- Backend API (Go) running on `localhost:8080` 

## Getting Started

```bash
# Clone the repo
cd frontend
npm install
npm run dev
```
---

## Project Structure
```bash
src/
├── components/         # Reusable UI components (e.g. PlayerSidebar, ObjectivesGrid)
├── hooks/              # Custom React hooks (e.g. useGameData)
├── pages/              # Top-level views (GameDetail, Home)
├── relics/             # Handlers for relic-related logic
├── App.js              # Main app with routing
├── index.js            # Entry point
└── config.js           # Backend API base URL
```
---

## Access
Access is granted via Tailscale until v6 can be made to work with the backend.

To access the tailnet, request access - link will be provided. From there login via a third party (git or google for example), once done download the application to $device and then you should be able to join the tailnet from there and access the magic URL ```ti4-stats.tail798251.ts.net```
>>>>>>> 7dc2949 (Squashed 'frontend/' content from commit cb858be)
