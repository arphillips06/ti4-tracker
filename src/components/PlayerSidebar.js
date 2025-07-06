import React from "react";

export default function PlayerSidebar({
  playersSorted,
  expandedPlayers,
  setExpandedPlayers,
  game,
  scoreObjective,
  unscoreObjective,
  secretCounts,
  setSecretCounts,
  secretObjectives,
  objectiveScores,
  gameId,
  setGame,
  setObjectiveScores,
  refreshGameState,
  custodiansScored,

}) {
  const custodiansScorerId = game?.AllScores?.find((s) => s.Type === "mecatol")?.PlayerID || null;
  return (
    <div style={{ flex: "0 1 300px" }}>
      {(playersSorted || []).map((entry) => (
        <div
          key={entry.player_id}
          className="card mb-3 border-start border-5"
          style={{ borderColor: entry.color }}
        >
          <div className="card-body">
            <div className="d-flex justify-content-between align-items-center">
              <div className="fw-semibold small">{entry.name}</div>
              <button
                className="btn btn-sm btn-outline-secondary"
                onClick={() =>
                  setExpandedPlayers((prev) => ({
                    ...prev,
                    [entry.player_id]: !prev[entry.player_id],
                  }))
                }
              >
                {expandedPlayers[entry.player_id] ? "−" : "+"}
              </button>
            </div>
            <div className="d-flex align-items-center gap-2 mt-1">
              <img
                src={`/faction-icons/${entry.factionKey}.webp`}
                alt={entry.faction}
                style={{
                  width: "24px",
                  height: "24px",
                  borderRadius: "50%",
                  objectFit: "contain",
                  backgroundColor: "transparent",
                }}
                onError={(e) => (e.target.style.display = "none")}
              />
              {game?.AllScores?.some(
                (s) => s.Type === "mecatol" && s.PlayerID === entry.player_id
              ) && (
                  <img
                    src="/MR-point/MR-scored.png"
                    alt="Custodians Point"
                    title="Custodians Point"
                    style={{ width: "20px", height: "20px" }}
                  />
                )}
              <div className="text-muted small fst-italic">{entry.faction}</div>
            </div>
            <div className="mt-1 small">Points: {entry.points}</div>
            <div className="d-flex flex-column gap-2 mt-2">
              {expandedPlayers[entry.player_id] && (
                <div className="mt-3">
                  <div className="small fw-semibold mb-1">Secrets</div>
                  <div className="d-flex align-items-center gap-2 mb-2">
                    <div className="dropdown">
                      <button
                        className="btn btn-sm btn-outline-secondary dropdown-toggle"
                        type="button"
                        data-bs-toggle="dropdown"
                        aria-expanded="false"
                      >
                        −
                      </button>
                      <ul className="dropdown-menu">
                        {game.AllScores?.filter(
                          (s) =>
                            s.PlayerID === entry.player_id &&
                            (s.Type || "").toLowerCase() === "secret"
                        ).map((s) => {
                          const obj = secretObjectives.find(
                            (o) => o.id === s.ObjectiveID
                          );
                          console.log("🧪 AllScores:", game?.AllScores);
                          return obj ? (
                            <li key={obj.id}>
                              <button
                                className="dropdown-item"
                                onClick={async () => {
                                  try {
                                    await fetch("http://localhost:8080/score/mecatol", {
                                      method: "POST",
                                      headers: { "Content-Type": "application/json" },
                                      body: JSON.stringify({
                                        game_id: parseInt(gameId),
                                        player_id: entry.player_id,
                                      }),
                                    });
                                    console.log("✅ Fetching updated game...");

                                    await refreshGameState(); // ✅ Force a proper re-fetch of all state
                                    console.log("✅ Finished refreshing.");
                                  } catch (err) {
                                    console.error("Failed to score Custodians:", err);
                                    alert("Failed to score Custodians. See console.");
                                  }
                                }}                              >
                                {obj.name}
                              </button>
                            </li>
                          ) : null;
                        })}
                      </ul>
                    </div>

                    <div className="d-flex gap-1">
                      {[0, 1, 2].map((i) => {
                        const scored = (secretCounts[entry.player_id] || 0) > i;
                        return (
                          <img
                            key={i}
                            src={`/objective-backgrounds/secret-${scored ? "active" : "inactive"}.jpg`}
                            alt={scored ? "Scored secret" : "Unscored secret"}
                            style={{
                              width: "16px",
                              height: "25px",
                              opacity: scored ? 1 : 0.4,
                            }}
                          />
                        );
                      })}
                    </div>

                    <select
                      className="form-select form-select-sm"
                      value=""
                      onChange={async (e) => {
                        const selectedId = parseInt(e.target.value);
                        if (selectedId) {
                          const success = await scoreObjective(
                            entry.player_id,
                            selectedId,
                            entry.name
                          );
                          if (success) {
                            setSecretCounts((prev) => ({
                              ...prev,
                              [entry.player_id]: Math.min(
                                3,
                                (prev[entry.player_id] || 0) + 1
                              ),
                            }));
                          }
                        }
                      }}
                    >
                      <option value="">+ Score a secret objective</option>
                      {["Action", "Status", "Agenda"].map((phase) => (
                        <optgroup key={phase} label={phase}>
                          {secretObjectives
                            .filter((o) => o.phase === phase.toLowerCase())
                            .map((obj) => (
                              <option
                                key={obj.id}
                                value={obj.id}
                                disabled={objectiveScores[obj.id]?.includes(
                                  entry.name
                                )}
                              >
                                {obj.name}
                              </option>
                            ))}
                        </optgroup>
                      ))}
                    </select>
                  </div>
                </div>
              )}

              {game.AllScores?.some(
                (s) =>
                  s.PlayerID === entry.player_id &&
                  s.AgendaTitle === "Mutiny"
              ) && (
                  <div className="mt-1 small text-success">Bonus: Mutiny</div>
                )}
              {game.AllScores?.some(
                (s) => s.PlayerID === entry.player_id && s.AgendaTitle === "Seed of an Empire"
              ) && (
                  <div className="mt-1 small text-success">Bonus: Seed of an Empire</div>
                )}
              {expandedPlayers[entry.player_id] && (
                <div className="mt-3 small">
                  <button
                    className="btn btn-warning btn-sm"
                    disabled={game?.AllScores?.some((s) => s.Type === "mecatol")}
                    onClick={async () => {
                      try {
                        await fetch("http://localhost:8080/score/mecatol", {
                          method: "POST",
                          headers: { "Content-Type": "application/json" },
                          body: JSON.stringify({
                            game_id: parseInt(gameId),
                            player_id: entry.player_id,
                          }),
                        });

                        const updatedGame = await fetch(`http://localhost:8080/games/${gameId}`).then((r) => r.json());
                        updatedGame.AllScores = updatedGame.all_scores || [];
                        setGame(updatedGame);

                        const updatedScores = await fetch(
                          `http://localhost:8080/games/${gameId}/objectives/scores`
                        ).then((r) => r.json());
                        const map = {};
                        (Array.isArray(updatedScores)
                          ? updatedScores
                          : updatedScores?.value || []
                        ).forEach((entry) => {
                          map[entry.objective_id ?? entry.name] =
                            entry.scored_by || [];
                        });
                        setObjectiveScores(map);
                      } catch (err) {
                        console.error("Failed to score Custodians:", err);
                        alert("Failed to score Custodians. See console.");
                      }
                    }}
                  >
                    {game?.AllScores?.some((s) => s.Type === "mecatol")
                      ? "Custodians Already Scored"
                      : "Score Custodians"}
                  </button>
                </div>
              )}
            </div>
          </div>
        </div>
      ))}
    </div>
  );
}
