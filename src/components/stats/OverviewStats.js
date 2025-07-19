import '../../pages/stats.css';

export default function OverviewStats({ stats }) {
  const averageRounds = stats.averageGameRounds?.toFixed(1) || "—";
  const averagePoints = stats.averagePlayerPoints?.toFixed(2) || "—";
  const mostPlayedFaction = stats.mostPlayedFaction || "—";
  const mostVictoriousFaction = stats.mostVictoriousFaction || "—";
  const totalPlayers = stats.totalUniquePlayers || "—";
  return (
    <div className="d-flex flex-wrap gap-4 mb-4 justify-content-center">
      <div className="stat-card">
        <div className="label">Total Games</div>
        <div>{stats.totalGames ?? "—"}</div>
      </div>
      <div className="stat-card">
        <div className="label">Average Rounds per Game</div>
        <div>{averageRounds}</div>
      </div>
      <div className="stat-card">
        <div className="label">Average Player Score</div>
        <div>{averagePoints}</div>
      </div>
      <div className="stat-card">
        <div className="label">Unique Players</div>
        <div>{totalPlayers}</div>
      </div>
      <div className="stat-card">
        <div className="label">Most Played Faction</div>
        <div>{mostPlayedFaction}</div>
      </div>
      <div className="stat-card">
        <div className="label">Most Victorious Faction</div>
        <div>{mostVictoriousFaction}</div>
      </div>
    </div>
  );
}
