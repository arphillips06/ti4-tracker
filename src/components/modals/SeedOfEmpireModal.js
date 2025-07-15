import React, { useState } from "react";
import StandardModal from "../shared/StandardModal";

export default function SeedOfEmpireModal({ show, onClose, onSubmit }) {
  const [result, setResult] = useState("for");

  const handleSubmit = () => {
    onSubmit(result);
    onClose();
  };

  return (
    <StandardModal
      show={show}
      onClose={onClose}
      onConfirm={handleSubmit}
      title="Seed of an Empire"
      confirmLabel="Confirm"
    >
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
    </StandardModal>
  );
}
