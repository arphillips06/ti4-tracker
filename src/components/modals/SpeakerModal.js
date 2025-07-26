import React from "react";
import API_BASE_URL from "../../config";

export default function SpeakerModal({ show, onClose, players, gameId, refetchGame, roundId }) {
  if (!show) return null;

  const handleAssignSpeaker = async (playerId, isInitial = false) => {

    if (!roundId && roundId !== 0) {
      alert("Please wait — current round not loaded yet.");
      return;
    }

    await fetch(`${API_BASE_URL}/games/${gameId}/speaker`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        player_id: playerId,
        round_id: roundId,
        is_initial: isInitial,
      }),
    });

    onClose();
    if (refetchGame) refetchGame();
  };

  return (

    <div className="modal-backdrop">
      <div className="modal-content p-4 bg-white rounded shadow" style={{ maxWidth: "400px", margin: "10% auto" }}>
        <h5>Assign Speaker</h5>
        <p>Select a player to become the speaker:</p>
        <ul className="list-group">
          {players.map((p) => (
            <li key={p.ID} className="list-group-item d-flex justify-content-between align-items-center">
              <span>{p.Player?.Name || 'Unknown'}</span>
              <button
                className="btn btn-sm btn-primary"
                onClick={() => {
                  handleAssignSpeaker(p.ID);
                }}
              >
                Assign
              </button>
            </li>
          ))}
        </ul>

        <div className="mt-3 text-end">
          <button className="btn btn-sm btn-secondary" onClick={onClose}>
            Cancel
          </button>
        </div>
      </div>
    </div>
  );
}
