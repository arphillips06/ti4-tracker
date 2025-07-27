import React from "react";
import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  Tooltip,
  ResponsiveContainer,
  Cell,
} from "recharts";

export default function VictoryPointBarChart({ vpHistogram }) {
  if (!vpHistogram || vpHistogram.length === 0) return null;

  const highlightColor = "#ffd700"; // highlight 10 VP
  const barColor = "#bbbbbb";

  // Fill 1–14 even if empty
  const filledHistogram = Array.from({ length: 14 }, (_, i) => {
    const vp = i + 1;
    const entry = vpHistogram.find((e) => e.vp === vp);
    return { vp, count: entry ? entry.count : 0 };
  });

  return (
    <div style={{ width: "180px", height: "80px", marginLeft: "auto" }}>
      <ResponsiveContainer width="100%" height="100%">
        <BarChart
          data={filledHistogram}
          margin={{ top: 0, right: 0, bottom: 10, left: 0 }}
        >
          <XAxis
            dataKey="vp"
            tick={{ fill: "#ccc", fontSize: 10 }}
            axisLine={false}
            tickLine={false}
          />
          <YAxis hide />
          <Tooltip
            contentStyle={{
              backgroundColor: "#222",
              border: "none",
              color: "#fff",
            }}
            formatter={(value) => `${value} game${value === 1 ? "" : "s"}`}
            labelFormatter={(label) => `VP: ${label}`}
          />
          <Bar dataKey="count" radius={[3, 3, 0, 0]}>
            {filledHistogram.map((entry, index) => (
              <Cell
                key={`bar-${index}`}
                fill={entry.vp === 10 ? highlightColor : barColor}
              />
            ))}
          </Bar>
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
