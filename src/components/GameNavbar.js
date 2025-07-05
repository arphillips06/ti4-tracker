import React from "react";

export default function GameNavbar({ mutinyUsed, setShowAgendaModal }) {
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
          </ul>
        </div>
      </div>
    </nav>
  );
}
