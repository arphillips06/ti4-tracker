import React, { useState } from "react";
import Modal from "react-bootstrap/Modal";
import Button from "react-bootstrap/Button";

export default function SeedOfEmpireModal({ show, onClose, onSubmit }) {
  const [result, setResult] = useState("for");

  const handleSubmit = () => {
    onSubmit(result);
    onClose();
  };

  return (
    <Modal show={show} onHide={onClose} backdrop="static" centered>
      <Modal.Header closeButton>
        <Modal.Title>Seed of an Empire</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <div className="form-check mb-2">
          <input
            className="form-check-input"
            type="radio"
            name="seedResult"
            id="for"
            checked={result === "for"}
            onChange={() => setResult("for")}
          />
          <label className="form-check-label" htmlFor="for">
            For (Most points gains 1)
          </label>
        </div>
        <div className="form-check">
          <input
            className="form-check-input"
            type="radio"
            name="seedResult"
            id="against"
            checked={result === "against"}
            onChange={() => setResult("against")}
          />
          <label className="form-check-label" htmlFor="against">
            Against (Fewest points gains 1)
          </label>
        </div>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={onClose}>
          Cancel
        </Button>
        <Button variant="primary" onClick={handleSubmit}>
          Confirm
        </Button>
      </Modal.Footer>
    </Modal>
  );
}
