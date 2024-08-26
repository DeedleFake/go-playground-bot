package play

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const compileURL = "https://play.golang.org/compile"

type Result struct {
	Errors string
	Events []Event
}

type Event struct {
	Message string
	Kind    string
}

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
