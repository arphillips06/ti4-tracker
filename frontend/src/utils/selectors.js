const byType = (scores, type) =>
    (scores || []).filter(s => (s.Type || s.type || "").toLowerCase() === type);

export const isAgendaUsed = (game, title) =>
    byType(game?.AllScores, "agenda").some(s => s.AgendaTitle === title || s.agenda_title === title);

export const isRelicUsed = (game, title) =>
    byType(game?.AllScores, "relic").some(s => s.RelicTitle === title || s.relic_title === title);

export const custodiansScorerId = (game) =>
    byType(game?.AllScores, "mecatol").find(Boolean)?.PlayerID
    ?? byType(game?.AllScores, "mecatol").find(Boolean)?.player_id
    ?? null;

export const hasScoredObjective = (game, playerId, objectiveId) =>
    (game?.AllScores || []).some(s =>
        (s.PlayerID ?? s.player_id) === playerId &&
        (s.ObjectiveID ?? s.objective_id) === objectiveId
    );

export const playerTotalPoints = (game, playerId) =>
    (game?.AllScores || [])
        .filter(s => (s.PlayerID ?? s.player_id) === playerId)
        .reduce((sum, s) => sum + (s.Points ?? s.points ?? 0), 0);

export const winnerPlayerId = (game) =>
    game?.winner_id ?? game?.WinnerID ?? null;

export const isMecatolScoreForPlayer = (s, playerId) =>
    (s?.Type || s?.type)?.toLowerCase() === "mecatol" &&
    (s.PlayerID ?? s.player_id) === playerId;

export const isImperialScoreForPlayer = (s, playerId) =>
    (s?.Type || s?.type)?.toLowerCase() === "imperial" &&
    (s.PlayerID ?? s.player_id) === playerId;

export const isAgendaScoreForPlayer = (s, playerId, title = null) =>
    (s?.Type || s?.type)?.toLowerCase() === "agenda" &&
    (s.PlayerID ?? s.player_id) === playerId &&
    (title ? (s.AgendaTitle || s.agenda_title) === title : true);

export const isAgendaScore = (s, title) =>
    (s?.Type || s?.type)?.toLowerCase() === "agenda" &&
    (s.AgendaTitle || s.agenda_title) === title;

export const isRelicScore = (s, title) =>
    (s?.Type || s?.type)?.toLowerCase() === "relic" &&
    (s.RelicTitle || s.relic_title) === title;