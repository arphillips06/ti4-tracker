package controllers

import (
	handle "github.com/arphillips06/TI4-stats/errors"
	"github.com/gin-gonic/gin"
)

type RespondingHandler func(*gin.Context) (int, any, error)

func Wrap(handler RespondingHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		status, payload, err := handler(c)
		if err != nil {
			handle.Handle(c, err)
			return
		}
		if payload == nil {
			c.Status(status)
			return
		}
		c.JSON(status, payload)
	}
}
