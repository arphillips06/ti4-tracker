# TI4 Score Tracker — Frontend

A web-based interface for tracking player scores, factions, and objectives in **Twilight Imperium: Fourth Edition**.

This is the **frontend** portion of the project, built with React and Vite. It connects to a Go-based backend.

---

## Technologies

- React (via Vite)
- Bootstrap (for styling)
- React Router
- Fetch API (to communicate with the Go backend)

---

### Getting Started

#### 1. Clone the repository
#### 2. Install Node.js
#### 3. Install project dependancies
```
npm install
```
This installs everything listed in package.json using npm.
#### 4. Start the development server
```
npm start
```

---
## Project Structure
```
src/
├── components/         # Reusable UI components (e.g. PlayerRow)
├── pages/              # Page views (Home, GameList, NewGamePage)
├── App.js              # Main application routes
└── main.jsx            # Application entry point
```

