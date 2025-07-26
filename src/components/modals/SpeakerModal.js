import React from "react";

export default function SpeakerModal({ show, onClose, players, onAssign }) {
  if (!show) return null;

  return (
    <div className="modal-backdrop">
      <div className="modal-content p-4 bg-white rounded shadow" style={{ maxWidth: "400px", margin: "10% auto" }}>
        <h5>Assign Speaker</h5>
        <p>Select a player to become the speaker:</p>
        <ul className="list-group">
          {players.map((p) => (
            <li key={p.player_id} className="list-group-item d-flex justify-content-between align-items-center">
              <span>{p.name}</span>
              <button
                className="btn btn-sm btn-primary"
                onClick={() => {
                  onAssign(p.player_id);
                  onClose();
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
