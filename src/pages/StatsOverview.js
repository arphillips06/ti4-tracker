import React, { useEffect, useState } from "react";
import OverviewStats from "../components/stats/OverviewStats";
import PlayerStats from "../components/stats/PlayerStats";
import FactionStats from "../components/stats/FactionStats";
import ObjectiveStats from "../components/stats/ObjectiveStats";
import API_BASE_URL from "../config";
import { Link } from "react-router-dom";
import './stats.css';

export default function StatsPage() {
  const [stats, setStats] = useState(null);
  const [view, setView] = useState("overview");

  useEffect(() => {
    fetch(`${API_BASE_URL}/stats/overview`)
      .then((res) => res.json())
      .then((data) => {
        setStats(data);
      })
      .catch((err) => console.error("Failed to load stats:", err));
  }, []);

  if (!stats) return <div className="p-4">Loading stats...</div>;

  return (
    <div className="p-4">

      <h1 className="mb-4">Twilight Imperium Stats</h1>

      {/* Styled View Switcher */}
      <div className="stats-nav mb-4">
        <Link to="/" className="nav-btn">Home</Link>
        <button
          className={view === "overview" ? "active" : ""}
          onClick={() => setView("overview")}
        >
          Overview
        </button>
        <button
          className={view === "players" ? "active" : ""}
          onClick={() => setView("players")}
        >
          Players
        </button>
        <button
          className={view === "factions" ? "active" : ""}
          onClick={() => setView("factions")}
        >
          Factions
        </button>
        <button
          className={view === "objectives" ? "active" : ""}
          onClick={() => setView("objectives")}
        >
          Objectives
        </button>

      </div>


      {/* Dynamic Content View */}
      {view === "overview" && <OverviewStats stats={stats} />}
      {view === "players" && <PlayerStats stats={stats} />}
      {view === "factions" && <FactionStats stats={stats} />}
      {view === "objectives" && <ObjectiveStats stats={stats} />}
    </div>
  );
}
