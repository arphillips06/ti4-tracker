# Database

This package manages initialization and loading of core static game data (factions and objectives).

## Files

- `setup.go` – Initializes the in-memory or file-based data storage used by the app.

## Subfolders

- `factions/` – Contains code for loading and managing Twilight Imperium factions.
- `objectives/` – Contains definitions for Stage I, Stage II, and Secret objectives.

The database layer is read-only and used to populate game data at startup.
