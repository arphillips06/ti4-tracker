import React from "react";

export default function GameNavbar({
  mutinyUsed,
  incentiveUsed,
  seedUsed,
  setShowAgendaModal,
  setShowCensureModal,
  setShowSeedModal,
  setAgendaModal,
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
                disabled={seedUsed}
              >
                Seed of an Empire {seedUsed ? "(Used)" : ""}
              </button>
            </li>
            <li>
              <button
                className="dropdown-item"
                onClick={() => setAgendaModal("Classified Document Leaks")}
              >
                Classified Document Leaks
              </button>
            </li>
            <li>
              <button
                className="dropdown-item"
                onClick={() => setAgendaModal("Incentive Program")}
                disabled={incentiveUsed}
              >
                Incentive Program {incentiveUsed ? "(Used)" : ""}
              </button>
            </li>
          </ul>
        </div>
      </div>
    </nav>
  );
}
