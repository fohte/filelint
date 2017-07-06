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
		rep  []*Report
	}{
		// replace to softtab indents
		{
			rule: IndentRule{Style: IndentStyleSoft, Size: 2},
			src:  []byte("\t"),
			want: []byte("  "),
			rep: []*Report{
				{
					position: &Position{1, -1},
					message:  `Expected indent with 2 space(s) but used hardtabs (\t)`,
				},
			},
		},
		{
			rule: IndentRule{Style: IndentStyleSoft, Size: 4},
			src:  []byte("\t"),
			want: []byte("    "),
			rep: []*Report{
				{
					position: &Position{1, -1},
					message:  `Expected indent with 4 space(s) but used hardtabs (\t)`,
				},
			},
		},
		{
			rule: IndentRule{Style: IndentStyleSoft, Size: 2},
			src:  []byte("\t\t"),
			want: []byte("    "),
			rep: []*Report{
				{
					position: &Position{1, -1},
					message:  `Expected indent with 2 space(s) but used hardtabs (\t)`,
				},
			},
		},
		{
			rule: IndentRule{Style: IndentStyleSoft, Size: 2},
			src:  []byte("  "),
			want: []byte("  "),
			rep:  []*Report{},
		},
		{
			rule: IndentRule{Style: IndentStyleSoft, Size: 2},
			src:  []byte("\t\n\t"),
			want: []byte("  \n  "),
			rep: []*Report{
				{
					position: &Position{1, -1},
					message:  `Expected indent with 2 space(s) but used hardtabs (\t)`,
				},
				{
					position: &Position{2, -1},
					message:  `Expected indent with 2 space(s) but used hardtabs (\t)`,
				},
			},
		},
		{
			rule: IndentRule{Style: IndentStyleSoft, Size: 2},
			src:  []byte("\t\r\n\t"),
			want: []byte("  \r\n  "),
			rep: []*Report{
				{
					position: &Position{1, -1},
					message:  `Expected indent with 2 space(s) but used hardtabs (\t)`,
				},
				{
					position: &Position{2, -1},
					message:  `Expected indent with 2 space(s) but used hardtabs (\t)`,
				},
			},
		},
		{
			rule: IndentRule{Style: IndentStyleSoft, Size: 2},
			src:  []byte("\n\t\n  "),
			want: []byte("\n  \n  "),
			rep: []*Report{
				{
					position: &Position{2, -1},
					message:  `Expected indent with 2 space(s) but used hardtabs (\t)`,
				},
			},
		},
		{
			rule: IndentRule{Style: IndentStyleSoft, Size: 4},
			src:  []byte(".\n  "),
			want: []byte(".\n    "),
			rep: []*Report{
				{
					position: &Position{2, -1},
					message:  `Expected indent with 4 space(s) but used 2 space(s)`,
				},
			},
		},

		// replace to hardtab indents
		{
			rule: IndentRule{Style: IndentStyleHard},
			src:  []byte(".\n  "),
			want: []byte(".\n\t"),
			rep: []*Report{
				{
					position: &Position{2, -1},
					message:  `Expected indent with hardtabs (\t) but used 2 space(s)`,
				},
			},
		},
		{
			rule: IndentRule{Style: IndentStyleHard, Size: 2},
			src:  []byte(".\n  "),
			want: []byte(".\n\t"),
			rep: []*Report{
				{
					position: &Position{2, -1},
					message:  `Expected indent with hardtabs (\t) but used 2 space(s)`,
				},
			},
		},
		{
			rule: IndentRule{Style: IndentStyleHard},
			src:  []byte(".\n  \n    "),
			want: []byte(".\n\t\n\t\t"),
			rep: []*Report{
				{
					position: &Position{2, -1},
					message:  `Expected indent with hardtabs (\t) but used 2 space(s)`,
				},
				{
					position: &Position{3, -1},
					message:  `Expected indent with hardtabs (\t) but used 2 space(s)`,
				},
			},
		},
		{
			rule: IndentRule{Style: IndentStyleHard},
			src:  []byte(".\r\n  \r\n    "),
			want: []byte(".\r\n\t\r\n\t\t"),
			rep: []*Report{
				{
					position: &Position{2, -1},
					message:  `Expected indent with hardtabs (\t) but used 2 space(s)`,
				},
				{
					position: &Position{3, -1},
					message:  `Expected indent with hardtabs (\t) but used 2 space(s)`,
				},
			},
		},
		{
			rule: IndentRule{Style: IndentStyleHard},
			src:  []byte(".\n\t\n  "),
			want: []byte(".\n\t\n\t"),
			rep: []*Report{
				{
					position: &Position{3, -1},
					message:  `Expected indent with hardtabs (\t) but used 2 space(s)`,
				},
			},
		},

		// for the javadoc comment style
		{
			rule: IndentRule{Style: IndentStyleSoft, Size: 2},
			src:  []byte("/**\n *\n */"),
			want: []byte("/**\n *\n */"),
			rep:  []*Report{},
		},
		{
			rule: IndentRule{Style: IndentStyleSoft, Size: 2},
			src:  []byte(".\n\t/**\n\t *\n\t */"),
			want: []byte(".\n  /**\n   *\n   */"),
			rep: []*Report{
				{
					position: &Position{2, -1},
					message:  `Expected indent with 2 space(s) but used hardtabs (\t)`,
				},
				{
					position: &Position{3, -1},
					message:  `Expected indent with 2 space(s) but used hardtabs (\t)`,
				},
				{
					position: &Position{4, -1},
					message:  `Expected indent with 2 space(s) but used hardtabs (\t)`,
				},
			},
		},
		{
			rule: IndentRule{Style: IndentStyleSoft, Size: 4},
			src:  []byte(".\n  /**\n   *\n   */"),
			want: []byte(".\n    /**\n     *\n     */"),
			rep: []*Report{
				{
					position: &Position{2, -1},
					message:  `Expected indent with 4 space(s) but used 2 space(s)`,
				},
				{
					position: &Position{3, -1},
					message:  `Expected indent with 4 space(s) but used 2 space(s)`,
				},
				{
					position: &Position{4, -1},
					message:  `Expected indent with 4 space(s) but used 2 space(s)`,
				},
			},
		},
	}

	for _, tt := range tests {
		got, _ := tt.rule.Lint(tt.src)
		assert.Equal(t, tt.want, got.Fixed)
		assert.Equal(t, tt.rep, got.Reports)
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
		src        [][]byte
		javadocPos []*columnRange
		want       int
	}{
		{
			src: lines(
				".",
				".",
			),
			want: 0,
		},
		{
			src: lines(
				".",
				"  .",
			),
			want: 2,
		},
		{
			src: lines(
				"  .",
				".",
			),
			want: 2,
		},
		{
			src: lines(
				".",
				"  .",
				"    .",
			),
			want: 2,
		},
		{
			src: lines(
				".",
				"  .",
				"  .",
			),
			want: 2,
		},
		{
			src: lines(
				".",
				"  .",
				"",
			),
			want: 2,
		},
		{
			src: lines(
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
			),
			want: 2,
		},
		{
			src: lines(
				".",
				"  .",
				"      .",
			),
			want: 0,
		}, // cannot predict

		// for the javadoc comment style
		{
			src: lines(
				"/**",
				" *",
				" */",
			),
			javadocPos: []*columnRange{{1, 3}},
			want:       0,
		},
		{
			src: lines(
				"/**",
				" *",
				" *",
				" */",
			),
			javadocPos: []*columnRange{{1, 4}},
			want:       0,
		},
		{
			src: lines(
				".",
				"  /**",
				"   *",
				"   */",
			),
			javadocPos: []*columnRange{{2, 4}},
			want:       2,
		},
		{
			src: lines(
				".",
				"  /**",
				"   *",
				"   *",
				"   */",
			),
			javadocPos: []*columnRange{{2, 5}},
			want:       2,
		},
	}

	for _, tt := range tests {
		got := detectSoftIndentWidth(tt.src, tt.javadocPos)
		assert.Equal(t, tt.want, got)
	}
}
