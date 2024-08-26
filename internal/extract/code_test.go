package extract_test

import (
	"slices"
	"testing"

	"github.com/DeedleFake/go-playground-bot/internal/extract"
)

func TestCodeBlocks(t *testing.T) {
	const source = "before\n```go\nstuff\n```\nafter"

	blocks := slices.Collect(extract.CodeBlocks(source))
	if len(blocks) != 1 {
		t.Fatalf("%#v", blocks)
	}
	if blocks[0].Language != "go" {
		t.Fatalf("%#v", blocks[0])
	}
	if blocks[0].Source != "stuff\n" {
		t.Fatalf("%#v", blocks[0])
	}
}
