import React from "react";
import API_BASE_URL from "../../config";
import './playersidebar.css';
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
  obsidianHolderId,
  triggerGraphUpdate,
}) {

  const supportScorers = new Set(
    game?.AllScores?.filter((s) => s.Type === "Support").map((s) => s.PlayerID)
  );
  const maxSupportScorers = (playersSorted?.length || 0) - 1;

  const custodiansScorerId = game?.AllScores?.find((s) => s.Type === "mecatol")?.PlayerID || null;

  // Step 1: Gather all CDL-revealed secret objective IDs across all players
  const cdlRevealedObjectiveIds = new Set(
    game?.AllScores
      ?.filter((s) => s.AgendaTitle === "Classified Document Leaks")
      .map((s) => s.ObjectiveID)
  );

  // Utility: Determine if a scored objective is still secret (not publicly revealed)
  const isStillSecret = (score) => {
    if ((score.Type || "").toLowerCase() !== "secret") return false;
    // If this secret objective is not revealed by CDL, it remains secret
    return !cdlRevealedObjectiveIds.has(score.ObjectiveID);
  };

  return (
    <div style={{ flex: "0 1 300px" }}>
      {(playersSorted || []).map((entry) => (
        <div
          key={entry.player_id}
          className="card mb-3 glass-box"
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
              <div className="small fst-italic faction-name">{entry.faction}</div>
            </div>

            <div className="mt-1 small">Points: {entry.points}</div>

            <div className="d-flex flex-column gap-2 mt-2">
              {expandedPlayers[entry.player_id] && (
                <div className="mt-3">
                  <div className="small fw-semibold mb-1">Support for the Throne</div>
                  <div className="d-flex align-items-center gap-2 mb-3">
                    <button
                      className="btn btn-sm btn-outline-danger"
                      disabled={
                        !game?.AllScores?.some(
                          (s) => s.Type === "Support" && s.PlayerID === entry.player_id
                        )
                      }
                      onClick={async () => {
                        const res = await fetch(`${API_BASE_URL}/games/${gameId}/support/${entry.player_id}`, {
                          method: "POST",
                          headers: { "Content-Type": "application/json" },
                          body: JSON.stringify({
                            round_id: game?.current_round_id,
                            action: "unscore",
                          }),
                        });

                        if (res.ok) {
                          await refreshGameState();
                          triggerGraphUpdate?.();
                        } else {
                          const err = await res.json();
                          alert(err.error || "Failed to remove Support");
                        }
                      }}
                    >
                      −
                    </button>

                    <span>
                      {(() => {
                        const total = game?.AllScores
                          ?.filter((s) => s.Type === "Support" && s.PlayerID === entry.player_id)
                          .reduce((sum, s) => sum + s.Points, 0) || 0;

                        return `${total} Support point${total === 1 ? "" : "s"}`;
                      })()}
                    </span>

                    <button
                      className="btn btn-sm btn-outline-success"
                      disabled={
                        (() => {
                          const totalSupportPoints = game?.AllScores
                            ?.filter((s) => s.Type === "Support")
                            .reduce((acc, s) => acc + s.Points, 0) || 0;
                          return totalSupportPoints >= maxSupportScorers;
                        })()
                      }
                      onClick={async () => {
                        const res = await fetch(`${API_BASE_URL}/games/${gameId}/support/${entry.player_id}`, {
                          method: "POST",
                          headers: { "Content-Type": "application/json" },
                          body: JSON.stringify({
                            round_id: game?.current_round_id,
                            action: "score",
                          }),
                        });

                        if (res.ok) {
                          await refreshGameState();
                          triggerGraphUpdate?.();
                        } else {
                          const err = await res.json();
                          alert(err.error || "Failed to score Support");
                        }
                      }}
                    >
                      +
                    </button>
                  </div>
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
                              isStillSecret(s)
                          ).map((s) => {
                            const obj = secretObjectives.find((o) => o.id === s.ObjectiveID);
                            return obj ? (
                              <li key={obj.id}>
                                <button className="dropdown-item">{obj.name}</button>
                              </li>
                            ) : null;
                          })}
                        </ul>
                      </div>

                      <div className="d-flex gap-1">
                        {(() => {

                          const extraSecret = parseInt(entry.player_id) === parseInt(obsidianHolderId) ? 1 : 0;
                          const baseSecrets = 3;
                          const maxSecrets = baseSecrets + extraSecret;

                          // Filter secrets excluding CDL revealed ones (public)
                          const scoredSecrets = game?.AllScores?.filter(
                            (s) =>
                              s.PlayerID === entry.player_id &&
                              (s.Type || "").toLowerCase() === "secret" &&
                              !cdlRevealedObjectiveIds.has(s.ObjectiveID)
                          ) || [];

                          return [...Array(maxSecrets)].map((_, i) => {
                            const secret = scoredSecrets[i]; // might be undefined
                            const scored = !!secret;
                            // If secret is scored but revealed public by CDL it won't be here (excluded)

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
                          });
                        })()}
                      </div>
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
                            triggerGraphUpdate?.();
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
                                disabled={game?.AllScores?.some(
                                  (s) =>
                                    s.PlayerID === entry.player_id &&
                                    s.ObjectiveID === obj.id &&
                                    ((s.Type?.toLowerCase() === "secret") || s.AgendaTitle === "Classified Document Leaks")
                                )}

                              >
                                {obj.name}
                              </option>
                            ))}
                        </optgroup>
                      ))}
                    </select>

                    {game.AllScores?.some(
                      (s) =>
                        s.PlayerID === entry.player_id &&
                        s.AgendaTitle === "Mutiny"
                    ) && (
                        <div className="mt-1 small text-success">Bonus: Mutiny</div>
                      )}
                    {game.AllScores?.some(
                      (s) =>
                        s.PlayerID === entry.player_id &&
                        s.AgendaTitle === "Seed of an Empire"
                    ) && (
                        <div className="mt-1 small text-success">Bonus: Seed of an Empire</div>
                      )}

                    <div className="mt-3 small">
                      <button
                        className="btn btn-warning btn-sm"
                        disabled={game?.AllScores?.some((s) => s.Type === "mecatol")}
                        onClick={async () => {
                          try {
                            await fetch(`${API_BASE_URL}/score/mecatol`, {
                              method: "POST",
                              headers: { "Content-Type": "application/json" },
                              body: JSON.stringify({
                                game_id: parseInt(gameId),
                                player_id: entry.player_id,
                              }),
                            });
                            await refreshGameState();
                            triggerGraphUpdate?.();

                            const updatedGame = await fetch(`${API_BASE_URL}/games/${gameId}`).then((r) => r.json());
                            updatedGame.AllScores = updatedGame.all_scores || [];
                            updatedGame.game_players = updatedGame.players || [];
                            setGame(updatedGame);

                            const updatedScores = await fetch(
                              `${API_BASE_URL}/games/${gameId}/objectives/scores`
                            ).then((r) => r.json());

                            const map = {};
                            (Array.isArray(updatedScores)
                              ? updatedScores
                              : updatedScores?.value || []
                            ).forEach((entry) => {
                              map[entry.objective_id ?? entry.name] = entry.scored_by || [];
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
                  </div>
                  {custodiansScored && (
                    <div className="mt-3 small">
                      <div className="fw-semibold mb-1">Imperial Objective</div>

                      <div className="d-flex align-items-center gap-2 flex-wrap">
                        {game?.AllScores?.filter(
                          (s) => s.Type === "imperial" && s.PlayerID === entry.player_id
                        ).map((_, i) => (
                          <img
                            key={i}
                            src="/imperial/imperial8.png"
                            alt="Imperial Point"
                            title="Imperial Point"
                            style={{ width: "32px", height: "48px" }}
                          />
                        ))}

                        <button
                          className="btn btn-sm btn-outline-primary"
                          onClick={async () => {
                            try {
                              const res = await fetch(`${API_BASE_URL}/score/imperial`, {
                                method: "POST",
                                headers: { "Content-Type": "application/json" },
                                body: JSON.stringify({
                                  game_id: parseInt(gameId),
                                  player_id: entry.player_id,
                                  round_id: game?.current_round_id,
                                }),
                              });

                              if (res.ok) {
                                await refreshGameState();
                                triggerGraphUpdate?.();
                              } else {
                                const err = await res.json();
                                alert(err.error || "Failed to score Imperial");
                              }
                            } catch (err) {
                              console.error("Failed to score Imperial:", err);
                              alert("Failed to score Imperial. See console.");
                            }
                          }}
                        >
                          +
                        </button>
                      </div>
                    </div>
                  )}</div>
              )}</div>
          </div>
        </div>
      ))}
    </div>
  )
}
