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
		Timeout:      1 * time.Second,
		MaxPerSecond: 10,
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

func TestCompileRateLimitOk(t *testing.T) {

	router := SetupServer()

	helloworld := `println("Hello, world!")`

	for i := 0; i < 10; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/compile", bytes.NewBufferString(helloworld))
		req.Header.Set("X-Real-IP", "2601:7:1c82:4097:59a0:a80b:2841:b8c8")
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
	}
}

func TestCompileRateLimitFail(t *testing.T) {

	router := SetupServer()

	helloworld := `for int:=0;i<100;i+=1{println("Hello, world!")}`

	countratelimited := 0

	for i := 0; i < 100; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/compile", bytes.NewBufferString(helloworld))
		req.Header.Set("X-Real-IP", "2601:7:1c82:4097:59a0:a80b:2841:b8c8")
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			countratelimited++
		}
	}
	assert.Equal(t, 90, countratelimited)
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

func TestTimeout(t *testing.T) {
	router := SetupServer()

	wtf := `
	while true {
		println("bitch")
	}
	`

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/compile", bytes.NewBufferString(wtf))
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response Response
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Empty(t, response.Output)
	assert.Equal(t, -1, response.Status)
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
	assert.Equal(t, "832040\n", response.Output)
	assert.Empty(t, response.Errors)
}
