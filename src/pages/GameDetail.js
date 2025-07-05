import { useParams } from 'react-router-dom';
import { useEffect, useState } from 'react';

export default function GameDetail() {
  const { gameId } = useParams();
  const [game, setGame] = useState(null);
  const [objectives, setObjectives] = useState([]);
  const [objectiveScores, setObjectiveScores] = useState({});
  const [localScored, setLocalScored] = useState({});
  const [scoringMode, setScoringMode] = useState(false);
  const [expandedPlayers, setExpandedPlayers] = useState({});
  const custodiansScored = game?.AllScores?.some((s) => s.Type === "mecatol") || false;
console.log("Custodians scored:", custodiansScored);


  useEffect(() => {
    fetch(`http://localhost:8080/games/${gameId}`)
      .then((res) => res.json())
      .then(setGame)
      .catch((err) => console.error("Error loading game:", err));

    fetch(`http://localhost:8080/games/${gameId}/objectives`)
      .then((res) => res.json())
      .then((data) => setObjectives(Array.isArray(data) ? data : data?.value || []))
      .catch((err) => console.error("Error loading objectives:", err));

    fetch(`http://localhost:8080/games/${gameId}/objectives/scores`)
      .then((res) => res.json())
      .then((data) => {
        const list = Array.isArray(data) ? data : data?.value || [];
        const map = {};
        list.forEach((entry) => {
          map[entry.objective_id] = entry.scored_by || [];
        });
        setObjectiveScores(map);
      })
      .catch((err) => console.error("Error loading objective scores:", err));
  }, [gameId]);

  const scoreObjective = async (playerId, objectiveId, playerName) => {
    const payload = {
      game_id: parseInt(gameId),
      player_id: playerId,
      objective_id: objectiveId,
    };

    try {
      await fetch("http://localhost:8080/score", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });

      setLocalScored((prev) => {
        const current = prev[objectiveId] || [];
        return {
          ...prev,
          [objectiveId]: [...new Set([...current, playerName])],
        };
      });

      const updatedGame = await fetch(`http://localhost:8080/games/${gameId}`).then((r) => r.json());
      setGame(updatedGame);

      const updatedObjectives = await fetch(`http://localhost:8080/games/${gameId}/objectives`).then((r) => r.json());
      setObjectives(Array.isArray(updatedObjectives) ? updatedObjectives : updatedObjectives?.value || []);

      const updatedScores = await fetch(`http://localhost:8080/games/${gameId}/objectives/scores`).then((r) => r.json());
      const map = {};
      (Array.isArray(updatedScores) ? updatedScores : updatedScores?.value || []).forEach((entry) => {
        map[entry.objective_id] = entry.scored_by || [];
      });
      setObjectiveScores(map);
    } catch (err) {
      console.error("Scoring failed:", err);
      alert("Scoring failed. See console.");
    }
  };

  const advanceRound = async () => {
    try {
      const res = await fetch(`http://localhost:8080/games/${gameId}/advance-round`, {
        method: "POST",
      });

      if (!res.ok) {
        const errorText = await res.text();
        throw new Error(errorText);
      }

      const gameRes = await fetch(`http://localhost:8080/games/${gameId}`);
      if (!gameRes.ok) throw new Error("Failed to fetch updated game after advancing");
      const updatedGame = await gameRes.json();
      setGame(updatedGame);

      const updatedObjectives = await fetch(`http://localhost:8080/games/${gameId}/objectives`).then((r) => r.json());
      setObjectives(Array.isArray(updatedObjectives) ? updatedObjectives : updatedObjectives?.value || []);

      const updatedScores = await fetch(`http://localhost:8080/games/${gameId}/objectives/scores`).then((r) => r.json());
      const map = {};
      (Array.isArray(updatedScores) ? updatedScores : updatedScores?.value || []).forEach((entry) => {
        map[entry.objective_id] = entry.scored_by || [];
      });
      setObjectiveScores(map);
    } catch (err) {
      console.error("Failed to advance round:", err);
      alert("Could not advance round. See console for details.");
    }
  };

  const getMergedPlayerData = (sort = true) => {
    const scoreMap = new Map(game?.scores?.map((s) => [s.player_id, s]) || []);

    const merged = (game?.players || []).map((p) => {
      const name = p.Player?.Name || "Unknown";
      const faction = p.Faction || "Unknown Faction";
      const factionKey =
        faction.replace(/^The\s+/i, "")
          .replace(/\s+/g, "")
          .replace(/[^a-zA-Z0-9]/g, "") +
        (faction.toLowerCase().includes("keleres") ? "FactionSymbol" : "");
      return {
        player_id: p.PlayerID,
        name,
        faction,
        factionKey,
        color: p.color || "#000",
        points: scoreMap.get(p.PlayerID)?.points || 0,
      };
    });

    return sort ? merged.sort((a, b) => b.points - a.points) : merged;
  };

  if (!game || !game.players) return <div className="p-6">Loading game data...</div>;

  const playersUnsorted = getMergedPlayerData(false);
  const playersSorted = getMergedPlayerData(true);

  return (
    <div className="p-6 max-w-7xl mx-auto">
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
          <button className="btn btn-outline-primary btn-sm" onClick={advanceRound}>
            Advance Round
          </button>
        </div>
      </div>

      <div className="d-flex flex-row align-items-start gap-4 flex-wrap">
        <div style={{ flex: "1 1 0" }}>
          <h4>Objectives</h4>
          <div className="d-flex flex-wrap justify-content-start" style={{ gap: "20px" }}>
            {objectives.map((obj) => {
              const objId = obj.Objective?.ID;
              const isStageTwo = obj.Objective?.stage === "II";
              const stageColor = isStageTwo ? "#00bfff" : "#ffd700";
              const glowColor = isStageTwo ? "rgba(0, 191, 255, 0.4)" : "rgba(255, 215, 0, 0.4)";
              const backgroundImage = isStageTwo
                ? "/objective-backgrounds/stage2.jpg"
                : "/objective-backgrounds/stage1.jpg";

              const scoredBy = [
                ...(objectiveScores[objId] || []),
                ...(localScored[objId] || []),
              ];

              return (
                <div
                  key={obj.ID}
                  className="card shadow"
                  style={{
                    width: "220px",
                    height: "330px",
                    backgroundImage: `url(${backgroundImage})`,
                    backgroundSize: "contain",
                    backgroundRepeat: "no-repeat",
                    backgroundPosition: "center",
                    border: `2px solid ${stageColor}`,
                    borderRadius: "12px",
                    color: "#fff",
                    fontFamily: "'Orbitron', sans-serif",
                    boxShadow: `0 0 10px ${glowColor}`,
                    position: "relative",
                  }}
                >
                  <div
                    style={{
                      position: "absolute",
                      inset: 0,
                      backgroundColor: "rgba(0, 0, 0, 0.6)",
                      padding: "12px",
                      display: "flex",
                      flexDirection: "column",
                      justifyContent: "space-between",
                    }}
                  >
                    <div>
                      <h5>{obj.Objective?.name || "Unnamed Objective"}</h5>
                      <p className="small fst-italic" style={{ color: "#ccc" }}>
                        {obj.Objective?.description || "No description provided."}
                      </p>
                    </div>
                    <div className="d-flex flex-wrap gap-2 mt-3">
                      {playersUnsorted.map((p) => {
                        const hasScored = scoredBy.includes(p.name);

                        return (
                          <div
                            key={p.player_id}
                            style={{
                              width: "32px",
                              height: "32px",
                              display: "flex",
                              alignItems: "center",
                              justifyContent: "center",
                            }}
                          >
                            {scoringMode ? (
                              <button
                                className="btn btn-sm p-1"
                                style={{
                                  backgroundColor: p.color,
                                  borderRadius: "6px",
                                  width: "100%",
                                  height: "100%",
                                  display: "flex",
                                  alignItems: "center",
                                  justifyContent: "center",
                                }}
                                onClick={() => scoreObjective(p.player_id, objId, p.name)}
                              >
                                <img
                                  src={`/faction-icons/${p.factionKey}.webp`}
                                  alt={p.faction}
                                  style={{ width: "20px", height: "20px", objectFit: "contain" }}
                                  onError={(e) => (e.target.style.display = "none")}
                                />
                              </button>
                            ) : hasScored ? (
                              <img
                                src={`/faction-icons/${p.factionKey}.webp`}
                                alt={p.faction}
                                style={{
                                  width: "24px",
                                  height: "24px",
                                  borderRadius: "50%",
                                  objectFit: "contain",
                                  backgroundColor: "transparent",
                                }}
                                onError={(e) => (e.target.style.display = "none")}
                              />
                            ) : (
                              <div style={{ width: "24px", height: "24px" }} />
                            )}
                          </div>
                        );
                      })}
                    </div>
                    <div className="text-end">
                      <span
                        className="badge"
                        style={{
                          backgroundColor: stageColor,
                          color: "#000",
                          fontWeight: "bold",
                          fontSize: "0.75rem",
                        }}
                      >
                        {obj.Objective?.type?.toUpperCase() || "PUBLIC"}
                      </span>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        </div>

        <div style={{ flex: "0 1 300px" }}>
          {playersSorted.map((entry) => (
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
                    {custodiansScored && game?.AllScores?.some(s => s.Type === "mecatol" && s.PlayerID === entry.player_id) && (
                      <img
                        src="public/MR-point/MR-scored.png"
                        alt="Custodians Point"
                        title="Custodians Point"
                        style={{ width: "20px", height: "20px" }}
                      />
                    )}
                  <div className="text-muted small fst-italic">{entry.faction}</div>
                </div>
                <div className="mt-1 small">Points: {entry.points}</div>
                {expandedPlayers[entry.player_id] && (
                  <div className="mt-3 small">
                    <button
                      className="btn btn-warning btn-sm"
                      disabled={custodiansScored}
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

                          // Refresh game and scores
                          const updatedGame = await fetch(`http://localhost:8080/games/${gameId}`).then((r) => r.json());
                          setGame(updatedGame);

                          const updatedScores = await fetch(`http://localhost:8080/games/${gameId}/objectives/scores`).then((r) => r.json());
                          const map = {};
                          (Array.isArray(updatedScores) ? updatedScores : updatedScores?.value || []).forEach((entry) => {
                            map[entry.objective_id ?? entry.name] = entry.scored_by || [];
                          });
                          setObjectiveScores(map);
                        } catch (err) {
                          console.error("Failed to score Custodians:", err);
                          alert("Failed to score Custodians. See console.");
                        }
                      }}
                    >
                      {custodiansScored ? "Custodians Already Scored" : "Score Custodians"}
                    </button>
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
