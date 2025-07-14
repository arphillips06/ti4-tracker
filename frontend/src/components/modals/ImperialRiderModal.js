import React, { useState } from "react";
import Modal from "react-bootstrap/Modal";
import Button from "react-bootstrap/Button";

export default function ImperialRiderModal({ show, onClose, players, onSubmit }) {
  const [selectedPlayerId, setSelectedPlayerId] = useState("");

  const handleSubmit = () => {
    if (selectedPlayerId) {
      onSubmit(selectedPlayerId);
      onClose();
    } else {
      alert("Please select a player.");
    }
  };

  return (
    <Modal show={show} onHide={onClose} backdrop="static" centered>
      <Modal.Header closeButton>
        <Modal.Title>Score Imperial Rider</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <label htmlFor="playerSelect" className="form-label">Select player to assign point</label>
        <select
          id="playerSelect"
          className="form-select"
          value={selectedPlayerId}
          onChange={(e) => setSelectedPlayerId(e.target.value)}
        >
          <option value="">-- Select player --</option>
          {players.map((p) => (
            <option key={p.player_id} value={p.player_id}>
              {p.name}
            </option>
          ))}
        </select>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={onClose}>
          Cancel
        </Button>
        <Button variant="primary" onClick={handleSubmit}>
          Assign Point
        </Button>
      </Modal.Footer>
    </Modal>
  );
}
