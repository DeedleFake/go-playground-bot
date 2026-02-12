//go:build goexperiment.jsonv2

package play

import (
	"encoding/json/v2"
	"fmt"
	"io"

	"github.com/DeedleFake/go-playground-bot/internal/pool"
)

func unmarshal(result *Result, src io.Reader) error {
	buf := pool.GetBuffer()
	defer pool.PutBuffer(buf)
	r := io.TeeReader(src, buf)

	err := json.UnmarshalRead(r, &result)
	if err != nil {
		return fmt.Errorf("decode result: %w\n%q", err, buf)
	}
	return nil
}
