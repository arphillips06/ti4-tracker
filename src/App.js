import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import 'bootstrap/dist/css/bootstrap.min.css';

import Home from './components/Home';
import GameList from './GameList';
// import GameForm from './GameForm'; // Uncomment when this exists
// import GameDetail from './GameDetail'; // Uncomment when this exists

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/games" element={<GameList />} />
        {/* Add more routes as you build them */}
        {/* <Route path="/new-game" element={<GameForm />} /> */}
        {/* <Route path="/games/:id" element={<GameDetail />} /> */}
      </Routes>
    </Router>
  );
}

export default App;
