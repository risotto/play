package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func SetupServer() *gin.Engine {
	s := &Server{
		Timeout: 5 * time.Second,
	}

	return s.SetupRouter()
}

func TestPingRoute(t *testing.T) {

	router := SetupServer()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}

func TestHelloWorld(t *testing.T) {

	router := SetupServer()

	helloworld := `println("Hello, world!")`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/compile", bytes.NewBufferString(helloworld))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, "Hello, world!\n", response.Output)
	assert.Equal(t, 0, response.Status)
	assert.Empty(t, response.Errors)
}

func TestError(t *testing.T) {
	router := SetupServer()

	wtf := "wtf"

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/compile", bytes.NewBufferString(wtf))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Empty(t, response.Output)
	assert.NotEqual(t, 0, response.Status)
	assert.NotEmpty(t, response.Errors)
}

func TestFibonacci(t *testing.T) {
	router := SetupServer()

	input := `
	func fib(n int) int {
		if n == 0 {
		return 0
		}

		if n == 1 {
		return 1
		}

		return fib(n-1) + fib(n-2)
	}

	println(fib(30))`

	w := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/compile", bytes.NewBufferString(input))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 0, response.Status)
	assert.Empty(t, response.Errors)
}
