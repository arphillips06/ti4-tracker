import { Link, useNavigate } from "react-router-dom";
import { useState } from "react";

export default function Home() {
  const navigate = useNavigate();
  const [gameId, setGameId] = useState("");

  const handleGoToGame = () => {
    if (gameId.trim()) {
      navigate(`/games/${gameId.trim()}`);
    }
  };

  return (
    <div className="container text-center mt-5">
      <h1 className="mb-4">Twilight Imperium Stats Tracker</h1>

      <img
        src="/ti-logo.png"
        alt="Twilight Imperium"
        className="img-fluid mb-4"
        style={{ maxHeight: "180px" }}
      />

      <div className="d-grid gap-3 col-6 mx-auto">
        <Link to="/new-game" className="btn btn-primary btn-lg">
          Start New Game
        </Link>
        <Link to="/games" className="btn btn-secondary btn-lg">
          View Past Games
        </Link>
      </div>

      {/* New section for quick game jump */}
      <div className="mt-5">
        <h4>Jump to Game</h4>
        <div className="input-group mb-3 justify-content-center" style={{ maxWidth: 400, margin: "0 auto" }}>
          <input
            type="text"
            className="form-control"
            placeholder="Enter Game ID"
            value={gameId}
            onChange={(e) => setGameId(e.target.value)}
          />
          <button className="btn btn-outline-success" onClick={handleGoToGame}>
            Go
          </button>
        </div>
      </div>
    </div>
  );
}
