// Package play runs code in the Go Playground.
package play

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"unsafe"

	"golang.org/x/tools/imports"
)

const compileURL = "https://play.golang.org/compile"

var noPackageError = regexp.MustCompile(`^prog\.go:1:1: expected 'package', found [a-zA-Z0-9_]+\n$`)

// Result is the result of running some code in the playground.
type Result struct {
	Errors      string
	Events      []Event
	IsTest      bool
	TestsFailed int
	VetErrors   string
	VetOK       bool
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
	data.Set("withVet", "true")

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
	if noPackageError.MatchString(result.Errors) {
		source, err := MainWrap(source)
		if err != nil {
			return result, fmt.Errorf("wrap with main: %w", err)
		}
		return Run(source) // TODO: Adjust error line numbers?
	}

	return result, nil
}

// MainWrap wraps source code in a main package and main function.
func MainWrap(source string) (string, error) {
	var buf bytes.Buffer
	buf.WriteString("package main\n\nfunc main() {\n")
	buf.WriteString(source)
	buf.WriteString("\n}")

	formatted, err := imports.Process("prog.go", buf.Bytes(), nil)
	return unsafe.String(unsafe.SliceData(formatted), len(formatted)), err
}
