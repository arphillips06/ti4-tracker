import React from "react";
import API_BASE_URL from "../../config";
import "./playersidebar.css";

export default function PlayerSidebar({
  playersSorted,
  expandedPlayers,
  setExpandedPlayers,
  game,
  allScores,
  setAllScores,
  scoreObjective,
  unscoreObjective,
  setSecretCounts,
  secretObjectives,
  gameId,
  setGame,
  setObjectiveScores,
  refreshGameState,
  custodiansScored,          // optional prop; we’ll fallback if missing
  obsidianHolderId,
  triggerGraphUpdate,
}) {
  // Normalize the score type safely (handles either `Type` or `type`)
  const scoreType = (s) => (s?.Type || s?.type || "").toLowerCase();


  // Derive whether Custodians is scored if the parent didn't pass the prop
  const custodiansScoredLocal = allScores?.some((s) => scoreType(s) === "mecatol");
  const showImperial = typeof custodiansScored === "boolean" ? custodiansScored : custodiansScoredLocal;

  // CDL: set of objective IDs that became public due to Classified Document Leaks
  const cdlRevealedObjectiveIds = new Set(
    (allScores || [])
      .filter((s) => s.AgendaTitle === "Classified Document Leaks")
      .map((s) => s.ObjectiveID)
  );

  // Is a scored secret still secret (not revealed by CDL)?
  const isStillSecret = (score) => {
    if (scoreType(score) !== "secret") return false;
    return !cdlRevealedObjectiveIds.has(score.ObjectiveID);
  };

  // Reusable SFTT action
  async function postSupportAction(playerId, action) {
    try {
      const res = await fetch(`${API_BASE_URL}/games/${gameId}/support/${playerId}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          round_id: game?.current_round_id,
          action, // "score" | "unscore"
        }),
      });
      if (!res.ok) {
        const err = await res.json().catch(() => ({}));
        throw new Error(err.error || `HTTP ${res.status}`);
      }
      await refreshGameState?.();
      triggerGraphUpdate?.();
    } catch (e) {
      console.error("Support action failed:", e);
      alert("Failed to update Support for the Throne.");
    }
  }

  return (
    <div className="flex-shrink-1" style={{ minWidth: "250px", flexBasis: "300px" }}>
      {(playersSorted || []).map((entry) => {
        return (
          <div
            key={entry.player_id}
            className="card mb-3 glass-box"
            style={{ borderColor: entry.color }}
          >
            <div className="card-body">
              <div className="d-flex justify-content-between align-items-center">
                <div className="fw-semibold small d-flex align-items-center gap-2">
                  {entry.name}
                  {game?.speaker_id === entry.id && (  // ✅ compare to player_id
                    <img
                      src="/speaker/speaker.webp"
                      alt="Speaker"
                      title="Speaker"
                      style={{ width: "30px", height: "auto", objectFit: "contain" }}
                    />
                  )}
                </div>
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
                  onError={(e) => (e.currentTarget.style.display = "none")}
                />
                {allScores?.some(
                  (s) => scoreType(s) === "mecatol" && s.PlayerID === entry.player_id
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
                    {/* ===== Support for the Throne ===== */}
                    <div className="small fw-semibold mb-1">Support for the Throne</div>
                    <div className="d-flex align-items-center gap-2 mb-3">
                      <button
                        className="btn btn-sm btn-outline-danger"
                        disabled={
                          !allScores?.some(
                            (s) => (s?.Type || s?.type) === "Support" && s.PlayerID === entry.player_id
                          )
                        }
                        onClick={() => postSupportAction(entry.player_id, "unscore")}
                      >
                        −
                      </button>

                      <span className="small">
                        {(() => {
                          const total =
                            allScores
                              ?.filter((s) => (s?.Type || s?.type) === "Support" && s.PlayerID === entry.player_id)
                              .reduce((sum, s) => sum + (s.Points || 0), 0) || 0;
                          return `${total} Support point${total === 1 ? "" : "s"}`;
                        })()}
                      </span>

                      <button
                        className="btn btn-sm btn-outline-success"
                        disabled={(() => {
                          const n = playersSorted?.length || 0;
                          const playerSupportPoints =
                            allScores
                              ?.filter((s) => (s?.Type || s?.type) === "Support" && s.PlayerID === entry.player_id)
                              .reduce((acc, s) => acc + (s.Points || 0), 0) || 0;
                          return playerSupportPoints >= Math.max(0, n - 1);
                        })()}
                        onClick={() => postSupportAction(entry.player_id, "score")}
                      >
                        +
                      </button>
                    </div>

                    {/* ===== Secrets ===== */}
                    <div className="mt-3">
                      <div className="small fw-semibold mb-1">Secrets</div>
                      <div className="d-flex align-items-center gap-2 mb-2">
                        {/* Unscore dropdown (only still-secret ones) */}
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
                            {allScores
                              ?.filter((s) => s.PlayerID === entry.player_id && isStillSecret(s))
                              .map((s) => {
                                const obj = secretObjectives.find((o) => o.id === s.ObjectiveID);
                                return obj ? (
                                  <li key={obj.id}>
                                    <button
                                      className="dropdown-item"
                                      onClick={async () => {
                                        const ok = await unscoreObjective?.(entry.player_id, obj.id);
                                        if (ok) {
                                          setSecretCounts?.((prev) => ({
                                            ...prev,
                                            [entry.player_id]: Math.max(0, (prev?.[entry.player_id] || 0) - 1),
                                          }));
                                          triggerGraphUpdate?.();
                                        }
                                      }}
                                    >
                                      {obj.name}
                                    </button>
                                  </li>
                                ) : null;
                              })}
                          </ul>
                        </div>

                        {/* Secret slots (icons) */}
                        <div className="d-flex gap-1">
                          {(() => {
                            const extraSecret =
                              parseInt(entry.player_id) === parseInt(obsidianHolderId) ? 1 : 0;
                            const baseSecrets = 3;
                            const maxSecrets = baseSecrets + extraSecret;

                            const scoredSecrets = (allScores || []).filter((s) => {
                              const isSecret = scoreType(s) === "secret";
                              const isThisPlayer = s.PlayerID === entry.player_id;
                              const isCDL = cdlRevealedObjectiveIds.has(s.ObjectiveID);
                              return isThisPlayer && isSecret && !isCDL;
                            });

                            return [...Array(maxSecrets)].map((_, i) => {
                              const secret = scoredSecrets[i];
                              const scored = !!secret;
                              return (
                                <img
                                  key={i}
                                  src={`/objective-backgrounds/secret-${scored ? "active" : "inactive"}.jpg`}
                                  alt={scored ? "Scored secret" : "Unscored secret"}
                                  style={{ width: "16px", height: "25px", opacity: scored ? 1 : 0.4 }}
                                />
                              );
                            });
                          })()}
                        </div>
                      </div>

                      {/* Score new secret objective */}
                      <select
                        className="secret-objective-select"
                        value=""
                        onChange={async (e) => {
                          const selectedId = parseInt(e.target.value);
                          if (selectedId) {
                            const success = await scoreObjective(entry.player_id, selectedId, entry.name);
                            if (success) {
                              setSecretCounts?.((prev) => ({
                                ...prev,
                                [entry.player_id]: Math.min(3, (prev[entry.player_id] || 0) + 1),
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
                                  disabled={allScores?.some(
                                    (s) =>
                                      s.PlayerID === entry.player_id &&
                                      s.ObjectiveID === obj.id &&
                                      (scoreType(s) === "secret" ||
                                        s.AgendaTitle === "Classified Document Leaks")
                                  )}
                                >
                                  {obj.name}
                                </option>
                              ))}
                          </optgroup>
                        ))}
                      </select>

                      {/* Small agenda badges (kept from your working file) */}
                      {allScores?.some((s) => s.PlayerID === entry.player_id && s.AgendaTitle === "Mutiny") && (
                        <div className="mt-1 small text-success">Bonus: Mutiny</div>
                      )}
                      {allScores?.some(
                        (s) => s.PlayerID === entry.player_id && s.AgendaTitle === "Seed of an Empire"
                      ) && <div className="mt-1 small text-success">Bonus: Seed of an Empire</div>}

                      {/* ===== Custodians (Mecatol) ===== */}
                      <div className="mt-3 small">
                        <button
                          className="btn btn-warning btn-sm"
                          disabled={allScores?.some((s) => scoreType(s) === "mecatol")}
                          onClick={async () => {
                            try {
                              const res = await fetch(`${API_BASE_URL}/score/mecatol`, {
                                method: "POST",
                                headers: { "Content-Type": "application/json" },
                                body: JSON.stringify({
                                  game_id: parseInt(gameId),
                                  player_id: entry.player_id,
                                }),
                              });
                              if (!res.ok) {
                                const err = await res.json().catch(() => ({}));
                                throw new Error(err.error || "Failed to score Custodians");
                              }

                              // Refresh
                              const [gameRes, objScoresRes] = await Promise.all([
                                fetch(`${API_BASE_URL}/games/${gameId}`).then((r) => r.json()),
                                fetch(`${API_BASE_URL}/games/${gameId}/objectives/scores`).then((r) => r.json()),
                              ]);

                              const updatedAllScores = gameRes.all_scores || [];
                              gameRes.game_players = gameRes.players || [];

                              setGame(gameRes);
                              setAllScores(updatedAllScores);

                              const map = {};
                              (Array.isArray(objScoresRes) ? objScoresRes : objScoresRes?.value || []).forEach(
                                (row) => {
                                  map[row.objective_id ?? row.name] = row.scored_by || [];
                                }
                              );
                              setObjectiveScores(map);

                              await refreshGameState?.();
                              triggerGraphUpdate?.();
                            } catch (err) {
                              console.error("Failed to score Custodians:", err);
                              alert("Failed to score Custodians. See console.");
                            }
                          }}
                        >
                          {allScores?.some((s) => scoreType(s) === "mecatol")
                            ? "Custodians Already Scored"
                            : "Score Custodians"}
                        </button>
                      </div>
                    </div>

                    {/* ===== Imperial (after Custodians) ===== */}
                    {showImperial && (
                      <div className="mt-3 small">
                        <div className="fw-semibold mb-1">Imperial Objective</div>
                        <div className="d-flex align-items-center gap-2 flex-wrap">
                          {allScores
                            ?.filter((s) => scoreType(s) === "imperial" && s.PlayerID === entry.player_id)
                            .map((_, i) => (
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
                    )}
                  </div>
                )}
              </div>
            </div>
          </div>
        );
      })}
    </div>
  );
}
