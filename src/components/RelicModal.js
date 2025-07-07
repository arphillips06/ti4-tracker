// src/components/RelicModal.js
import React, { useState } from "react";
import Modal from "react-bootstrap/Modal";
import Button from "react-bootstrap/Button";

export default function RelicModal({ show, onClose, title, players, onSubmit, description }) {
  const [selectedPlayerId, setSelectedPlayerId] = useState("");

  const handleSubmit = () => {
    if (!selectedPlayerId) return alert("Select a player.");
    onSubmit(selectedPlayerId);
    onClose();
  };

  return (
    <Modal show={show} onHide={onClose} backdrop="static" centered>
      <Modal.Header closeButton>
        <Modal.Title>{title}</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <label className="form-label">{description}</label>
        <select
          className="form-select"
          value={selectedPlayerId}
          onChange={(e) => setSelectedPlayerId(e.target.value)}
        >
          <option value="">-- Select Player --</option>
          {players.map((p) => (
            <option key={p.player_id} value={p.player_id}>{p.name}</option>
          ))}
        </select>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={onClose}>Cancel</Button>
        <Button variant="primary" onClick={handleSubmit}>Assign Point</Button>
      </Modal.Footer>
    </Modal>
  );
}
