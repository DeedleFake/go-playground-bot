package extract

import (
	"iter"
	"strings"
	"unsafe"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

type CodeBlock struct {
	Language string
	Source   string
}

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

func allSegments(segments *text.Segments) iter.Seq2[int, text.Segment] {
	return func(yield func(int, text.Segment) bool) {
		for i := range segments.Len() {
			if !yield(i, segments.At(i)) {
				return
			}
		}
	}
}
