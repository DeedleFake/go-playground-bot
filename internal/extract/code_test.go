package extract_test

import (
	"slices"
	"testing"

	"deedles.dev/xiter"
	"github.com/DeedleFake/go-playground-bot/internal/extract"
)

const source = "before\n```go\nstuff\n```\nafter"

func TestCodeBlocks(t *testing.T) {
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

func BenchmarkCodeBlocks(b *testing.B) {
	for range b.N {
		xiter.Drain(extract.CodeBlocks(source))
	}
}
