import React from "react";
import API_BASE_URL from "../../config";
import { postJSON } from "../../utils/helpers";

export default function SpeakerModal({ show, onClose, players, gameId, refetchGame, roundId }) {
  if (!show) return null;

  const handleAssignSpeaker = async (playerId, isInitial = false) => {
    if (!roundId && roundId !== 0) {
      alert("Please wait — current round not loaded yet.");
      return;
    }

    try {
      await postJSON(`/games/${gameId}/speaker`, {
        player_id: playerId,
        round_id: roundId,
        is_initial: isInitial,
      });

      onClose();
      refetchGame?.();
    } catch (err) {
      console.error("Failed to assign speaker:", err);
      alert("Failed to assign speaker. See console for details.");
    }
  };

  return (

    <div className="modal-backdrop" role="dialog" aria-modal="true" onClick={onClose}>
      <div
        className="modal-content p-4 bg-white rounded shadow"
        style={{ maxWidth: "400px", margin: "10% auto" }}
        onClick={(e) => e.stopPropagation()}
      >
        <h5>Assign Speaker</h5>
        <p>Select a player to become the speaker:</p>
        <ul className="list-group">
          {players.map((p) => (
            <li
              key={p.id ?? p.ID}
              className="list-group-item d-flex justify-content-between align-items-center"
            >
              <span>{p.name ?? p.Player?.Name ?? "Unknown"}</span>
              <button
                className="btn btn-sm btn-primary"
                onClick={() => handleAssignSpeaker(p.id ?? p.ID)}
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
