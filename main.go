package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/risotto/play/pkg/server"
)

func main() {
	s := &server.Server{
		Timeout:      3 * time.Second,
		MaxPerSecond: 10,
		SizeLimit:    10000,
	}

	gin.SetMode("release")
	gr := gin.New()
	gr.Use(gin.Logger(), gin.Recovery())

	r := s.SetupRouter(gr)
	// Listen and serve router on port 4000
	r.Run(":4000")
}
