// src/components/GameModals.jsx
import React from "react";
import AgendaModals from "./AgendaModals";
import ImperialRiderModal from "./ImperialRiderModal";
import RelicModal from "./RelicModal";
import { handleShardSubmit } from "../../utils/relicHandlers";
import { handleAssignObsidian } from "../../utils/obsidianhandler";
import { handleMutinySubmit, handleSeedSubmit, handleIncentiveSubmit, handlePoliticalCensureSubmit, handleClassifiedSubmit } from "../../utils/agendaHandlers";

export default function GameModals({
  modals,
  toggleModal,
  gameId,
  game,
  refreshGameState,
  playersSorted,
  playersUnsorted,
  scoreImperialRider,
  scoreCrown,
  setGame,
  agendaModal,
  setAgendaModal,
  secretObjectives,
  groupedScoredSecrets,
  mutinyVotes,
  setMutinyVotes,
  mutinyResult,
  setMutinyResult,
  mutinyAbstained,
  setMutinyAbstained
}) {
  return (
    <>
      <AgendaModals
        agendaModal={agendaModal}
        setAgendaModal={setAgendaModal}
        showAgendaModal={modals.agenda}
        setShowAgendaModal={(val) => toggleModal("agenda", val)}
        showSeedModal={modals.seed}
        setShowSeedModal={(val) => toggleModal("seed", val)}
        showCensureModal={modals.censure}
        setShowCensureModal={(val) => toggleModal("censure", val)}
        mutinyVotes={mutinyVotes}
        setMutinyVotes={setMutinyVotes}
        mutinyResult={mutinyResult}
        setMutinyResult={setMutinyResult}
        mutinyAbstained={mutinyAbstained}
        setMutinyAbstained={setMutinyAbstained}
        playersUnsorted={playersUnsorted}
        gameId={gameId}
        game={game}
        refreshGameState={refreshGameState}
        setGame={setGame}
        secretObjectives={secretObjectives}
        groupedScoredSecrets={groupedScoredSecrets}
        handleMutinySubmit={() =>
          handleMutinySubmit({
            gameId,
            game,
            mutinyResult,
            mutinyVotes,
            mutinyAbstained,
            refreshGameState,
            setGame,
            setShowAgendaModal: () => toggleModal("agenda", false),
          })
        }
        handleSeedSubmit={(result) =>
          handleSeedSubmit({
            gameId,
            result,
            game,
            refreshGameState,
          })
        }
        handlePoliticalCensureSubmit={(data) =>
          handlePoliticalCensureSubmit({
            ...data,
            gameId,
            game,
            refreshGameState,
          })
        }
        handleClassifiedSubmit={(playerId, objectiveId) =>
          handleClassifiedSubmit({
            gameId,
            playerId,
            objectiveId,
            refreshGameState,
            setAgendaModal,
          })
        }
        handleIncentiveSubmit={(result) =>
          handleIncentiveSubmit({
            gameId,
            result,
            refreshGameState,
            setAgendaModal,
          })
        }
      />

      <ImperialRiderModal
        show={modals.imperial}
        onClose={() => toggleModal("imperial", false)}
        players={playersSorted}
        onSubmit={scoreImperialRider}
      />

      <RelicModal
        show={modals.crown}
        onClose={() => toggleModal("crown", false)}
        title="The Crown of Emphidia"
        players={playersSorted}
        onSubmit={scoreCrown}
        description="Choose a player to gain 1 point from The Crown of Emphidia"
      />

      <RelicModal
        show={modals.shard}
        onClose={() => toggleModal("shard", false)}
        title="Shard of the Throne"
        players={playersSorted}
        onSubmit={(playerId) =>
          handleShardSubmit(playerId, gameId, refreshGameState, () =>
            toggleModal("shard", false)
          )
        }
        description="Choose a player to gain control of Shard of the Throne (gain 1 point)"
      />

      <RelicModal
        show={modals.obsidian}
        onClose={() => toggleModal("obsidian", false)}
        title="The Obsidian"
        players={playersSorted}
        onSubmit={(playerId) =>
          handleAssignObsidian(playerId, gameId, refreshGameState, () =>
            toggleModal("obsidian", false)
          )
        }
        description="Choose a player to gain The Obsidian (gain 1 extra secret objective slot)"
      />
    </>
  );
}
