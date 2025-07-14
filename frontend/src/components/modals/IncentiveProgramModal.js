import React, { useState } from 'react';
import Modal from 'react-bootstrap/Modal';
import Button from 'react-bootstrap/Button';

export default function IncentiveProgramModal({ show, onClose, onSubmit }) {
  const [result, setResult] = useState("for");

  return (
    <Modal show={show} onHide={onClose}>
      <Modal.Header closeButton>
        <Modal.Title>Resolve Incentive Program</Modal.Title>
      </Modal.Header>
      <Modal.Body>
        <div className="mb-3">
          <label className="form-label">Outcome</label>
          <select
            className="form-select"
            value={result}
            onChange={(e) => setResult(e.target.value)}
          >
            <option value="for">For (Reveal Stage I)</option>
            <option value="against">Against (Reveal Stage II)</option>
          </select>
        </div>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={onClose}>Cancel</Button>
        <Button variant="primary" onClick={() => onSubmit(result)}>Submit</Button>
      </Modal.Footer>
    </Modal>
  );
}
