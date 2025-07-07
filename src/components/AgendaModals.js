import React from "react";

import MutinyModal from "../components/MutinyModal";
import PoliticalCensureModal from "../components/PoliticalCensureModal";
import SeedOfEmpireModal from "../components/SeedOfEmpireModal";
import ClassifiedDocumentLeaksModal from "../components/ClassifiedDocumentLeaksModal";
import IncentiveProgramModal from "../components/IncentiveProgramModal";

export default function AgendaModals({
  game,
  gameId,
  playersUnsorted,
  secretObjectives,
  groupedScoredSecrets,
  mutinyResult,
  setMutinyResult,
  mutinyAbstained,
  setMutinyAbstained,
  mutinyVotes,
  setMutinyVotes,
  showAgendaModal,
  setShowAgendaModal,
  showSeedModal,
  setShowSeedModal,
  showCensureModal,
  setShowCensureModal,
  agendaModal,
  setAgendaModal,
  handleMutinySubmit,
  handleSeedSubmit,
  handlePoliticalCensureSubmit,
  handleClassifiedSubmit,
  handleIncentiveSubmit,
}) {
  return (
    <>
      <MutinyModal
        show={showAgendaModal}
        onClose={() => setShowAgendaModal(false)}
        mutinyResult={mutinyResult}
        setMutinyResult={setMutinyResult}
        mutinyAbstained={mutinyAbstained}
        setMutinyAbstained={setMutinyAbstained}
        mutinyVotes={mutinyVotes}
        setMutinyVotes={setMutinyVotes}
        players={playersUnsorted}
        onSubmit={handleMutinySubmit}
      />

      <SeedOfEmpireModal
        show={showSeedModal}
        onClose={() => setShowSeedModal(false)}
        onSubmit={handleSeedSubmit}
      />

      <PoliticalCensureModal
        show={showCensureModal}
        onClose={() => setShowCensureModal(false)}
        onSubmit={handlePoliticalCensureSubmit}
        players={playersUnsorted.map((p) => ({
          ...p,
          agendaScores:
            game.all_scores?.filter(
              (s) => s.PlayerID === p.player_id && s.Type?.toLowerCase() === "agenda"
            ) || [],
        }))}
      />

      <ClassifiedDocumentLeaksModal
        show={agendaModal === "Classified Document Leaks"}
        players={playersUnsorted}
        secretObjectives={secretObjectives}
        scoredSecrets={groupedScoredSecrets}
        onClose={() => setAgendaModal(null)}
        onSubmit={handleClassifiedSubmit}
      />

      <IncentiveProgramModal
        show={agendaModal === "Incentive Program"}
        onClose={() => setAgendaModal(null)}
        onSubmit={handleIncentiveSubmit}
      />
    </>
  );
}