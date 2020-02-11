package main

import server "github.com/jamesjarvis/risotto-play/pkg/server"

func main() {
	r := server.SetupRouter()
	// Listen and serve router on port 80
	r.Run(":4000")
}
