import React, { useState } from "react";

export default function ClassifiedDocumentLeaksModal({
    show,
    players,
    secretObjectives,
    scoredSecrets,
    onClose,
    onSubmit,
}) {
    const [selectedPlayerId, setSelectedPlayerId] = useState("");
    const [selectedObjectiveId, setSelectedObjectiveId] = useState("");

    if (!show) return null; // ⛔️ Don't render if not shown

    const handleSubmit = () => {
        if (selectedPlayerId && selectedObjectiveId) {
            onSubmit(parseInt(selectedPlayerId), parseInt(selectedObjectiveId));
        }
    };

    return (
        <div className="modal-backdrop">
            <div className="modal-content p-4 bg-white rounded shadow">
                <h5>Classified Document Leaks</h5>
                <p>
                    Select the player and their scored secret objective to convert into a
                    public one.
                </p>
                <div className="mb-3">
                    <label className="form-label">Player</label>

                    <select
                        className="form-select"
                        value={selectedPlayerId}
                        onChange={(e) => setSelectedPlayerId(e.target.value)}
                    >
                        <option value="">Select player</option>
                        {players.map((p) => (
                            <option key={p.player_id} value={p.player_id}>
                                {p.name}
                            </option>
                        ))}
                    </select>
                </div>

                <div className="mb-3">
                    <label className="form-label">Secret Objective</label>
                    {console.log("Selected Player ID:", selectedPlayerId)}
                    {console.log("Secrets for player:", scoredSecrets?.[parseInt(selectedPlayerId)])}

                    <select
                        className="form-select"
                        value={selectedObjectiveId}
                        onChange={(e) => setSelectedObjectiveId(e.target.value)}
                        disabled={!selectedPlayerId}
                    >
                        <option value="">Select objective</option>
                        {(scoredSecrets[parseInt(selectedPlayerId)] || []).map((obj) => (
                            <option key={obj.id} value={obj.id}>
                                {obj.name}
                            </option>
                        ))}
                    </select>
                </div>

                <div className="d-flex justify-content-end gap-2">
                    <button className="btn btn-secondary" onClick={onClose}>
                        Cancel
                    </button>
                    <button
                        className="btn btn-primary"
                        onClick={handleSubmit}
                        disabled={!selectedPlayerId || !selectedObjectiveId}
                    >
                        Confirm
                    </button>
                </div>
            </div>
        </div>
    );
}
