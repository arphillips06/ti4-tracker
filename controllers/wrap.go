package controllers

import (
	handle "github.com/arphillips06/TI4-stats/errors"
	"github.com/gin-gonic/gin"
)

type H func(*gin.Context) (int, any, error)

func Wrap(h H) gin.HandlerFunc {
	return func(c *gin.Context) {
		status, payload, err := h(c)
		if err != nil {
			handle.Handle(c, err)
			return
		}
		c.JSON(status, payload)
	}
}
