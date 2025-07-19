// src/components/VictoryBanner.jsx
import React from "react";

export default function VictoryBanner({ winner, finished }) {
  if (!finished || !winner) return null;

  return (
    <div className="victory-banner">
      <div className="glow-bg" />
      <div className="victory-content">
        <h2 className="victory-text">{winner.Name.toUpperCase()} CLAIMS THE THRONE!</h2>
        <p className="victory-subtext">Their dominion over the galaxy is undisputed.</p>
      </div>
    </div>
  );
}
