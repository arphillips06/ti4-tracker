export function normalizeGame(raw) {
  const g = { ...raw };
  g.AllScores = (g.all_scores || g.AllScores || []).map(s => ({
    ...s,
    PlayerID: s.PlayerID ?? s.player_id,
    ObjectiveID: s.ObjectiveID ?? s.objective_id,
    RoundID: s.RoundID ?? s.round_id,
    Type: s.Type ?? s.type,
    Points: s.Points ?? s.points,
    AgendaTitle: s.AgendaTitle ?? s.agenda_title,
    RelicTitle: s.RelicTitle ?? s.relic_title,
  }));
  g.game_players = g.players || g.game_players || [];
  g.winner_id = g.winner_id ?? g.WinnerID ?? null;
  return g;
}
