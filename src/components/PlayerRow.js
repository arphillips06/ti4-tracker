// src/components/PlayerRow.js
import React from "react";
import factionImageMap from "../data/factionIcons";

const PlayerRow = ({ player, points }) => {
  const factionImage = factionImageMap[player.faction];

  return (
    <div className="border-2 rounded-md p-4 mb-4">
      <div className="flex items-center space-x-2">
        {factionImage && (
          <img
            src={`/faction-icons/${factionImage}`}
            alt={player.faction}
            className="w-6 h-6 object-contain"
          />
        )}
        <span className="font-semibold">{player.name}</span>
      </div>
      <p className="italic text-gray-600">{player.faction}</p>
      <p>Points: {points}</p>
    </div>
  );
};

export default PlayerRow;
