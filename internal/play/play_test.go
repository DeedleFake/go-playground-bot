package play_test

import (
	"testing"

	"github.com/DeedleFake/go-playground-bot/internal/play"
)

func TestMainWrap(t *testing.T) {
	tests := []struct {
		name          string
		before, after string
	}{
		{
			"JustFormat",
			`println("test")`,
			"package main\n\nfunc main() {\n\tprintln(\"test\")\n}\n",
		},
		{
			"OneImport",
			`fmt.Println("test")`,
			"package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"test\")\n}\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			source, err := play.MainWrap(test.before)
			if err != nil {
				t.Fatal(err)
			}
			if source != test.after {
				t.Fatal(source)
			}
		})
	}
}
