import React from "react";
import "../../pages/stats.css";

export default function GameStats({ stats }) {
  if (!stats || !stats.gameLengthStats) return null;

  const {
    longest_by_rounds,
    shortest_by_rounds,
    longest_by_time,
    shortest_by_time,
    average_round_time
  } = stats.gameLengthStats;

  const renderRow = (label, stat) => {
    if (!stat) return null;
    return (
      <tr key={label}>
        <td>{label}</td>
        <td>Game #{stat.game_id}</td>
        <td>{stat.round_count}</td>
        <td>{stat.duration}</td>
      </tr>
    );
  };

  return (
    <div className="stats-section">
      <h2 className="chart-header">Game Length Stats</h2>

      <table className="stats-table">
        <thead>
          <tr>
            <th>Stat</th>
            <th>Game</th>
            <th>Rounds</th>
            <th>Duration</th>
          </tr>
        </thead>
        <tbody>
          {renderRow("Longest Game (Rounds)", longest_by_rounds)}
          {renderRow("Shortest Game (Rounds)", shortest_by_rounds)}
          {renderRow("Longest Game (Time)", longest_by_time)}
          {renderRow("Shortest Game (Time)", shortest_by_time)}
          <tr>
            <td>Average Round Length</td>
            <td colSpan="3">{average_round_time}</td>
          </tr>
        </tbody>
      </table>
    </div>
  );
}
