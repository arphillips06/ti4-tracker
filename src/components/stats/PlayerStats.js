import React from "react";

export default function PlayerStats({ stats }) {
  return (
    <div className="stats-section">
      <h3>Player Win Rates</h3>
      <table className="stats-table">
        <thead>
          <tr>
            <th>Player</th>
            <th>Games Played</th>
            <th>Games Won</th>
            <th>Win %</th>
          </tr>
        </thead>
        <tbody>
          {stats.playerWinRates.map((p, i) => (
            <tr key={i}>
              <td>{p.player || "(Unnamed)"}</td>
              <td>{p.gamesPlayed}</td>
              <td>{p.gamesWon}</td>
              <td>
                {p.gamesPlayed > 0
                  ? ((p.gamesWon / p.gamesPlayed) * 100).toFixed(1)
                  : "0.0"}
                %
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
