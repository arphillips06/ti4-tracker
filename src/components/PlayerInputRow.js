// src/components/PlayerRow.js
import React from "react";


const COLORS = [
  "#ff3333", // Red
  "#3333ff", // Blue
  "#008000", // Green
  "#ffff00", // Yellow
  "#b300b3", // Purple
  "#000000", // Black
  "#ffa500", // Orange
  "#ff00ff", // Pink
  "#ffffff", // White
];

const colorName = (hex) => {
  const names = {
    "#ff3333": "Red",
    "#3333ff": "Blue",
    "#008000": "Green",
    "#ffff00": "Yellow",
    "#b300b3": "Purple",
    "#000000": "Black",
    "#ffa500": "Orange",
    "#ff00ff": "Pink",
    "#ffffff": "White",
  };
  return names[hex] || hex;
};

export default function PlayerRow({
  index,
  value,
  onFactionChange,
  onColorChange,
  onNameChange,
  factions,
  selectedFactions,
  selectedColors,
}) {
  return (
    <div className="flex space-x-4 mb-2 items-center">
      <input
        type="text"
        placeholder={`Player ${index + 1}`}
        className="border p-2 rounded w-1/4"
        value={value.name}
        onChange={(e) => onNameChange(index, e.target.value)}
      />
      <select
        className="border p-2 rounded w-1/2"
        value={value.faction || ""}
        onChange={(e) => onFactionChange(index, e.target.value)}
      >
        <option value="">Select Faction</option>
        {factions.map((faction) => (
          <option
            key={faction.key}
            value={faction.key}
            disabled={
              selectedFactions.includes(faction.key) &&
              faction.key !== value.faction
            }
          >
            {faction.label}
          </option>
        ))}
      </select>
      <select
        className="border p-2 rounded w-1/4"
        value={value.color || ""}
        onChange={(e) => onColorChange(index, e.target.value)}
      >
        <option value="">Color</option>
        {COLORS.map((color) => (
          <option
            key={color}
            value={color}
            disabled={selectedColors.includes(color) && color !== value.color}
            style={{
              backgroundColor: color,
              color: color === "#ffffff" ? "#000000" : "#ffffff",
            }}
          >
            {colorName(color)}
          </option>
        ))}
      </select>
    </div>
  );
}
