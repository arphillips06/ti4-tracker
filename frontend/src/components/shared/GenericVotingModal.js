import React, { useState, useEffect } from "react";
import StandardModal from "./StandardModal";

export default function GenericVotingModal({
  title,
  players,
  onSubmit,
  onCancel,
  show, // 👈 add this
  showAbstain = false,
  showDirectionToggle = false,
  defaultDirection = "for",
}) {
    const [selectedPlayerId, setSelectedPlayerId] = useState("");
    const [abstain, setAbstain] = useState(false);
    const [direction, setDirection] = useState(defaultDirection);

    const handleSubmit = () => {
        console.log("🧪 handleSubmit called");

        // Prevent submission if no player selected and not abstaining
        if (!selectedPlayerId && !abstain) {
            console.warn("🚫 No player selected and not abstaining");
            return;
        }

        onSubmit({
            playerId: parseInt(selectedPlayerId, 10), // ✅ force number
            abstain,
            direction,
        });
    };

    return (
        <StandardModal show={true} onClose={onCancel}>
            <h5 className="mb-3">{title}</h5>

            {!abstain && (
                <>
                    <label className="form-label small">Select Player</label>
                    <select
                        className="form-select form-select-sm mb-3"
                        value={selectedPlayerId}
                        onChange={(e) => setSelectedPlayerId(e.target.value)}
                    >
                        <option value="">-- Choose a player --</option>
                        {players.map((p) => (
                            <option key={p.player_id} value={p.player_id}>
                                {p.name}
                            </option>
                        ))}
                    </select>
                </>
            )}

            {showDirectionToggle && (
                <div className="mb-3">
                    <label className="form-label small">Vote Direction</label>
                    <select
                        className="form-select form-select-sm"
                        value={direction}
                        onChange={(e) => setDirection(e.target.value)}
                    >
                        <option value="for">For</option>
                        <option value="against">Against</option>
                    </select>
                </div>
            )}

            {showAbstain && (
                <div className="form-check mb-3">
                    <input
                        type="checkbox"
                        className="form-check-input"
                        id="abstainCheck"
                        checked={abstain}
                        onChange={() => setAbstain((prev) => !prev)}
                    />
                    <label className="form-check-label small" htmlFor="abstainCheck">
                        Abstain
                    </label>
                </div>
            )}

            <div className="d-flex gap-2">
                <button className="btn btn-primary btn-sm" onClick={handleSubmit}>
                    Submit
                </button>
                <button className="btn btn-secondary btn-sm" onClick={onCancel}>
                    Cancel
                </button>
            </div>
        </StandardModal>
    );
}
