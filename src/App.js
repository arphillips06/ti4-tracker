import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import 'bootstrap/dist/css/bootstrap.min.css';

import Home from './pages/Home';
import GameList from './pages/GameList';
import NewGamePage from './pages/NewGamePage'; // or './pages/GameForm' if that's what you name it
import GameDetail from './pages/GameDetail';


// import GameForm from './GameForm'; // Uncomment when this exists
// import GameDetail from './GameDetail'; // Uncomment when this exists

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/games" element={<GameList />} />
        <Route path="/new-game" element={<NewGamePage />} />
        <Route path="/games/:gameId" element={<GameDetail />} />

        {/* Add more routes as you build them */}
        {/* <Route path="/new-game" element={<GameForm />} /> */}
        {/* <Route path="/games/:id" element={<GameDetail />} /> */}
      </Routes>
    </Router>
  );
}

export default App;
