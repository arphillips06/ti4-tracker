import React, { useEffect, useState, useCallback } from "react";
import axios from "axios";
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";
import "./shared/graphs.css";
import API_BASE_URL from "../../config";

export default function ScoreGraph({ gameId, refreshSignal }) {
  const [data, setData] = useState([]);
  const [players, setPlayers] = useState([]);

  const fetchScores = useCallback(async () => {
    try {
      const res = await axios.get(`${API_BASE_URL}/games/${gameId}/scores-by-round`);
      const raw = res.data;
      if (!Array.isArray(raw)) {
        console.error("Invalid score data format:", raw);
        return;
      }

      const cumulativeScores = {};
      const playerNames = new Set();

      // Collect player names
      raw.forEach(({ scores }) =>
        scores.forEach(({ player }) => playerNames.add(player))
      );

      for (const name of playerNames) {
        cumulativeScores[name] = 0;
      }

      const chartData = [];
      const zeroEntry = { round: 0 };
      playerNames.forEach((name) => (zeroEntry[name] = 0));
      chartData.push(zeroEntry);

      raw.forEach(({ round, scores }) => {
        const entry = { round };
        for (const name of playerNames) {
          entry[name] = cumulativeScores[name];
        }
        scores.forEach(({ player, points }) => {
          cumulativeScores[player] += points;
          entry[player] = cumulativeScores[player];
        });
        chartData.push(entry);
      });

      setPlayers([...playerNames]);
      setData(chartData);
    } catch (err) {
      console.error("Failed to fetch round scores:", err);
    }
  }, [gameId]);

  useEffect(() => {
    fetchScores();
  }, [fetchScores, refreshSignal]);

  return (
    <div className="graph-container">
      <h3 className="chart-section-title">Score Progression</h3>
      <div style={{ height: "400px" }}>
        <ResponsiveContainer width="100%" height="100%">
          <LineChart data={data}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis
              dataKey="round"
              type="number"
              domain={["dataMin", "dataMax"]}
              tickFormatter={(tick) => Math.round(tick)}
              allowDecimals={false}
              label={{ value: "Round", position: "insideBottomRight", offset: -5 }}
            />
            <YAxis
              allowDecimals={false}
              label={{ value: "Score", angle: -90, position: "insideLeft" }}
            />
            <Tooltip />
            <Legend />
            {players.map((player, idx) => (
              <Line
                key={player}
                type="monotone"
                dataKey={player}
                stroke={`hsl(${(idx * 60) % 360}, 70%, 50%)`}
                dot={{ r: 3 }}
              />
            ))}
          </LineChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}
