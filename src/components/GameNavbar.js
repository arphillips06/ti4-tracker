import React from "react";

export default function GameNavbar({ 
  mutinyUsed, 
  setShowAgendaModal, 
  setShowCensureModal, 
  setShowSeedModal,
  setAgendaModal
 }) {
  return (
    <nav className="navbar navbar-expand-lg navbar-dark bg-dark mb-4">
      <div className="container-fluid">
        <span className="navbar-brand">Game Controls</span>
        <div className="dropdown">
          <button
            className="btn btn-outline-light dropdown-toggle"
            type="button"
            data-bs-toggle="dropdown"
          >
            Agendas
          </button>
          <ul className="dropdown-menu">
            <li>
              <button
                className="dropdown-item"
                onClick={() => setShowAgendaModal(true)}
                disabled={mutinyUsed}
              >
                Mutiny {mutinyUsed ? "(Used)" : ""}
              </button>
            </li>
            <li>
              <button
                className="dropdown-item"
                onClick={() => setShowCensureModal(true)}
              >
                Political Censure
              </button>
              <button
                className="dropdown-item"
                onClick={() => setShowSeedModal(true)}
              >
                Seed of an Empire
              </button>
            </li>
            <li>
              <button className="dropdown-item" onClick={() => setAgendaModal("Classified Document Leaks")}>
                Classified Document Leaks
              </button>
            </li>
          </ul>
        </div>
      </div>
    </nav>
  );
}
