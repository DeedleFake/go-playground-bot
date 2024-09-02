// Package play runs code in the Go Playground.
package play

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"unsafe"

	"github.com/DeedleFake/go-playground-bot/internal/pool"
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
func Run(ctx context.Context, source string) (result Result, err error) {
	data := make(url.Values)
	data.Set("body", source)
	data.Set("withVet", "true")

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		compileURL,
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return result, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("post: %w", err)
	}
	defer rsp.Body.Close()

	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)

	_, err = io.Copy(buf, rsp.Body)
	if err != nil {
		return result, fmt.Errorf("read body: %w", err)
	}

	err = json.Unmarshal(buf.Bytes(), &result)
	if err != nil {
		return result, fmt.Errorf("decode result: %w\n%q", err, buf)
	}
	if noPackageError.MatchString(result.Errors) {
		source, err := MainWrap(source)
		if err != nil {
			return result, fmt.Errorf("wrap with main: %w", err)
		}
		return Run(ctx, source) // TODO: Adjust error line numbers?
	}

	return result, nil
}

// MainWrap wraps source code in a main package and main function.
func MainWrap(source string) (string, error) {
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)

	buf.WriteString("package main\n\nfunc main() {\n")
	buf.WriteString(source)
	buf.WriteString("\n}")

	formatted, err := imports.Process("prog.go", buf.Bytes(), nil)
	return unsafe.String(unsafe.SliceData(formatted), len(formatted)), err
}
