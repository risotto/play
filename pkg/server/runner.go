package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

// Response is the response from when a file is ran
type Response struct {
	Errors string    `json:"errors"`
	Output string    `json:"output"`
	Status int       `json:"status"`
	Time   time.Time `json:"time"`
}

// ServerError is returned in the event of the server crapping itself
type ServerError struct {
	Error string `json:"error"`
}

// Server is the server config
type Server struct {
	Timeout time.Duration
}

// RunCode takes Risotto code as a string, then does the following:
// Checks to see if it is in the cache already, if so, return the output from the cache
// Saves the file as a temporary file (since there seems to be no other way of running it)
// Runs the file using the risotto binary which needs to be located in the PATH of whatever this runs on
// Then removes the local file
// Then stores the output in the cache (in another goprocess to speed things up a little)
// Then returns the output
func (s *Server) RunCode(b []byte) (*Response, error) {

	tmpfile, err := ioutil.TempFile("", "somecode")
	if err != nil {
		return nil, err
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(b); err != nil {
		return nil, err
	}
	if err := tmpfile.Close(); err != nil {
		return nil, err
	}

	response := s.RunRisotto(tmpfile.Name())

	return response, nil
}

// RunRisotto runs risotto
func (s *Server) RunRisotto(filename string) *Response {
	var stdOut bytes.Buffer
	var stdErr bytes.Buffer

	cmd := exec.Command("rst", filename)
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	if err := cmd.Start(); err != nil {
		return &Response{
			Errors: err.Error(),
			Output: "",
			Status: -1,
		}
	}

	// Wait for the process to finish or kill it after a timeout (whichever happens first):
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-time.After(s.Timeout):
		if err := cmd.Process.Kill(); err != nil {
			return &Response{
				Errors: err.Error(),
				Output: "",
				Status: -1,
			}
		}

		return &Response{
			Errors: fmt.Sprintf("Timeout limit of %v exceeded.\nCalm the hecc down pls I don't pay for a good cluster...", s.Timeout),
			Output: "",
			Status: -1,
		}
	case err := <-done:
		status := cmd.ProcessState.ExitCode()

		if err != nil {
			return &Response{
				Errors: err.Error(),
				Output: "",
				Status: status,
			}
		}

		return &Response{
			Errors: stdErr.String(),
			Output: stdOut.String(),
			Status: status,
		}
	}
}

func getMillis(t *time.Time) int64 {
	return t.UnixNano() / 1000000
}
