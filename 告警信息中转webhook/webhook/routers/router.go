package routers

import (
	"github.com/gin-gonic/gin"
	"webhook/server"
)

func Router() {
	r := gin.Default()
	r.POST("/api/v2/alerts", server.Alerts)
	r.Run(":19093")
}
