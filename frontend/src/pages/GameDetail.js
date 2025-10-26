import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import {
  isAgendaUsed,
  isRelicUsed,
  custodiansScorerId as getCustodiansScorerId,
  winnerPlayerId
} from "../utils/selectors";

import PlayerSidebar from "../components/players/PlayerSidebar";
import GameControls from "../components/layout/GameControls";
import ObjectivesGrid from "../components/objectives/ObjectivesGrid";
import GameNavbar from "../components/layout/GameNavbar";
import useGroupedScoredSecrets from "../hooks/useGroupedScoredSecrets";
import useMergedPlayerData from "../hooks/useMergedPlayerData";
import useObjectiveActions from "../hooks/useObjectiveActions";
import useModalControls from "../hooks/useModalControls";
import { handleScoreImperialRider } from "../utils/imperialRiderHandler";
import { handleScoreCrown } from "../utils/relicHandlers";
import "../GameDetail.css";
import ScoreGraph from "../components/graphs/ScoreGraph";
import VictoryBanner from "../components/layout/VictoryBanner";
import GameModals from "../components/modals/GameModals";
import useGameData from "../hooks/useGameData";
import SpeakerModal from "../components/modals/SpeakerModal";

export default function GameDetail() {
  const { gameId } = useParams();
  const { modals, toggleModal } = useModalControls();

  const scoreImperialRider = (playerId) =>
    handleScoreImperialRider(playerId, gameId, refreshGameState, () =>
      toggleModal("imperial", false)
    );

  const scoreCrown = (playerId) =>
    handleScoreCrown(playerId, gameId, refreshGameState);
  
  const [graphRefreshKey, setGraphRefreshKey] = useState(0);
  const triggerGraphUpdate = () => setGraphRefreshKey((prev) => prev + 1);

  const {
    game,
    objectives,
    secretObjectives,
    secretCounts,
    setSecretCounts,
    objectiveScores,
    setGame,
    setObjectiveScores,
    refreshGameState,
    crownUsed,
    obsidianHolderId,
  } = useGameData(gameId);
console.log("DEBUG: game object", game);

  const mutinyUsed = isAgendaUsed(game, "Mutiny");
  const incentiveUsed = isAgendaUsed(game, "Incentive Program");
  const seedUsed = isAgendaUsed(game, "Seed of an Empire");

  const [scoringMode, setScoringMode] = useState(false);
  const [expandedPlayers, setExpandedPlayers] = useState({});
  const [localScored, setLocalScored] = useState({});
  const [mutinyVotes, setMutinyVotes] = useState([]);
  const [mutinyResult, setMutinyResult] = useState("for");
  const [mutinyAbstained, setMutinyAbstained] = useState(false);
  const [assigningObjective, setAssigningObjective] = useState(null);
  const [agendaModal, setAgendaModal] = useState(null);
  const [showSpeakerModal, setShowSpeakerModal] = useState(false);

  const playersUnsorted = useMergedPlayerData(game, false);
  const playersSorted = useMergedPlayerData(game, true);
  const groupedScoredSecrets = useGroupedScoredSecrets(
    game?.AllScores,
    secretObjectives
  );

  const obsidianUsed = isRelicUsed(game, "The Obsidian");
  const bookUsed = isRelicUsed(game, "Book Of Latvina");
  const shardUsed = isRelicUsed(game, "The Crown of Emphidia");

  const [showScoreGraph, setShowScoreGraph] = useState(false);
  const [allScores, setAllScores] = useState(game?.AllScores || []);

  const custodianOwnerId = getCustodiansScorerId(game);

  const handleOpenSpeakerModal = () => setShowSpeakerModal(true);
  const { scoreObjective, unscoreObjective, assignObjective, advanceRound } =
    useObjectiveActions(gameId, refreshGameState, setLocalScored);

  if (!game) return <div className="p-4 text-warning">Loading game data...</div>;
  const winner_id = winnerPlayerId(game);

  return (
    <>
      <GameNavbar
        gameId={gameId}
        gameNumber={game.game_number}
        showScoreGraph={showScoreGraph}
        setShowScoreGraph={setShowScoreGraph}
        mutinyUsed={mutinyUsed}
        incentiveUsed={incentiveUsed}
        seedUsed={seedUsed}
        setShowSpeakerModal={setShowSpeakerModal}
        setShowAgendaModal={(val) => toggleModal("agenda", val)}
        setShowCensureModal={(val) => toggleModal("censure", val)}
        setShowSeedModal={(val) => toggleModal("seed", val)}
        setAgendaModal={setAgendaModal}
        setShowImperialModal={(val) => toggleModal("imperial", val)}
        setShowShardModal={(val) => toggleModal("shard", val)}
        setShowCrownModal={(val) => toggleModal("crown", val)}
        setShowLatvinaModal={(val) => toggleModal("latvina", val)}
        crownUsed={crownUsed}
        obsidianUsed={obsidianUsed}
        bookUsed={bookUsed}
        setShowObsidianModal={(val) => toggleModal("obsidian", val)}
        refreshGameState={refreshGameState}
        onOpenSpeakerModal={handleOpenSpeakerModal}
      />

      <VictoryBanner
        winner={game?.winner}
        finished={game?.finished_at}
        victoryPathSummary={game?.victory_path}
      />

      <div className="p-6 max-w-7xl mx-auto">
        <GameControls
          game={game}
          scoringMode={scoringMode}
          setScoringMode={setScoringMode}
          onAdvanceRound={advanceRound}
        />

        <div className="d-flex flex-column flex-md-row align-items-start gap-4 flex-wrap">
          {showScoreGraph ? (
            <div style={{ flex: 1 }}>
              <ScoreGraph gameId={gameId} refreshSignal={graphRefreshKey} />
            </div>
          ) : (
            <ObjectivesGrid
              game={game}
              refreshGameState={refreshGameState}
              objectives={objectives}
              playersUnsorted={playersUnsorted}
              gameId={game.id}
              objectiveScores={objectiveScores}
              localScored={localScored}
              scoringMode={scoringMode}
              scoreObjective={scoreObjective}
              useObjectiveDecks={game?.use_objective_decks}
              assigningObjective={assigningObjective}
              setAssigningObjective={setAssigningObjective}
              assignObjective={assignObjective}
            />
          )}

          <GameModals
            modals={modals}
            toggleModal={toggleModal}
            gameId={gameId}
            game={game}
            refreshGameState={refreshGameState}
            playersSorted={playersSorted}
            playersUnsorted={playersUnsorted}
            scoreImperialRider={scoreImperialRider}
            scoreCrown={scoreCrown}
            setGame={setGame}
            agendaModal={agendaModal}
            setAgendaModal={setAgendaModal}
            secretObjectives={secretObjectives}
            groupedScoredSecrets={groupedScoredSecrets}
            mutinyVotes={mutinyVotes}
            setMutinyVotes={setMutinyVotes}
            mutinyResult={mutinyResult}
            setMutinyResult={setMutinyResult}
            mutinyAbstained={mutinyAbstained}
            setMutinyAbstained={setMutinyAbstained}
          />

          <div style={{ flex: "0 1 300px" }}>
            <PlayerSidebar
              playersSorted={playersSorted}
              expandedPlayers={expandedPlayers}
              setExpandedPlayers={setExpandedPlayers}
              game={game}
              scoreObjective={scoreObjective}
              unscoreObjective={unscoreObjective}
              secretCounts={secretCounts}
              setSecretCounts={setSecretCounts}
              secretObjectives={secretObjectives}
              objectiveScores={objectiveScores}
              gameId={gameId}
              setGame={setGame}
              setObjectiveScores={setObjectiveScores}
              refreshGameState={refreshGameState}
              custodiansScored={!!custodianOwnerId}
              obsidianHolderId={obsidianHolderId}
              triggerGraphUpdate={triggerGraphUpdate}
              allScores={game?.AllScores}
              setAllScores={setAllScores}
            />
            <SpeakerModal
              show={showSpeakerModal}
              onClose={() => setShowSpeakerModal(false)}
              players={game.players}
              gameId={game.id}
              refetchGame={refreshGameState}
              roundId={game?.current_round?.id ?? game?.current_round}
            />
          </div>
        </div>
      </div>
    </>
  );
}
