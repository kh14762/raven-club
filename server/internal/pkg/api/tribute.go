package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"net/http"
	"ravens-club/internal/pkg/tribute"
	"strconv"
)

var TributeControllerModule = fx.Module("tribute-controller", fx.Invoke(initTributeController))

func initTributeController(engine *gin.Engine, ts tribute.TributeService) {
	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	engine.GET("/api/tribute", func(c *gin.Context) {
		tributeList := ts.List()
		c.JSON(http.StatusOK, gin.H{"tributes": tributeList})
	})

	engine.GET("/api/tribute/:id", func(c *gin.Context) {
		param := c.Param("id")
		id, _ := strconv.Atoi(param)
		t, ok := ts.Get(uint16(id))
		if ok {
			c.JSON(http.StatusOK, gin.H{"tribute": t})
		} else {
			c.JSON(http.StatusNotFound, NewError(fmt.Sprintf("unknown device id: %d", id)))
		}
	})
}
