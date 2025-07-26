import React from "react";
import { useNavigate } from "react-router-dom";
import "../layout/navbar.css";


export default function GameNavbar({
  gameId,
  showScoreGraph,
  setShowScoreGraph,
  mutinyUsed,
  incentiveUsed,
  seedUsed,
  setShowAgendaModal,
  setShowCensureModal,
  setShowSeedModal,
  setAgendaModal,
  setShowImperialModal,
  setShowCrownModal,
  setShowShardModal,
  crownUsed,
  obsidianUsed,
  setShowObsidianModal,
  setShowSpeakerModal,
}) {
  const navigate = useNavigate();

  return (
    <nav className="navbar navbar-expand-lg navbar-dark bg-dark mb-4 px-3">
      <div className="container-fluid justify-content-between">
        {/* Left Side */}
        <div className="d-flex align-items-center gap-3">
          <button className="btn btn-sm btn-secondary" onClick={() => navigate("/")}>
            Home
          </button>
          <span className="navbar-text text-light fw-bold">Game #{gameId}</span>
          <button
            className={`btn ${showScoreGraph ? "btn-secondary" : "btn-outline-info"} ms-2`}
            onClick={() => setShowScoreGraph(!showScoreGraph)}
          >
            {showScoreGraph ? "Back to Objectives" : "Score"}
          </button>
        </div>

        {/* Right Side: Button Group */}
        <div className="d-flex align-items-center gap-2">
          <button
            className="btn btn-outline-danger"
            onClick={() => setShowSpeakerModal(true)}
          >
            Assign Speaker
          </button>

          {/* Imperial Rider */}
          <button
            className="btn btn-outline-orange"
            onClick={() => setShowImperialModal(true)}
          >
            Score Imperial Rider
          </button>

          {/* Relics */}
          <div className="btn-group">
            <button
              className="btn btn-outline-warning dropdown-toggle"
              type="button"
              data-bs-toggle="dropdown"
            >
              Relics
            </button>
            <ul className="dropdown-menu custom-dropdown dropdown-menu-end">
              <li>
                <button
                  className="dropdown-item"
                  onClick={() => setShowCrownModal(true)}
                  disabled={crownUsed}
                >
                  The Crown of Emphidia {crownUsed ? "(Used)" : ""}
                </button>
              </li>
              <li>
                <button className="dropdown-item" onClick={() => setShowShardModal(true)}>
                  Shard of the Throne
                </button>
              </li>
              <li>
                <button
                  className="dropdown-item"
                  onClick={() => setShowObsidianModal(true)}
                  disabled={obsidianUsed}
                >
                  The Obsidian {obsidianUsed ? "(Used)" : ""}
                </button>
              </li>
            </ul>
          </div>

          {/* Agendas */}
          <div className="btn-group">
            <button
              className="btn btn-outline-primary dropdown-toggle"
              type="button"
              data-bs-toggle="dropdown"
            >
              Agendas
            </button>
            <ul className="dropdown-menu custom-dropdown dropdown-menu-end">
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
                <button className="dropdown-item" onClick={() => setShowCensureModal(true)}>
                  Political Censure
                </button>
              </li>
              <li>
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
      </div>
    </nav>
  );
}
