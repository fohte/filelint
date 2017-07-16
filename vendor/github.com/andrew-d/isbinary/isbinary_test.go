package isbinary

import (
	"bytes"
	"fmt"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func toUTF8(s string) []byte {
	out := []byte{}
	for _, char := range s {
		var curr [utf8.UTFMax]byte
		n := utf8.EncodeRune(curr[:], char)
		out = append(out, curr[:n]...)
	}

	return out
}

var (
	testCases = []struct {
		Input  []byte
		Result bool
	}{
		{[]byte("some text"), false},
		{[]byte("some text with a \x00 char"), true},
		{[]byte("\xEF\xBB\xBF text with utf-8 BOM"), false},
		{[]byte("text with suspicious \xFF\xFF\xFF\xFF\xFF\xFF\xFF"), true},
		{toUTF8("utf8 text  世界"), false},
		{toUTF8("utf8 世界 with null \x00"), true},
		{toUTF8("utf8 世界 with suspicious \x01\x01\x01\x01\x01"), true},
	}
)

func TestIsBinary(t *testing.T) {
	for i, tcase := range testCases {
		assert.Equal(t, tcase.Result, Test(tcase.Input),
			"test case %d did not have expected result", i)
	}
}

func TestIsBinaryReader(t *testing.T) {
	for i, tcase := range testCases {
		r := bytes.NewReader(tcase.Input)
		res, err := TestReader(r)
		assert.NoError(t, err)
		assert.Equal(t, tcase.Result, res,
			"test case %d did not have expected result", i)
	}
}

func TestIsBinaryUTF8(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	// TODO: when we add support for four-byte UTF-8 characters, bump the upper
	// limit to 0x10FFFF
	const (
		// Start after the ASCII characters
		UTF8CharMin = 0x80

		// Upper limit of two- or three-byte UTF-8 characters.
		UTF8CharMax = 0xFFFF
	)

	// Run through all possible UTF-8 characters and verify that they don't get
	// detected as binary.
	var buf [utf8.UTFMax]byte
	for i := UTF8CharMin; i <= UTF8CharMax; i++ {
		if i%(UTF8CharMax/10) == 0 {
			fmt.Printf("running test case 0x%06x/0x%06x\n", i, UTF8CharMax)
		}

		r := rune(i)
		n := utf8.EncodeRune(buf[:], r)
		assert.False(t, Test(buf[:n]),
			"encoding rune %U should be detected as non-binary", r)
	}
}
