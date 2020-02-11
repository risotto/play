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
	Errors string `json:"errors"`
	Output string `json:"output"`
	Status int    `json:"status"`
}

// ServerError is returned in the event of the server crapping itself
type ServerError struct {
	Error string `json:"error"`
}

// RunCode takes Risotto code as a string, then does the following:
// Checks to see if it is in the cache already, if so, return the output from the cache
// Saves the file as a temporary file (since there seems to be no other way of running it)
// Runs the file using the risotto binary which needs to be located in the PATH of whatever this runs on
// Then removes the local file
// Then stores the output in the cache (in another goprocess to speed things up a little)
// Then returns the output
func RunCode(b []byte) (*Response, error) {

	// TODO: Check if it is in the cache

	// Create the filename
	now := time.Now()
	currentMilliseconds := getMillis(&now)
	filename := fmt.Sprintf("/tmp/%v.rst", currentMilliseconds)

	// Save to a file
	if err := ioutil.WriteFile(filename, b, 0644); err != nil {
		return nil, err
	}

	// Run the risotto file
	response, err := runRisotto(filename)
	if err != nil {
		return nil, err
	}

	// Save to cache

	// Delete the file
	err = os.Remove(filename)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func runRisotto(filename string) (*Response, error) {
	cmd := exec.Command("/Users/jarjames/executables/rst", filename)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	bytesErr, err := ioutil.ReadAll(stderr)
	if err != nil {
		return nil, err
	}

	bytesOut, err := ioutil.ReadAll(stdout)
	if err != nil {
		return nil, err
	}

	cmd.Wait() // We ignore the error for now

	status := cmd.ProcessState.ExitCode()
	if status == 0 {
		return &Response{
			Errors: bytes.NewBuffer(bytesErr).String(),
			Output: bytes.NewBuffer(bytesOut).String(),
			Status: cmd.ProcessState.ExitCode(),
		}, nil
	}

	return &Response{
		Errors: bytes.NewBuffer(bytesErr).String(),
		Output: bytes.NewBuffer(bytesOut).String(),
		Status: cmd.ProcessState.ExitCode(),
	}, nil
}

func getMillis(t *time.Time) int64 {
	return t.UnixNano() / 1000000
}
