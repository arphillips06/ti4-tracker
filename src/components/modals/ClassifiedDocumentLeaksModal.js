import React, { useState } from "react";
import PlayerSelect from "../shared/PlayerSelect";
import StandardModal from "../shared/StandardModal";

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

  const handleSubmit = () => {
    if (selectedPlayerId && selectedObjectiveId) {
      onSubmit(parseInt(selectedPlayerId), parseInt(selectedObjectiveId));
      onClose();
    }
  };

  const availableObjectives =
    scoredSecrets[parseInt(selectedPlayerId)] || [];

  return (
    <StandardModal
      show={show}
      onClose={onClose}
      onConfirm={handleSubmit}
      title="Classified Document Leaks"
      confirmDisabled={!selectedPlayerId || !selectedObjectiveId}
      confirmLabel="Confirm"
    >
      <p>
        Select the player and their scored secret objective to convert into a
        public one.
      </p>

      <PlayerSelect
        players={players}
        value={selectedPlayerId}
        onChange={setSelectedPlayerId}
        label="Player"
      />

      <div className="mb-3">
        <label className="form-label">Secret Objective</label>
        <select
          className="form-select"
          value={selectedObjectiveId}
          onChange={(e) => setSelectedObjectiveId(e.target.value)}
          disabled={!selectedPlayerId}
        >
          <option value="">Select objective</option>
          {availableObjectives.map((obj) => (
            <option key={obj.id ?? obj.ID} value={obj.id ?? obj.ID}>
              {obj.name ?? obj.Name}
            </option>
          ))}
        </select>
      </div>
    </StandardModal>
  );
}
