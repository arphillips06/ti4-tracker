import React from "react";
import "./PlayerInputRow.css";
import colorNames from "../../data/colourNames";
import factionImageMap from "../../data/factionIcons";

export default function PlayerInputRow({
  index,
  value,
  onNameChange,
  onFactionChange,
  onColorChange,
  factions,
  selectedFactions,
  selectedColors,
}) {
  const availableFactions = factions.filter(
    (f) => !selectedFactions.includes(f.key) || f.key === value.faction
  );
  const factionIcon = value.faction
    ? `/faction-icons/${factionImageMap[value.faction] || "default.webp"}`
    : null;

  const glowColor = value.color || "transparent";

  return (
    <div
      className="player-card"
      style={{ "--glow-color": glowColor }}
    >
      <div className="player-info">
        <input
          type="text"
          placeholder="Player Name"
          value={value.name}
          onChange={(e) => onNameChange(index, e.target.value)}
          className="player-input top"
        />
        <select
          value={value.faction}
          onChange={(e) => onFactionChange(index, e.target.value)}
          className="player-input bottom"
        >
          <option value="">Select Faction</option>
          {availableFactions.map((f) => (
            <option key={f.key} value={f.key}>
              {f.label}
            </option>
          ))}
        </select>
      </div>

      <div
        className="player-color"
        style={{ backgroundColor: value.color || "#2e2e38" }}
      >
        <select
          value={value.color}
          onChange={(e) => onColorChange(index, e.target.value)}
          className="player-input bottom"
          style={{
            backgroundColor: value.color || "#111122",
            color: value.color === "#ffff00" ? "#000" : "#e0c873",
            border: `1px solid ${value.color || "#e0c87388"}`,
          }}
        >

          <option value="">Colour</option>
          {[
            "#ff3333", "#008000", "#3333ff", "#000000", "#ffff00",
            "#ffa500", "#b300b3", "#ff00ff", "#ffffff"
          ].map((color) => (
            <option
              key={color}
              value={color}
              disabled={selectedColors.includes(color) && color !== value.color}
            >
              {colorNames[color] || color}
            </option>
          ))}
        </select>
      </div>

      <div className="faction-icon-wrapper">
        {factionIcon && (
          <img
            src={factionIcon}
            alt={value.faction}
            className="faction-icon"
          />
        )}
      </div>
    </div>
  );
}