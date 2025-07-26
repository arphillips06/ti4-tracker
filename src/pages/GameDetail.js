import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";

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
import '../GameDetail.css';
import ScoreGraph from '../components/graphs/ScoreGraph';
import VictoryBanner from "../components/layout/VictoryBanner";
import GameModals from "../components/modals/GameModals";
import useGameData from "../hooks/useGameData";
import SpeakerModal from "../components/modals/SpeakerModal";
import API_BASE_URL from "../config"

export default function GameDetail() {
  const { gameId } = useParams();
  const { modals, toggleModal } = useModalControls();

  const scoreImperialRider = (playerId) => {
    return handleScoreImperialRider(
      playerId,
      gameId,
      refreshGameState,
      () => toggleModal("imperial", false)
    );
  };

  const scoreCrown = (playerId) =>
    handleScoreCrown(playerId, gameId, refreshGameState);

  const [graphRefreshKey, setGraphRefreshKey] = useState(0);
  const triggerGraphUpdate = () => {
    setGraphRefreshKey((prev) => prev + 1);
  };

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

  const isAgendaUsed = (title) =>
    game?.AllScores?.some(
      (s) => s.Type?.toLowerCase() === "agenda" && s.AgendaTitle === title
    );
  const isRelicUsed = (title) =>
    game?.AllScores?.some(
      (s) => s.Type?.toLowerCase() === "relic" && s.RelicTitle === title
    );

  const mutinyUsed = isAgendaUsed("Mutiny");
  const incentiveUsed = isAgendaUsed("Incentive Program");
  const seedUsed = isAgendaUsed("Seed of an Empire");

  const [scoringMode, setScoringMode] = useState(false);
  const [expandedPlayers, setExpandedPlayers] = useState({});
  const [localScored, setLocalScored] = useState({});
  const [mutinyVotes, setMutinyVotes] = useState([]);
  const [mutinyResult, setMutinyResult] = useState("for");
  const [mutinyAbstained, setMutinyAbstained] = useState(false);
  const [custodiansScorerId, setCustodiansScorerId] = useState(null);
  const [assigningObjective, setAssigningObjective] = useState(null);
  const [agendaModal, setAgendaModal] = useState(null);
  const [showSpeakerModal, setShowSpeakerModal] = useState(false);

  const playersUnsorted = useMergedPlayerData(game, false);
  const playersSorted = useMergedPlayerData(game, true);
  const groupedScoredSecrets = useGroupedScoredSecrets(game?.AllScores, secretObjectives);
  const obsidianUsed = isRelicUsed("The Obsidian");
  const [showScoreGraph, setShowScoreGraph] = useState(false);
  const [allScores, setAllScores] = useState(game?.AllScores || []);
  const assignSpeaker = async (playerId, isInitial = false) => {

    const roundId = game?.current_round?.id;
    if (!roundId) {
      alert("Cannot assign speaker: round not loaded.");
      return;
    }

    await fetch(`${API_BASE_URL}/games/${game.id}/speaker`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        player_id: playerId,
        round_id: roundId,
        is_initial: isInitial,
      }),
    });

    await refreshGameState();
  };
  const handleOpenSpeakerModal = () => {
    console.log("🧪 current_round.id =", game?.current_round?.id);
    setShowSpeakerModal(true);
  };

  const {
    scoreObjective,
    unscoreObjective,
    assignObjective,
    advanceRound,
  } = useObjectiveActions(gameId, refreshGameState, setLocalScored);

  useEffect(() => {
    const scorer = game?.AllScores?.find((s) => s.Type === "mecatol");
    setCustodiansScorerId(scorer?.PlayerID || null);
  }, [game]);
  useEffect(() => {
    console.log("📣 showSpeakerModal state is now:", showSpeakerModal);
  }, [showSpeakerModal]);

  if (!game) return <div className="p-4 text-warning">Loading game data...</div>;
  const winner_id = game.winner_id || game.winner_id;
  const winnerPlayer = game.players?.find(p => p.PlayerID === winner_id || p.player_id === winner_id);

  return (
    <>
      <GameNavbar
        gameId={gameId}
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
        crownUsed={crownUsed}
        obsidianUsed={obsidianUsed}
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
        <div className="d-flex flex-row align-items-start gap-4 flex-wrap">
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
              gameId={game.id || game.ID}
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
              custodiansScored={!!game?.AllScores?.some((s) => s.Type === "mecatol")}
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
