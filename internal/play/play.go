// Package play runs code in the Go Playground.
package play

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const compileURL = "https://play.golang.org/compile"

// Result is the result of running some code in the playground.
type Result struct {
	Errors string
	Events []Event
}

// Event is an output event from code run in the playground.
type Event struct {
	// Message is the content of the output.
	Message string

	// Kind is the mechanism of the output, either stdout or stderr.
	Kind string
}

// Run runs some code in the Go Playground. If the code is submitted
// successfully and a result is returned, err will be nil, even if the
// playground itself failed to run due to an error in the code
// submitted.
func Run(source string) (result Result, err error) {
	data := make(url.Values)
	data.Set("body", source)

	rsp, err := http.PostForm(compileURL, data)
	if err != nil {
		return result, fmt.Errorf("post: %w", err)
	}
	defer rsp.Body.Close()

	buf, err := io.ReadAll(rsp.Body)
	if err != nil {
		return result, fmt.Errorf("read body: %w", err)
	}

	err = json.Unmarshal(buf, &result)
	if err != nil {
		return result, fmt.Errorf("decode result: %w", err)
	}
	return result, nil
}
