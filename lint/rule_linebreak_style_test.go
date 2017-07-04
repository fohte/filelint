package lint

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLinebreakStyle(t *testing.T) {
	tests := []struct {
		key     string
		want    LinebreakStyle
		wanterr error
	}{
		{"lf", UnixStyleLinebreak, nil},
		{"LF", UnixStyleLinebreak, nil},
		{"crlf", WindowsStyleLinebreak, nil},
		{"CRLF", WindowsStyleLinebreak, nil},
		{"lflf", nil, ErrUnknownLinebreakStyle},
	}

	for _, tt := range tests {
		got, err := NewLinebreakStyle(tt.key)
		if !bytes.Equal(tt.want, got) {
			t.Errorf("NewLinebreakStyle(%q) == %q (got: %q)", tt.key, tt.want, got)
		}
		if err != tt.wanterr {
			t.Errorf("NewLinebreakStyle(%q) should throw %v (got: %v)", tt.key, tt.wanterr, err)
		}
	}
}

func TestLinebreakRule_Lint(t *testing.T) {
	tests := []struct {
		rule LinebreakRule
		src  []byte
		want []byte
	}{
		{
			LinebreakRule{Style: UnixStyleLinebreak},
			[]byte("\r\n"),
			[]byte("\n"),
		},
		{
			LinebreakRule{Style: UnixStyleLinebreak},
			[]byte("\n"),
			[]byte("\n"),
		},
		{
			LinebreakRule{Style: WindowsStyleLinebreak},
			[]byte("\n"),
			[]byte("\r\n"),
		},
		{
			LinebreakRule{Style: WindowsStyleLinebreak},
			[]byte("\r\n"),
			[]byte("\r\n"),
		},
	}

	for _, tt := range tests {
		got, _ := tt.rule.Lint(tt.src)
		assert.Equal(t, tt.want, got.Fixed)
	}
}
