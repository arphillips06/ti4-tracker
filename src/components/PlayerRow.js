const COLORS = [
  { name: "Green", value: "#008000" },
  { name: "Red", value: "#ff3333" },
  { name: "Blue", value: "#3333ff" },
  { name: "Black", value: "#000000" },
  { name: "Orange", value: "#ffa500" },
  { name: "Yellow", value: "#ffff00" },
  { name: "Pink", value: "#ff00ff" },
  { name: "Purple", value: "#b300b3" },
  { name: "White", value: "#ffffff" },
];

export default function PlayerRow({ index, player, onChange, factions }) {
  return (
    <div className="flex gap-4 mb-2 items-center">
      {/* Player Name */}
      <input
        type="text"
        className="border p-2 rounded flex-1"
        placeholder={`Player ${index + 1}`}
        value={player.name}
        onChange={(e) => onChange(index, "name", e.target.value)}
      />

      {/* Faction Dropdown */}
      <select
        className="border p-2 rounded"
        value={player.faction}
        onChange={(e) => onChange(index, "faction", e.target.value)}
      >
        <option value="">Select Faction</option>
        {factions.map(f => (
          <option key={f} value={f}>{f}</option>
        ))}
      </select>

      {/* Color Dropdown */}
      <select
        className="border p-2 rounded"
        value={player.color}
        onChange={(e) => onChange(index, "color", e.target.value)}
        style={{ backgroundColor: player.color, color: player.color === "#ffff00" ? "#000" : "#fff" }}
      >
        {COLORS.map(c => (
          <option
            key={c.value}
            value={c.value}
            style={{ backgroundColor: c.value, color: c.value === "#ffff00" ? "#000" : "#fff" }}
          >
            {c.name}
          </option>
        ))}
      </select>
    </div>
  );
}
