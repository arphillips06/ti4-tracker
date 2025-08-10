import React, { useState } from "react";
import StandardModal from "../shared/StandardModal";

export default function IncentiveProgramModal({ show, onClose, onSubmit }) {
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
      title="Resolve Incentive Program"
      confirmLabel="Submit"
    >
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
    </StandardModal>
  );
}
