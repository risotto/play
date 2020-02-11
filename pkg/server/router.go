package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter sets up the router
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Submit Risotto code to the server
	r.POST("/compile", func(c *gin.Context) {
		// Retrieve the text from the message, and then just straight up run the code I guess?
		rawData, err := c.GetRawData()
		if err != nil {
			log.Fatal(err)
		}

		response, err := RunCode(rawData)
		if err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, response)
	})

	return r
}
