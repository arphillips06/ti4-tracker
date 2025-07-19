import React from "react";
import StandardModal from "../shared/StandardModal";
import "./modal.css";

export default function MutinyModal({
  show,
  onClose,
  onSubmit,
  mutinyResult,
  setMutinyResult,
  mutinyAbstained,
  setMutinyAbstained,
  mutinyVotes,
  setMutinyVotes,
  players,
}) {
  const handleSubmit = async () => {
    await onSubmit();
    onClose();
  };

  return (
    <StandardModal
      show={show}
      onClose={onClose}
      onConfirm={handleSubmit}
      title="Resolve Mutiny Agenda"
      confirmLabel="Submit"
    >
      <div className="mb-3">
        <label className="form-label">How did the agenda resolve?</label>
        <select
          className="form-select"
          value={mutinyResult}
          onChange={(e) => setMutinyResult(e.target.value)}
        >
          <option value="for">For</option>
          <option value="against">Against</option>
        </select>
      </div>

      <div className="form-check mb-3">
        <input
          className="form-check-input"
          type="checkbox"
          id="mutiny-abstain"
          checked={mutinyAbstained}
          onChange={() => setMutinyAbstained(!mutinyAbstained)}
        />
        <label className="form-check-label" htmlFor="mutiny-abstain">
          All players abstained
        </label>
      </div>

      {!mutinyAbstained && (
        <fieldset className="mt-3">
          <legend className="form-label">Who voted "For"?</legend>
          {players.map((p) => (
            <div key={p.player_id} className="form-check mb-2">
              <input
                className="form-check-input"
                type="checkbox"
                id={`mutiny-${p.player_id}`}
                checked={mutinyVotes.includes(p.player_id)}
                onChange={(e) => {
                  if (e.target.checked) {
                    setMutinyVotes([...mutinyVotes, p.player_id]);
                  } else {
                    setMutinyVotes(
                      mutinyVotes.filter((id) => id !== p.player_id)
                    );
                  }
                }}
              />
              <label
                className="form-check-label"
                htmlFor={`mutiny-${p.player_id}`}
              >
                {p.name}
              </label>
            </div>
          ))}
        </fieldset>
      )}
    </StandardModal>
  );
}
