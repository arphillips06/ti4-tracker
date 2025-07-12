import React, { useEffect, useState } from "react";
import API_BASE_URL from "../../config";
export default function ObjectivesGrid({
  game,
  refreshGameState,
  gameId,
  objectives,
  playersUnsorted,
  objectiveScores,
  localScored,
  scoringMode,
  scoreObjective,
  useObjectiveDecks,
  setAssigningObjective,
  assigningObjective,
  assignObjective,
}) {
  const factionImageMap = {
    "Arborec": "Arborec.webp",
    "Argent Flight": "ArgentFlight.webp",
    "Barony of Letnev": "BaronyofLetnev.webp",
    "Clan of Saar": "ClanofSaar.webp",
    "Council Keleres": "CouncilKeleresFactionSymbol.webp",
    "Embers of Muaat": "EmbersofMuaat.webp",
    "Emirates of Hacan": "EmiratesofHacan.webp",
    "Empyrean": "Empyrean.webp",
    "Federation of Sol": "FederationofSol.webp",
    "Ghosts of Creuss": "GhostsofCreuss.webp",
    "L1Z1X Mindnet": "L1Z1XMindnet.webp",
    "Mahat Gene-Sorcerers": "MahatGeneSorcerers.webp",
    "Mentak Coalition": "MentakCoalition.webp",
    "Naalu Collective": "NaaluCollective.webp",
    "Naaz-Rokha Alliance": "NaazRokhaAlliance.webp",
    "Nekro Virus": "NekroVirus.webp",
    "Nomad": "Nomad.webp",
    "Sardakk N'orr": "SardakkNorr.webp",
    "Titans of Ul": "TitansofUl.webp",
    "Universities of Jol-Nar": "UniversitiesofJolNar.webp",
    "Vuil'raith Cabal": "VuilraithCabal.webp",
    "Winnu": "Winnu.webp",
    "Xxcha Kingdom": "XxchaKingdom.webp",
    "Yin Brotherhood": "YinBrotherhood.webp",
    "Yssaril Tribes": "YssarilTribes.webp",
  };

  const safeObjectives = objectives || [];

  const safePlayers = (playersUnsorted || []).map((p) => {
    const faction = p.Faction || p.faction || p.Player?.Faction || p.Player?.faction;
    const name = p.name || p.Player?.name;
    const color = p.color || p.Player?.color;
    const id = p.PlayerID || p.player_id || p.Player?.ID || p.id;


    return {
      ...p,
      id,
      name,
      color,
      faction,
      factionKey: factionImageMap[faction] ? factionImageMap[faction] : "fallback.webp"

    };
  });

  const safeScores = objectiveScores || {};
  const normalizedScores = {};
  Object.entries(safeScores).forEach(([objId, scoreEntries]) => {
    normalizedScores[objId] = scoreEntries.map(
      (s) => s.player_id || s.Player?.ID
    );
  });

  const safeLocalScored = localScored || {};
  const usingDecks = String(useObjectiveDecks).toLowerCase() === "true";

  const [publicObjectives, setPublicObjectives] = useState([]);
  const isManualMode = !usingDecks;
  const availableStageI = publicObjectives.filter((o) => o.stage === "I");
  const availableStageII = publicObjectives.filter((o) => o.stage === "II");
  const currentRoundId = game?.current_round || 0;

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
                        src={`/faction-icons/${p.factionKey}`}
                        alt={p.faction}
                        onError={(e) => {
                          console.warn("🚫 Could not load faction icon for", p.factionKey);
                          e.target.style.display = "none";
                        }}
                        style={{
                          width: "20px",
                          height: "20px",
                          objectFit: "contain",
                        }}
                      />
                    </button>
                  ) : hasScored ? (
                    <img
                      src={`/faction-icons/${p.factionKey}`}
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
                  const selectedId = parseInt(e.target.value);
                  if (selectedId) {
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
