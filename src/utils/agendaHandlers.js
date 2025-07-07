// src/utils/agendaHandlers.js
const API_URL = "http://localhost:8080";

export async function handleMutinySubmit({ gameId, game, mutinyResult, mutinyAbstained, mutinyVotes, refreshGameState }) {
  await fetch(`${API_URL}/agenda/mutiny`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      game_id: parseInt(gameId),
      round_id: game.current_round_id || 0,
      result: mutinyResult,
      for_votes: mutinyAbstained ? [] : mutinyVotes,
    }),
  });

  await refreshGameState();
}

export async function handleSeedSubmit({ gameId, game, result, refreshGameState }) {
  await fetch(`${API_URL}/agenda/seed`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      game_id: parseInt(gameId),
      round_id: game.current_round_id || 0,
      result,
    }),
  });

  await refreshGameState();
}

export async function handleIncentiveSubmit({ gameId, result, refreshGameState, setAgendaModal }) {
  await fetch(`${API_URL}/agenda/incentive-program`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      game_id: parseInt(gameId),
      outcome: result,
    }),
  });

  await refreshGameState();
  setAgendaModal(null);
}

export async function handlePoliticalCensureSubmit({ gameId, game, playerId, gained, refreshGameState }) {
  await fetch(`${API_URL}/agenda/political-censure`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      game_id: parseInt(gameId),
      round_id: game.current_round_id || 0,
      player_id: playerId,
      gained,
    }),
  });

  await refreshGameState();
}

export async function handleClassifiedSubmit({ gameId, playerId, objectiveId, refreshGameState, setAgendaModal }) {
  await fetch(`${API_URL}/agenda/classified-document-leaks`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      game_id: parseInt(gameId),
      player_id: playerId,
      objective_id: objectiveId,
    }),
  });

  await refreshGameState();
  setAgendaModal(null);
}
