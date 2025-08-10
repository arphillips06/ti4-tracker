import React from "react";
import "./shared/FactionObjectiveModal.css";

export default function FactionObjectiveModal({ faction, objectiveStats, onClose }) {
  const stageI = Object.entries(objectiveStats).filter(([_, s]) => s.type === "stage1");
  const stageII = Object.entries(objectiveStats).filter(([_, s]) => s.type === "stage2");
  const secret = Object.entries(objectiveStats).filter(([_, s]) => s.type === "secret");

  const renderTable = (title, data) => (
    <>
      <h5>{title}</h5>
      {data.length === 0 ? <p>No data</p> : (
        <table className="table table-sm table-striped">
          <thead>
            <tr>
              <th>Objective</th>
              <th>Appeared</th>
              <th>Scored</th>
              <th>% Scored</th>
            </tr>
          </thead>
          <tbody>
            {data.map(([name, s]) => (
              <tr key={name}>
                <td>{name}</td>
                <td>{s.appearedCount}</td>
                <td>{s.scoredCount}</td>
                <td>{s.appearedCount > 0 ? ((s.scoredCount / s.appearedCount) * 100).toFixed(1) + "%" : "0%"}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </>
  );

  return (
    <div className="faction-modal-overlay">
      <div className="faction-modal">
        <h4>{faction} - Objective Stats</h4>
        {renderTable("Stage I", stageI)}
        {renderTable("Stage II", stageII)}
        {renderTable("Secret Objectives", secret)}
        <button className="btn btn-secondary" onClick={onClose}>Close</button>
      </div>
    </div>
  );
}
