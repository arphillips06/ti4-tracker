import { Link, useNavigate } from "react-router-dom";
import { useState } from "react";

export default function Home() {
  const navigate = useNavigate();
  const [gameId, setGameId] = useState("");

const handleGoToGame = () => {
  const trimmed = gameId.trim();
  if (!trimmed || !/^\d+$/.test(trimmed)) return;

  navigate(`/games/${trimmed}`);
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
        <Link to="/games" className="btn btn-primary btn-lg">
          View Past Games
        </Link>
        <Link to="/stats" className="btn btn-primary btn-lg">
          View Game Stats
        </Link>
      </div>
      {/* New section for quick game jump */}
      <div className="mt-5">
        <form
          onSubmit={(e) => {
            e.preventDefault();
            handleGoToGame();
          }}
          style={{ textAlign: "center", marginTop: "2rem" }}
        >
          <h4 style={{ color: "#e0c873", marginBottom: "0.5rem" }}>Jump to Game</h4>

          <div style={{ display: "flex", justifyContent: "center" }}>
            <div
              style={{
                display: "flex",
                backgroundColor: "transparent",
                borderRadius: "6px",
                overflow: "hidden",
                border: "2px solid #e0c873",
              }}
            >
              <input
                type="text"
                placeholder="Enter Game ID"
                inputMode="numeric"
                pattern="[0-9]*"
                value={gameId}
                onChange={(e) => {
                  const val = e.target.value;
                  // Only allow digits
                  if (/^\d*$/.test(val)) {
                    setGameId(val);
                  }
                }}
                style={{
                  padding: "10px 14px",
                  backgroundColor: "#111122",
                  color: "#e0c873",
                  border: "none",
                  fontSize: "16px",
                  fontFamily: "Merriweather, serif",
                  width: "250px",
                  outline: "none",
                }}
              />
              <button
                type="submit"
                className="btn-primary"
                style={{
                  border: "none",
                  backgroundColor: "transparent",
                  fontFamily: "Merriweather, serif",
                }}
              >
                GO
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
}
