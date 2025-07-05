import React from "react";

export default function ObjectivesGrid({
  objectives,
  playersUnsorted,
  objectiveScores,
  localScored,
  scoringMode,
  scoreObjective,
}) {
  return (
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
                            onClick={() =>
                              scoreObjective(p.player_id, objId, p.name)
                            }
                          >
                            <img
                              src={`/faction-icons/${p.factionKey}.webp`}
                              alt={p.faction}
                              style={{
                                width: "20px",
                                height: "20px",
                                objectFit: "contain",
                              }}
                              onError={(e) =>
                                (e.target.style.display = "none")
                              }
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
                            onError={(e) =>
                              (e.target.style.display = "none")
                            }
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
  );
}
