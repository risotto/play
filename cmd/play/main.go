package main

import (
	"time"

	"github.com/jamesjarvis/risotto-play/pkg/server"
)

func main() {
	s := &server.Server{
		Timeout: 5 * time.Second,
	}

	r := s.SetupRouter()
	// Listen and serve router on port 80
	r.Run(":4000")
}
