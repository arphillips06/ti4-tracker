import React, { useEffect, useState } from "react";
import API_BASE_URL from "../../config";
import factionImageMap from "../../data/factionIcons";

// ✅ Move these ABOVE where they're used
const getFactionKey = (faction) => {
  return factionImageMap[faction] || "fallback";
};

const normalizePlayer = (p) => {
  const rawFaction =
    p.faction || p.Faction || p.Player?.Faction || p.Player?.faction || "Unknown";

  const factionKey =
    p.factionKey || p.FactionKey || getFactionKey(rawFaction);

  return {
    id: p.Player?.ID || p.player_id || p.PlayerID || p.ID,
    name: p.name || p.Name || p.Player?.name || p.Player?.Name || "Unknown",
    faction: rawFaction,
    factionKey,
    color: p.color || "#000",
  };
};

export default function ObjectivesGrid({
  game,
  gameId,
  objectives,
  playersUnsorted,
  localScored,
  scoringMode,
  scoreObjective,
  useObjectiveDecks,
  setAssigningObjective,
  assigningObjective,
  assignObjective,
}) {
  const safeObjectives = objectives || [];
  const safePlayers = (playersUnsorted || []).map(normalizePlayer); // ✅ now safe
  const rawScores = game?.ScoresByObjective || {};
  const safeLocalScored = localScored || {};
  const usingDecks = String(useObjectiveDecks).toLowerCase() === "true";
  const [publicObjectives, setPublicObjectives] = useState([]);
  const isManualMode = !usingDecks;
  const availableStageI = publicObjectives.filter((o) => o.stage === "I");
  const availableStageII = publicObjectives.filter((o) => o.stage === "II");
  const currentRoundId = game?.current_round || 1;

  const normalizedScores = {};
  Object.entries(rawScores).forEach(([objId, entries]) => {
    normalizedScores[objId] = entries.map(
      (s) => s.PlayerID || s.Player?.ID
    );
  });

  useEffect(() => {
    if (!usingDecks) {
      fetch(`${API_BASE_URL}/objectives/public/all`)
        .then((res) => res.json())
        .then((data) => setPublicObjectives(data))
        .catch((err) => console.error("Failed to load public objectives:", err));
    }
  }, [usingDecks]);

  const displayObjectives = safeObjectives;

  const renderObjectiveCard = (obj) => {
    const objId = obj.Objective?.ID;
    const isStageTwo = obj.Objective?.stage === "II";
    const isCDL = obj.IsCDL;

    const stageColor = isCDL ? "#d63384" : (isStageTwo ? "#00bfff" : "#ffd700");
    const glowColor = isCDL ? "rgba(214, 51, 132, 0.4)" : (isStageTwo ? "rgba(0, 191, 255, 0.4)" : "rgba(255, 215, 0, 0.4)");

    const backgroundImage = isCDL
      ? "/objective-backgrounds/secret-active.jpg"
      : isStageTwo
        ? "/objective-backgrounds/stage2.jpg"
        : "/objective-backgrounds/stage1.jpg";

    const scoredBy = [
      ...(normalizedScores[objId] || []),
      ...(safeLocalScored[objId] || []),
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
            {safePlayers.map((p) => {
              const hasScored = scoredBy.includes(p.id);
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
                      onClick={() =>
                        scoreObjective(p.id, objId, p.name)
                      }
                    >
                      <img
                        src={`/faction-icons/${p.factionKey}.webp`}
                        alt={p.faction}
                        onError={(e) => {
                          console.warn("🚫 Could not load faction icon for", p.factionKey);
                          if (!e.target.src.includes("fallback.webp")) {
                            e.target.src = "/faction-icons/fallback.webp";
                          }
                        }}
                        style={{
                          width: "24px",
                          height: "24px",
                          borderRadius: "50%",
                          objectFit: "contain",
                          backgroundColor: "transparent",
                        }}
                      />

                    </button>
                  ) : hasScored ? (
                    <img
                      src={`/faction-icons/${p.factionKey}.webp`}
                      alt={p.faction}
                      onError={(e) => {
                        console.warn("🚫 Could not load faction icon for", p.factionKey);
                        e.target.src = "/faction-icons/fallback.webp";
                      }}
                      style={{
                        width: "24px",
                        height: "24px",
                        borderRadius: "50%",
                        objectFit: "contain",
                        backgroundColor: "transparent",
                      }}
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
              {obj.IsCDL ? "CDL" : (obj.Objective?.type?.toUpperCase() || "PUBLIC")}
            </span>
          </div>
        </div>
      </div>
    );
  };

  return (
    <div style={{ flex: "1 1 0" }}>
      <h4>Objectives</h4>

      {isManualMode && (
        <>
          <div className="mb-3 text-warning small fst-italic">
            Manual mode active (Use Objective Decks disabled).
          </div>
          <div className="d-flex gap-3 mb-4">
            <button
              className="btn btn-warning"
              onClick={() =>
                setAssigningObjective({
                  roundId: currentRoundId,
                  stage: "I",
                })
              }
            >
              + Assign Stage I Objective
            </button>
            <button
              className="btn btn-info"
              onClick={() =>
                setAssigningObjective({
                  roundId: currentRoundId,
                  stage: "II",
                })
              }
            >
              + Assign Stage II Objective
            </button>
          </div>

          {assigningObjective && (
            <div className="mb-4">
              <label className="form-label">
                Select a Stage {assigningObjective.stage} Objective to assign
              </label>

              <select
                className="form-select"
                onChange={async (e) => {
                  const selectedId = parseInt(e.target.value, 10);
                  if (selectedId && assigningObjective?.roundId) {
                    await assignObjective(gameId, assigningObjective.roundId, selectedId);
                    setAssigningObjective(null);
                  }
                }}
              >
                <option value="">-- Select Objective --</option>
                {(assigningObjective.stage === "I" ? availableStageI : availableStageII).map(
                  (obj) => (
                    <option key={obj.ID} value={obj.ID}>
                      {obj.name} ({obj.description})
                    </option>
                  )
                )}
              </select>
            </div>
          )}
        </>
      )}

      <div className="d-flex flex-wrap justify-content-start" style={{ gap: "20px" }}>
        {displayObjectives.map(renderObjectiveCard)}
      </div>
    </div>
  );
}
