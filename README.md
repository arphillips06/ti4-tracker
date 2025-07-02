# Twilight Imperium Stats Tracker

A local self-hosted web application for tracking stats and scoring in **Twilight Imperium 4** games. Built with Go, this app lets you assign players, select factions, track scores round-by-round, and analyze overall performance across games.

---

## Features

- Assign players to games and select from official factions  
- Track scoring by public and secret objectives  
- Monitor round-by-round progress and game outcomes  
- Analyze stats over multiple games *(planned)*  
- Predefined database of official objectives and factions  
- RESTful API for frontend integration  

---

## Getting Started

### Prerequisites

- Go 1.18 or later  

### Clone and Run

```
bash
git clone https://github.com/arphillips06/twilight-imperium-tracker.git
cd twilight-imperium-tracker
go run main.go
```
---
### Commands to use
These were tested from powershell. The localhost will open the port 8080, ```127.0.0.1``` is set explicty at the moment.

```
//make users - all users MUST be present before making a game
$player = @{ Name = "Alice" } | ConvertTo-Json
Invoke-RestMethod -Method POST -Uri http://localhost:8080/players -ContentType "application/json" -Body $player

// create a game. Users MUST be created first and factions MUST match what is in /database/factions/factions.go
$game = @{
    WinningPoints = 10
    Players = @(
        @{ Name = "Alice"; Faction = "Arborec" },
        @{ Name = "Bob"; Faction = "Argent Flight" }
        @{ Name = 'Charlie'; Faction = "Barony of Letnev"}
    )
} | ConvertTo-Json -Depth 10

Invoke-RestMethod -Method POST -Uri http://localhost:8080/games -ContentType "application/json" -Body $game

//Advance the round and draw a new objective
Invoke-RestMethod -Method POST -Uri http://localhost:8080/games/1/advance-round

//show ALL objectives selected for the game, at the moment all will be shown regardless of the revealed state. The 'unrevealed' cards won't have a roundID attached.
Invoke-RestMethod -Method GET -Uri http://localhost:8080/games/1/objectives

//get all players in a game
Invoke-RestMethod -Uri -Method GET http://localhost:8080/games/1/players
```
Please note that were there is an ```int``` in the URL this is the ID of the game. If any more games are made then that number should be updated to reflect the new request.

---
### Contributing

Contributions can be made by pulling the codebase into a branch and then making a PR to be reviewed. 
