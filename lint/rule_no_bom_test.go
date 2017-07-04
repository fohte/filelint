package lint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoBOMRule_Lint(t *testing.T) {
	tests := []struct {
		rule NoBOMRule
		src  []byte
		want []byte
	}{
		{
			rule: NoBOMRule{},
			src:  []byte{0xEF, 0xBB, 0xBF},
			want: []byte(""),
		},
		{
			rule: NoBOMRule{},
			src:  []byte{0xBB, 0xEF, 0xBF},
			want: []byte{0xBB, 0xEF, 0xBF},
		},
		{
			rule: NoBOMRule{},
			src:  []byte{0xBB, 0xBF},
			want: []byte{0xBB, 0xBF},
		},
		{
			rule: NoBOMRule{},
			src:  []byte{0xBF},
			want: []byte{0xBF},
		},
		{
			rule: NoBOMRule{},
			src:  []byte{' ', 0xEF, 0xBB, 0xBF},
			want: []byte{' ', 0xEF, 0xBB, 0xBF},
		},
	}

	for _, tt := range tests {
		got, _ := tt.rule.Lint(tt.src)
		assert.Equal(t, tt.want, got.Fixed)
	}
}
