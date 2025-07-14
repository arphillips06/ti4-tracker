export const isAgendaUsed = (game, title) =>
  game?.all_scores?.some(
    (s) => s.Type?.toLowerCase() === "agenda" && s.AgendaTitle === title
  );
