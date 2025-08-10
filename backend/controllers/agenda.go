package controllers

import (
	"net/http"

	"github.com/arphillips06/TI4-stats/helpers"
	"github.com/arphillips06/TI4-stats/models"
	"github.com/arphillips06/TI4-stats/services"
	"github.com/gin-gonic/gin"
)

// ResolveMutinyAgenda godoc
// @Summary      Apply "Mutiny" agenda
// @Description  Awards 1 point to players who voted "for" and finalizes the result.
// @Tags         agendas
// @Accept       json
// @Produce      json
// @Param        body  body      models.AgendaResolution  true  "Game and resolution context (if applicable)"
// @Success      204
// @Failure      400  {object}  map[string]string  "error"
// @Failure      500  {object}  map[string]string  "error"
// @Router       /agendas/mutiny [post]
func ResolveMutinyAgenda(c *gin.Context) {
	helpers.HandleRequest(c, services.ApplyMutinyAgenda)
}

// HandlePoliticalCensure godoc
// @Summary      Apply "Political Censure" agenda
// @Description  Bans the elected player from voting on the next agenda.
// @Tags         agendas
// @Accept       json
// @Produce      json
// @Param        body  body      models.PoliticalCensureRequest  true  "Game ID and elected player"
// @Success      204
// @Failure      400  {object}  map[string]string  "error"
// @Failure      500  {object}  map[string]string  "error"
// @Router       /agendas/political-censure [post]
func HandlePoliticalCensure(c *gin.Context) {
	helpers.HandleRequest(c, services.ApplyPoliticalCensure)

}

// HandleSeedOfEmpire godoc
// @Summary      Apply "Seed of an Empire" agenda
// @Description  Grants a VP to the elected player and creates a public 1-point objective for them.
// @Tags         agendas
// @Accept       json
// @Produce      json
// @Param        body  body      models.SeedOfEmpireResolution  true  "Game ID and elected player"
// @Success      204
// @Failure      400  {object}  map[string]string  "error"
// @Failure      500  {object}  map[string]string  "error"
// @Router       /agendas/seed-of-empire [post]
func HandleSeedOfEmpire(c *gin.Context) {
	helpers.HandleRequest(c, services.ApplySeedOfEmpire)
}

// HandleClassifiedDocumentLeaks godoc
// @Summary      Apply "Classified Document Leaks"
// @Description  Selects a scored secret objective to become public; the scorer keeps the point but loses a secret slot.
// @Tags         agendas
// @Accept       json
// @Produce      json
// @Param        body  body      models.ClassifiedDocumentLeaksRequest  true  "Game ID, player, and target secret objective"
// @Success      204
// @Failure      400  {object}  map[string]string  "error"
// @Failure      500  {object}  map[string]string  "error"
// @Router       /agendas/classified-document-leaks [post]
func HandleClassifiedDocumentLeaks(c *gin.Context) {
	helpers.HandleRequest(c, services.ApplyClassifiedDocumentLeaks)
}

// HandleIncentiveProgram godoc
// @Summary      Apply "Incentive Program"
// @Description  Grants 1 point to all players who voted with the outcome (for or against).
// @Tags         agendas
// @Accept       json
// @Produce      json
// @Param        body  body      models.IncentiveProgramRequest  true  "Game ID and outcome"
// @Success      200  {object}  map[string]string  "message"
// @Failure      400  {object}  map[string]string  "error"
// @Failure      500  {object}  map[string]string  "error"
// @Router       /agendas/incentive-program [post]
func HandleIncentiveProgram(c *gin.Context) {
	req, ok := helpers.BindJSON[models.IncentiveProgramRequest](c)
	if !ok {
		return
	}

	err := services.ApplyIncentiveProgramEffect(req.GameID, req.Outcome)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Incentive Program applied"})
}
