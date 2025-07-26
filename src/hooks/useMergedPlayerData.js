export default function useMergedPlayerData(game, sort = true) {
  const scoreMap = new Map(game?.scores?.map((s) => [s.player_id, s]) || []);

  const merged = (game?.game_players || []).map((gp) => {
    const p = gp.Player || {};
    const name = p.Name || "Unknown";
    const faction = gp.Faction || "Unknown Faction";

    const factionKey =
      faction.replace(/^The\s+/i, "")
        .replace(/\s+/g, "")
        .replace(/[^a-zA-Z0-9]/g, "") +
      (faction.toLowerCase().includes("keleres") ? "FactionSymbol" : "");

    return {
      player_id: p.ID,
      id: gp.ID,
      name,
      faction,
      factionKey,
      color: gp.color || "#000",
      points: scoreMap.get(p.ID)?.points || 0,
    };
  });

  return sort ? merged.sort((a, b) => b.points - a.points) : merged;
}
