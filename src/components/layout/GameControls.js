export default function GameControls({
  game,
  scoringMode,
  setScoringMode,
  onAdvanceRound,
}) {
  if (!game) return null;
  return (
    <div className="d-flex justify-content-between align-items-center mb-4">
      <h2 className="h4">
        Round {game.current_round} | {game.winning_points} Point Game
      </h2>
      <div className="d-flex gap-3 align-items-center">
        <div className="form-check form-switch">
          <input
            className="form-check-input"
            type="checkbox"
            checked={scoringMode}
            onChange={(e) => setScoringMode(e.target.checked)}
          />
          <label className="form-check-label">Score Objectives</label>
        </div>
        <button className="btn btn-outline-primary btn-sm" onClick={onAdvanceRound}>
          Advance Round
        </button>
      </div>
    </div>
  );
}
