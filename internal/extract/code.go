// Package extract provides functionality for extracting information
// from Discord messages.
package extract

import (
	"iter"
	"strings"
	"unsafe"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// CodeBlock provides information about a block of code extracted from
// a message.
type CodeBlock struct {
	// Language is the language of the code block. If no language was
	// specified, this will be empty.
	Language string

	// Source is the actual content of the code block.
	Source string
}

// CodeBlocks extracts code blocks from the provided message contents.
func CodeBlocks(content string) iter.Seq[CodeBlock] {
	return func(yield func(CodeBlock) bool) {
		contentbytes := unsafe.Slice(unsafe.StringData(content), len(content))
		root := goldmark.DefaultParser().Parse(text.NewReader(contentbytes))
		ast.Walk(root, walker(func(n *ast.FencedCodeBlock) bool {
			var cb CodeBlock
			if n.Info != nil {
				cb.Language = content[n.Info.Segment.Start:n.Info.Segment.Stop]
			}

			var source strings.Builder
			for _, segment := range allSegments(n.Lines()) {
				source.WriteString(content[segment.Start:segment.Stop])
			}
			cb.Source = source.String()

			return yield(cb)
		}))
	}
}

// walker returns an [ast.Walker] that calls the provided yield
// function for each node in the tree that is an instance of T. It
// skips the children of nodes that match.
func walker[T ast.Node](yield func(T) bool) ast.Walker {
	return func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		cb, ok := n.(T)
		if !ok {
			return ast.WalkContinue, nil
		}

		if !yield(cb) {
			return ast.WalkStop, nil
		}
		return ast.WalkSkipChildren, nil
	}
}

// allSegments returns an iterator over the text segments represented
// by segments.
func allSegments(segments *text.Segments) iter.Seq2[int, text.Segment] {
	return func(yield func(int, text.Segment) bool) {
		for i := range segments.Len() {
			if !yield(i, segments.At(i)) {
				return
			}
		}
	}
}
