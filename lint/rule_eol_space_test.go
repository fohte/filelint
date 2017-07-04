package lint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoEOLSpaceRule_Lint(t *testing.T) {
	tests := []struct {
		rule NoEOLSpaceRule
		src  []byte
		want []byte
	}{
		{
			rule: NoEOLSpaceRule{},
			src:  []byte(" "),
			want: []byte(""),
		},
		{
			rule: NoEOLSpaceRule{},
			src:  []byte("  "),
			want: []byte(""),
		},
		{
			rule: NoEOLSpaceRule{},
			src:  []byte("\t"),
			want: []byte(""),
		},
		{
			rule: NoEOLSpaceRule{},
			src:  []byte(" \t"),
			want: []byte(""),
		},
		{
			rule: NoEOLSpaceRule{},
			src:  []byte("\t "),
			want: []byte(""),
		},
		{
			rule: NoEOLSpaceRule{},
			src:  []byte(" \n "),
			want: []byte("\n"),
		},
		{
			rule: NoEOLSpaceRule{},
			src:  []byte(" \r\n "),
			want: []byte("\r\n"),
		},
		{
			rule: NoEOLSpaceRule{},
			src:  []byte(" . "),
			want: []byte(" ."),
		},
	}

	for _, tt := range tests {
		got, _ := tt.rule.Lint(tt.src)
		assert.Equal(t, tt.want, got.Fixed)
	}
}
