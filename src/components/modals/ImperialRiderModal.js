import React, { useState } from "react";
import Modal from "react-bootstrap/Modal";
import Button from "react-bootstrap/Button";
import PlayerSelect from "../shared/PlayerSelect";
import StandardModal from "../shared/StandardModal";

export default function ImperialRiderModal({ show, onClose, players, onSubmit }) {
  const [selectedPlayerId, setSelectedPlayerId] = useState("");
  const [selectedOutcome, setSelectedOutcome] = useState("");

const handleSubmit = () => {
  if (selectedPlayerId && selectedOutcome) {
    onSubmit(selectedPlayerId, selectedOutcome);
    onClose();
  } else {
    alert("Please select both a player and an outcome.");
  }
};

  return (
    <StandardModal
      show={show}
      onClose={onClose}
      onConfirm={handleSubmit}
      title="Imperial Rider"
      confirmLabel="Submit"
      confirmDisabled={!selectedPlayerId || !selectedOutcome}
    >
      <div className="mb-3">
        <label className="form-label">Prediction</label>
        <select
          className="form-select"
          value={selectedOutcome}
          onChange={(e) => setSelectedOutcome(e.target.value)}
        >
          <option value="">-- Choose Outcome --</option>
          <option value="for">For</option>
          <option value="against">Against</option>
        </select>
      </div>

      <PlayerSelect
        players={players}
        value={selectedPlayerId}
        onChange={setSelectedPlayerId}
        label="Select player to assign point"
      />
    </StandardModal>
  );
}
