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
