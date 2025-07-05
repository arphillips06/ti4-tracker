import Modal from 'react-bootstrap/Modal';
import Button from 'react-bootstrap/Button';

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
  players
}) {
  return (
    <Modal show={show} onHide={onClose}>
      <Modal.Header closeButton>
        <Modal.Title>Resolve Mutiny Agenda</Modal.Title>
      </Modal.Header>
      <Modal.Body>
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

        <div className="form-check mb-2">
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

        <div>
          <label className="form-label">Who voted "For"?</label>
          {players.map((p) => (
            <div key={p.player_id} className="form-check">
              <input
                className="form-check-input"
                type="checkbox"
                value={p.player_id}
                id={`mutiny-${p.player_id}`}
                checked={mutinyVotes.includes(p.player_id)}
                disabled={mutinyAbstained}
                onChange={(e) => {
                  if (e.target.checked) {
                    setMutinyVotes([...mutinyVotes, p.player_id]);
                  } else {
                    setMutinyVotes(mutinyVotes.filter((id) => id !== p.player_id));
                  }
                }}
              />
              <label className="form-check-label" htmlFor={`mutiny-${p.player_id}`}>
                {p.name}
              </label>
            </div>
          ))}
        </div>
      </Modal.Body>
      <Modal.Footer>
        <Button variant="secondary" onClick={onClose}>
          Cancel
        </Button>
        <Button variant="primary" onClick={onSubmit}>
          Submit
        </Button>
      </Modal.Footer>
    </Modal>
  );
}
