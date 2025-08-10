import React, { useState } from "react";
import Modal from "react-bootstrap/Modal";
import Button from "react-bootstrap/Button";
import PlayerSelect from "../shared/PlayerSelect";


export default function PoliticalCensureModal({ show, onClose, onSubmit, players }) {
    const scoredEntry = players.find((p) =>
        p.agendaScores?.some((s) => s.AgendaTitle === "Political Censure")
    );
    const alreadyAssignedPlayerId = scoredEntry?.player_id || null;

    const [selectedPlayerId, setSelectedPlayerId] = useState(alreadyAssignedPlayerId);
    const [gained, setGained] = useState(alreadyAssignedPlayerId ? false : true);

    const handleSubmit = () => {
        if (!selectedPlayerId) {
            alert("Please select a player.");
            return;
        }

        onSubmit({ playerId: selectedPlayerId, gained });
        onClose();
    };

    return (
        <Modal show={show} onHide={onClose} backdrop="static" centered>
            <Modal.Header closeButton>
                <Modal.Title>Political Censure</Modal.Title>
            </Modal.Header>
            <Modal.Body>
                <div className="mb-3">
                    <PlayerSelect
                        players={players}
                        value={selectedPlayerId}
                        onChange={(val) => setSelectedPlayerId(parseInt(val))}
                        label="Player Affected"
                        disabledIds={
                            alreadyAssignedPlayerId !== null
                                ? players.map((p) => p.player_id).filter((id) => id !== alreadyAssignedPlayerId)
                                : []
                        }
                    />
                </div>
                <div className="form-check mb-2">
                    <input
                        className="form-check-input"
                        type="radio"
                        name="censureDirection"
                        id="gain"
                        checked={gained}
                        onChange={() => setGained(true)}
                        disabled={alreadyAssignedPlayerId !== null}
                    />
                    <label className="form-check-label" htmlFor="gain">
                        Gain VP (Received Political Censure)
                    </label>
                </div>
                <div className="form-check">
                    <input
                        className="form-check-input"
                        type="radio"
                        name="censureDirection"
                        id="lose"
                        checked={!gained}
                        onChange={() => setGained(false)}
                    />
                    <label className="form-check-label" htmlFor="lose">
                        Lose VP (Lost Political Censure)
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
