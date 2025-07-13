import React from "react";

export default function PlayerStats({ stats }) {
  const avgPoints = stats.playerAveragePoints || [];
  const factionData = stats.topFactionsPerPlayer || [];
  const finishes = stats.playerMostCommonFinishes || [];
  const stdevMap = {};
  (stats.playerPointStdevs || []).forEach((entry) => {
    stdevMap[entry.player] = entry.stdev;
  });

  // Map for faction data
  const factionMap = {};
  factionData.forEach((entry) => {
    factionMap[entry.player] = entry.factions;
  });
  // Map standard deviation by player for quick lookup
  (stats.playerPointStdevs || []).forEach((entry) => {
    stdevMap[entry.player] = entry.stdev;
  });

  // Map for most common finishes
  const finishMap = {};
  finishes.forEach((entry) => {
    finishMap[entry.player] = entry;
  });

  const sorted = [...avgPoints].sort((a, b) => b.averagePoints - a.averagePoints);

  return (
    <div className="stats-section">
      <h2>Player Average Points per Game</h2>

      <table className="stats-table">

        <thead>
          <tr>
            <th>Player</th>
            <th>Games Played</th>
            <th>Total Points</th>
            <th>Average Points</th>
            <th>Stdev</th>
            <th>Top Factions</th>
            <th>Most Common Finish</th>
          </tr>
        </thead>
        <tbody>
          {sorted.map((p) => {
            const factions = factionMap[p.player] || {};
            const topFactions = Object.entries(factions)
              .sort((a, b) => b[1] - a[1])
              .slice(0, 3); // Top 3 factions

            const finish = finishMap[p.player];
            const finishDisplay = finish
              ? `${finish.position} (${finish.count}/${finish.totalGames}, ${((finish.count / finish.totalGames) * 100).toFixed(1)}%)`
              : "—";

            return (
              <tr key={p.player}>
                <td>{p.player || "(Unnamed)"}</td>
                <td>{p.gamesPlayed}</td>
                <td>{p.totalPoints.toFixed(1)}</td>
                <td>{p.averagePoints.toFixed(2)}</td>
                <td>{stdevMap[p.player]?.toFixed(2) || "0.00"}</td>

                <td>
                  {topFactions.map(([faction, count]) => (
                    <div key={faction}>
                      {faction} ({count})
                    </div>
                  ))}
                </td>
                <td>{finishDisplay}</td>
              </tr>
            );
          })}
        </tbody>
      </table>
      <h2 className="mt-5">Secret Objective Scoring Rate</h2>
      <table className="stats-table">
        <thead>
          <tr>
            <th>Player</th>
            <th>Secret Objectives Possible</th>
            <th>Secret Objectives Scored</th>
            <th>Score Rate (%)</th>
          </tr>
        </thead>
        <tbody>
          {stats.secretObjectiveRates
            .sort((a, b) => b.secretScoreRate - a.secretScoreRate)
            .map((p) => (
              <tr key={p.player}>
                <td>{p.player}</td>
                <td>{p.secretAppeared}</td>
                <td>{p.secretScored}</td>
                <td>{p.secretScoreRate.toFixed(1)}%</td>
              </tr>
            ))}
        </tbody>
      </table>

    </div>
  );
}
