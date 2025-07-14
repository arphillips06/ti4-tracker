package controllers

import (
	"net/http"

	"github.com/arphillips06/TI4-stats/helpers"
	"github.com/arphillips06/TI4-stats/models"
	"github.com/arphillips06/TI4-stats/services"
	"github.com/gin-gonic/gin"
)

// ResolveMutinyAgenda applies the "Mutiny" agenda resolution.
// It awards 1 point to players who voted "for", then ends the process.
func ResolveMutinyAgenda(c *gin.Context) {
	input, ok := helpers.BindJSON[models.AgendaResolution](c)
	if !ok {
		return
	}
	err := services.ApplyMutinyAgenda(*input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// HandlePoliticalCensure bans the elected player from voting on the next agenda.
func HandlePoliticalCensure(c *gin.Context) {
	input, ok := helpers.BindJSON[models.PoliticalCensureRequest](c)
	if !ok {
		return
	}
	err := services.ApplyPoliticalCensure(*input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// HandleSeedOfEmpire grants a victory point to the elected player
// and creates a public 1-point objective on their behalf.
func HandleSeedOfEmpire(c *gin.Context) {
	input, ok := helpers.BindJSON[models.SeedOfEmpireResolution](c)
	if !ok {
		return
	}
	err := services.ApplySeedOfEmpire(*input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// HandleClassifiedDocumentLeaks turns a selected scored secret objective into a public one.
// The player retains the point but loses a secret objective slot.
func HandleClassifiedDocumentLeaks(c *gin.Context) {
	input, ok := helpers.BindJSON[models.ClassifiedDocumentLeaksRequest](c)
	if !ok {
		return
	}
	err := services.ApplyClassifiedDocumentLeaks(*input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// HandleIncentiveProgram grants 1 point to all players who voted
// in favor or against, depending on the outcome of the vote.
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
