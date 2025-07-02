import { Link } from "react-router-dom";

export default function Home() {
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
    </div>
  );
}
