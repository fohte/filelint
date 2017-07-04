package lint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndentRule_Lint(t *testing.T) {
	tests := []struct {
		rule IndentRule
		src  []byte
		want []byte
	}{
		// replace to softtab indents
		{
			rule: IndentRule{Style: IndentStyleSoft, Size: 2},
			src:  []byte("\t"),
			want: []byte("  "),
		},
		{
			rule: IndentRule{Style: IndentStyleSoft, Size: 4},
			src:  []byte("\t"),
			want: []byte("    "),
		},
		{
			rule: IndentRule{Style: IndentStyleSoft, Size: 2},
			src:  []byte("\t\t"),
			want: []byte("    "),
		},
		{
			rule: IndentRule{Style: IndentStyleSoft, Size: 2},
			src:  []byte("  "),
			want: []byte("  "),
		},
		{
			rule: IndentRule{Style: IndentStyleSoft, Size: 2},
			src:  []byte("\t\n\t"),
			want: []byte("  \n  "),
		},
		{
			rule: IndentRule{Style: IndentStyleSoft, Size: 2},
			src:  []byte("\t\r\n\t"),
			want: []byte("  \r\n  "),
		},
		{
			rule: IndentRule{Style: IndentStyleSoft, Size: 2},
			src:  []byte("\n\t\n  "),
			want: []byte("\n  \n  "),
		},
		{
			rule: IndentRule{Style: IndentStyleSoft, Size: 4},
			src:  []byte(".\n  "),
			want: []byte(".\n    "),
		},

		// replace to hardtab indents
		{
			rule: IndentRule{Style: IndentStyleHard},
			src:  []byte(".\n  "),
			want: []byte(".\n\t"),
		},
		{
			rule: IndentRule{Style: IndentStyleHard, Size: 2},
			src:  []byte(".\n  "),
			want: []byte(".\n\t"),
		},
		{
			rule: IndentRule{Style: IndentStyleHard},
			src:  []byte(".\n  \n    "),
			want: []byte(".\n\t\n\t\t"),
		},
		{
			rule: IndentRule{Style: IndentStyleHard},
			src:  []byte(".\r\n  \r\n    "),
			want: []byte(".\r\n\t\r\n\t\t"),
		},
		{
			rule: IndentRule{Style: IndentStyleHard},
			src:  []byte(".\n\t\n  "),
			want: []byte(".\n\t\n\t"),
		},
	}

	for _, tt := range tests {
		got, _ := tt.rule.Lint(tt.src)
		assert.Equal(t, tt.want, got.Fixed)
	}
}

func TestDetectSoftIndentWidth(t *testing.T) {
	lines := func(ls ...string) [][]byte {
		bs := [][]byte{}
		for _, l := range ls {
			bs = append(bs, []byte(l))
		}
		return bs
	}

	tests := []struct {
		src  [][]byte
		want int
	}{
		{lines(
			".",
			".",
		), 0},
		{lines(
			".",
			"  .",
		), 2},
		{lines(
			"  .",
			".",
		), 2},
		{lines(
			".",
			"  .",
			"    .",
		), 2},
		{lines(
			".",
			"  .",
			"  .",
		), 2},
		{lines(
			".",
			"  .",
			"",
		), 2},
		{lines(
			"{",
			"  {",
			"    .",
			"",
			"    .",
			"",
			"    .",
			"",
			"    .",
			"",
			"    .",
			"  }",
			"}",
		), 2},
		{lines(
			".",
			"  .",
			"      .",
		), 0}, // cannot predict
	}

	for _, tt := range tests {
		got := detectSoftIndentWidth(tt.src)
		assert.Equal(t, tt.want, got)
	}
}
