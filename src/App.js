import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import 'bootstrap/dist/css/bootstrap.min.css';
import 'bootstrap/dist/js/bootstrap.bundle.min.js';
import StatsOverview from './pages/StatsOverview';
import './pages/stats.css';

import Home from './pages/Home';
import GameList from './pages/GameList';
import NewGamePage from './pages/NewGamePage';
import GameDetail from './pages/GameDetail';
import Achievements from './pages/Achievements';

function App() {
  const backgroundStyle = {
    backgroundImage: "url('/space-bg.png')",
    backgroundSize: "cover",
    backgroundPosition: "center center",
    backgroundRepeat: "no-repeat",
    backgroundAttachment: "fixed",
    minHeight: "100vh",
    width: "100%",
  };

  return (
    <div style={backgroundStyle}>
      
      <Router>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/games" element={<GameList />} />
          <Route path="/new-game" element={<NewGamePage />} />
          <Route path="/games/:gameId" element={<GameDetail />} />
          <Route path="/stats" element={<StatsOverview />} />
          <Route path="/achievements" element={<Achievements />} />
        </Routes>
      </Router>
    </div>
  );
}

export default App;
