// src/components/VictoryBanner.jsx
import React from "react";

export default function VictoryBanner({ winner, finished, victoryPathSummary }) {
  if (!finished || !winner) return null;

  const renderPathSummary = () => {
    if (!victoryPathSummary) return null;

    const {
      path,
      frequency,
      uniqueness_percent
    } = victoryPathSummary;

    return (
      <div className="victory-extra">
        <p className="victory-subtext">
          Victory Path: {`Stage I:${path.stage1} Stage II:${path.stage2scored} Secret:${path.secrets} Custodians:${path.custodians} Imperial:${path.imperial} Relics:${path.relics} Agenda:${path.agenda} Action Cards:${path.action_card} Support:${path.support}`}
        </p>
        <p className="victory-subtext">
          Seen in {frequency} game{frequency !== 1 ? "s" : ""} – Uniqueness: {uniqueness_percent}%
        </p>
      </div>
    );
  };

  return (
    <div className="victory-banner">
      <div className="glow-bg" />
      <div className="victory-content">
        <h2 className="victory-text">{winner.Name.toUpperCase()} CLAIMS THE THRONE!</h2>
        <p className="victory-subtext">Their dominion over the galaxy is undisputed.</p>
        {renderPathSummary()}
      </div>
    </div>
  );
}
