import React from "react";
import {
  ResponsiveContainer,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  Legend,
} from "recharts";
import "./shared/graphs.css";

export default function FactionWinRateChart({ dataMap }) {
  if (!dataMap) return null;

  const data = Object.entries(dataMap).map(([faction, winRate]) => ({
    faction,
    winRate,
  }));

  return (
    <div className="graph-container">
      <h3 className="chart-section-title">Faction Win Rate</h3>
      <ResponsiveContainer width="100%" height={400}>
        <BarChart data={data} margin={{ top: 20, right: 30, left: 20, bottom: 80 }}>
          <XAxis
            dataKey="faction"
            angle={-45}
            textAnchor="end"
            interval={0}
            tick={{ fill: "#fff", fontWeight: 500, fontSize: 12 }}
            label={{
              value: "Faction",
              position: "insideBottom",
              offset: -80,
              fill: "#fff",
              fontWeight: 600,
            }}
          />

          <YAxis
            domain={[0, 100]}
            tick={{ fill: "#fff", fontWeight: 500 }}
            label={{
              value: "Win Rate (%)",
              angle: -90,
              position: "insideLeft",
              fill: "#fff",
              fontWeight: 600,
            }}
          />
          <Tooltip
            cursor={false}
            content={({ active, payload, label }) =>
              active && payload?.length ? (
                <div
                  style={{
                    backgroundColor: "rgba(0, 0, 0, 0.85)",
                    padding: "8px",
                    color: "#ffd700",
                    border: "1px solid #ffd700",
                    borderRadius: "4px",
                  }}
                >
                  <strong>{label}</strong>
                  <br />
                  Win Rate: {payload[0].value.toFixed(2)}%
                </div>
              ) : null
            }
          />
          <Bar
            dataKey="winRate"
            name="Win Rate (%)"
            shape={({ x, y, width, height }) => (
              <rect
                x={x}
                y={y}
                width={width}
                height={height}
                fill="#61dafb"
                style={{ pointerEvents: "none" }} // disables highlight behavior
              />
            )}
            isAnimationActive={false}
          />


        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}