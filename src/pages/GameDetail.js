import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";

import PlayerSidebar from "../components/PlayerSidebar";
import GameControls from "../components/GameControls";
import ObjectivesGrid from "../components/ObjectivesGrid";
import GameNavbar from "../components/GameNavbar";
import AgendaModals from "../components/AgendaModals";
import useGroupedScoredSecrets from "../hooks/useGroupedScoredSecrets";
import useMergedPlayerData from "../hooks/useMergedPlayerData";
import useObjectiveActions from "../hooks/useObjectiveActions";
import ImperialRiderModal from "../components/ImperialRiderModal";
import { handleScoreImperialRider } from "../utils/imperialRiderHandler";
import RelicModal from "../components/RelicModal";
import { handleScoreCrown } from "../utils/relicHandler";
import { handleShardSubmit } from "../utils/relicHandlers";


import {
  handleMutinySubmit,
  handleSeedSubmit,
  handleIncentiveSubmit,
  handlePoliticalCensureSubmit,
  handleClassifiedSubmit,
} from "../utils/agendaHandlers";
import useGameData from "../hooks/useGameData";


export default function GameDetail() {
  const { gameId } = useParams();
  const scoreImperialRider = (playerId) => {
    return handleScoreImperialRider(
      playerId,
      gameId,
      refreshGameState,
      setShowImperialModal
    );
  };
  const [showCrownModal, setShowCrownModal] = useState(false);

  const scoreCrown = (playerId) =>
    handleScoreCrown(playerId, gameId, refreshGameState);


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
    obsidianHolder,
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
  const [showAgendaModal, setShowAgendaModal] = useState(false);
  const [mutinyVotes, setMutinyVotes] = useState([]);
  const [mutinyResult, setMutinyResult] = useState("for");
  const [mutinyAbstained, setMutinyAbstained] = useState(false);
  const [custodiansScorerId, setCustodiansScorerId] = useState(null);
  const [showCensureModal, setShowCensureModal] = useState(false);
  const [showSeedModal, setShowSeedModal] = useState(false);
  const [assigningObjective, setAssigningObjective] = useState(null);
  const [agendaModal, setAgendaModal] = useState(null);
  const [showImperialModal, setShowImperialModal] = useState(false);
  const playersUnsorted = useMergedPlayerData(game, false);
  const playersSorted = useMergedPlayerData(game, true);
  const groupedScoredSecrets = useGroupedScoredSecrets(game?.all_scores, secretObjectives);
  const [showShardModal, setShowShardModal] = useState(false);
  const [showObsidianModal, setShowObsidianModal] = useState(false);
  const obsidianUsed = isRelicUsed("The Obsidian");
const [obsidianHolderId, setObsidianHolderId] = useState(null);

  const {
    scoreObjective,
    unscoreObjective,
    assignObjective,
    advanceRound,
  } = useObjectiveActions(gameId, refreshGameState, setLocalScored);

  useEffect(() => {
    const scorer = game?.all_scores?.find((s) => s.Type === "mecatol");
    setCustodiansScorerId(scorer?.PlayerID || null);
  }, [game]);

  if (!game) return <div className="p-4 text-warning">Loading game data...</div>;

  return (
    <>
      <GameNavbar
        mutinyUsed={mutinyUsed}
        incentiveUsed={incentiveUsed}
        seedUsed={seedUsed}
        setShowAgendaModal={setShowAgendaModal}
        setShowCensureModal={setShowCensureModal}
        setShowSeedModal={setShowSeedModal}
        setAgendaModal={setAgendaModal}
        setShowImperialModal={setShowImperialModal}
        setShowShardModal={setShowShardModal}
        setShowCrownModal={setShowCrownModal}
        crownUsed={crownUsed}
        obsidianUsed={obsidianUsed}
        setShowObsidianModal={setShowObsidianModal}
      />

      <div className="p-6 max-w-7xl mx-auto">
        <GameControls
          game={game}
          scoringMode={scoringMode}
          setScoringMode={setScoringMode}
          onAdvanceRound={advanceRound}
        />
        <div className="d-flex flex-row align-items-start gap-4 flex-wrap">
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

          <AgendaModals
            agendaModal={agendaModal}
            setAgendaModal={setAgendaModal}
            showAgendaModal={showAgendaModal}
            setShowAgendaModal={setShowAgendaModal}
            showSeedModal={showSeedModal}
            setShowSeedModal={setShowSeedModal}
            showCensureModal={showCensureModal}
            setShowCensureModal={setShowCensureModal}
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
                setShowAgendaModal,
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
            show={showImperialModal}
            onClose={() => setShowImperialModal(false)}
            players={playersSorted}
            onSubmit={scoreImperialRider}
          />
          <RelicModal
            show={showCrownModal}
            onClose={() => setShowCrownModal(false)}
            title="The Crown of Emphidia"
            players={playersSorted}
            onSubmit={scoreCrown}
            description="Choose a player to gain 1 point from The Crown of Emphidia"
          />
          <RelicModal
            show={showShardModal}
            onClose={() => setShowShardModal(false)}
            title="Shard of the Throne"
            players={playersSorted}
            onSubmit={(playerId) =>
              handleShardSubmit(playerId, gameId, refreshGameState, () => setShowShardModal(false))
            }
            description="Choose a player to gain control of Shard of the Throne (gain 1 point)"
          />
          <RelicModal
            show={showObsidianModal}
            onClose={() => setShowObsidianModal(false)}
            title="The Obsidian"
            players={playersSorted}
            onSubmit={(playerId) => {
              setObsidianHolderId(parseInt(playerId));
              setShowObsidianModal(false);
              console.log("🧿 Obsidian given to player ID:", playerId);
            }}

            description="Choose a player to gain The Obsidian (gain 1 extra secret objective slot)"
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
              custodiansScored={!!game?.all_scores?.some((s) => s.Type === "mecatol")}
              obsidianHolderId={obsidianHolderId}
            />
          </div>
        </div>
      </div>
    </>
  );
}