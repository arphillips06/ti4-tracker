// src/utils/agendaHandlers.js
import { postJSON, submitAndRefresh } from "./helpers";

export async function handleMutinySubmit({ gameId, game, mutinyResult, mutinyAbstained, mutinyVotes, refreshGameState }) {
  await submitAndRefresh({
    requestFn: () =>
      postJSON("/agenda/mutiny", {
        game_id: parseInt(gameId),
        round_id: game.current_round_id || 0,
        result: mutinyResult,
        for_votes: mutinyAbstained ? [] : mutinyVotes,
      }),
    refreshGameState,
  });
}

export async function handleSeedSubmit({ gameId, game, result, refreshGameState }) {
  await submitAndRefresh({
    requestFn: () =>
      postJSON("/agenda/seed", {
        game_id: parseInt(gameId),
        round_id: game.current_round_id || 0,
        result,
      }),
    refreshGameState,
  });
}

export async function handleIncentiveSubmit({ gameId, result, refreshGameState, setAgendaModal }) {
  await submitAndRefresh({
    requestFn: () =>
      postJSON("/agenda/incentive-program", {
        game_id: parseInt(gameId),
        outcome: result,
      }),
    refreshGameState,
    closeModal: () => setAgendaModal(null),
  });
}

export async function handlePoliticalCensureSubmit({ gameId, game, playerId, gained, refreshGameState }) {
  await submitAndRefresh({
    requestFn: () =>
      postJSON("/agenda/political-censure", {
        game_id: parseInt(gameId),
        round_id: game.current_round_id || 0,
        player_id: playerId,
        gained,
      }),
    refreshGameState,
  });
}

export async function handleClassifiedSubmit({ gameId, playerId, objectiveId, refreshGameState, setAgendaModal }) {
  await submitAndRefresh({
    requestFn: () =>
      postJSON("/agenda/classified-document-leaks", {
        game_id: parseInt(gameId),
        player_id: playerId,
        objective_id: objectiveId,
      }),
    refreshGameState,
    closeModal: () => setAgendaModal(null),
  });
}
