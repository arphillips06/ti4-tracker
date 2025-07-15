// src/components/shared/PlayerSelect.js
import React from "react";

export default function PlayerSelect({
  players = [],
  value,
  onChange,
  label = "Select a player",
  placeholder = "-- Select Player --",
  disabledIds = [],
  id = "player-select",
}) {
  return (
    <div className="mb-3">
      {label && <label htmlFor={id} className="form-label">{label}</label>}
      <select
        id={id}
        className="form-select"
        value={value || ""}
        onChange={(e) => onChange(e.target.value)}
      >
        <option value="">{placeholder}</option>
        {players.map((p) => (
          <option
            key={p.player_id}
            value={p.player_id}
            disabled={disabledIds.includes(p.player_id)}
          >
            {p.name}
          </option>
        ))}
      </select>
    </div>
  );
}
