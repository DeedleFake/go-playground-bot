//go:build !goexperiment.jsonv2

package play

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/DeedleFake/go-playground-bot/internal/pool"
)

func unmarshal(result *Result, src io.Reader) error {
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)

	_, err := io.Copy(buf, src)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	err = json.Unmarshal(buf.Bytes(), &result)
	if err != nil {
		return fmt.Errorf("decode result: %w\n%q", err, buf)
	}

	return nil
}
