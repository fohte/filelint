package lint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFinalNewlineRule_Lint(t *testing.T) {
	tests := []struct {
		rule FinalNewlineRule
		src  []byte
		want []byte
	}{
		{
			rule: FinalNewlineRule{Num: 0},
			src:  []byte("\n\n"),
			want: []byte(""),
		},
		{
			rule: FinalNewlineRule{Num: 0},
			src:  []byte("\r\n\r\n"),
			want: []byte(""),
		},
		{
			rule: FinalNewlineRule{Num: 0},
			src:  []byte(".\n"),
			want: []byte("."),
		},
		{
			rule: FinalNewlineRule{Num: 0},
			src:  []byte(".\r\n"),
			want: []byte("."),
		},
		{
			rule: FinalNewlineRule{Num: 0},
			src:  []byte("\n.\n"),
			want: []byte("\n."),
		},
		{
			rule: FinalNewlineRule{Num: 0},
			src:  []byte("\r\n.\r\n"),
			want: []byte("\r\n."),
		},
		{
			rule: FinalNewlineRule{Num: 1},
			src:  []byte("\n\n"),
			want: []byte("\n"),
		},
		{
			rule: FinalNewlineRule{Num: 1},
			src:  []byte("\r\n\r\n"),
			want: []byte("\r\n"),
		},
		{
			rule: FinalNewlineRule{Num: 1},
			src:  []byte("."),
			want: []byte(".\n"),
		},
		{
			rule: FinalNewlineRule{Num: 1},
			src:  []byte("\n."),
			want: []byte("\n.\n"),
		},
		{
			rule: FinalNewlineRule{Num: 1},
			src:  []byte("\r\n."),
			want: []byte("\r\n.\r\n"),
		},
		{
			rule: FinalNewlineRule{Num: 1},
			src:  []byte("\n.\n"),
			want: []byte("\n.\n"),
		},
		{
			rule: FinalNewlineRule{Num: 1},
			src:  []byte("\r\n.\r\n"),
			want: []byte("\r\n.\r\n"),
		},
	}

	for _, tt := range tests {
		got, _ := tt.rule.Lint(tt.src)
		assert.Equal(t, tt.want, got.Fixed)
	}
}
