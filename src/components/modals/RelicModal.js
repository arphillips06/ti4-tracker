// src/components/RelicModal.js
import React, { useState } from "react";
import Modal from "react-bootstrap/Modal";
import Button from "react-bootstrap/Button";
import PlayerSelect from "../shared/PlayerSelect";


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
        <PlayerSelect
          players={players}
          value={selectedPlayerId}
          onChange={setSelectedPlayerId}
          label={description}
        />
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={onClose}>Cancel</Button>
        <Button variant="primary" onClick={handleSubmit}>Assign Point</Button>
      </Modal.Footer>
    </Modal>
  );
}
