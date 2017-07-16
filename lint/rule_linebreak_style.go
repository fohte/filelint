package lint

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

var metadataLinebreakRule = &MetaData{
	Name: "linebreak",

	// this rule should called before all rules
	rank: 0,
}

var (
	ErrUnknownLinebreakStyle = errors.New("unknown line break style")
)

type LinebreakStyle []byte

var (
	UnixStyleLinebreak    = LinebreakStyle{'\n'}
	WindowsStyleLinebreak = LinebreakStyle{'\r', '\n'}
)

func NewLinebreakStyle(str string) (LinebreakStyle, error) {
	switch strings.ToLower(str) {
	case "lf":
		return UnixStyleLinebreak, nil
	case "crlf":
		return WindowsStyleLinebreak, nil
	}
	return nil, ErrUnknownLinebreakStyle
}

func (s LinebreakStyle) text() string {
	switch {
	case bytes.Equal(s, UnixStyleLinebreak):
		return "LF"
	case bytes.Equal(s, WindowsStyleLinebreak):
		return "CRLF"
	}
	return ""
}

type LinebreakRule struct {
	Style LinebreakStyle
}

func NewLinebreakRule(ops map[string]interface{}) (Rule, error) {
	rule := &LinebreakRule{}

	if v, ok := ops["style"]; ok {
		if value, ok := v.(string); ok {
			style, err := NewLinebreakStyle(value)
			if err != nil {
				return nil, fmt.Errorf("linebreak.style is invalid: %v: %q", err, value)
			}
			rule.Style = style
		} else {
			return nil, fmt.Errorf("linebreak.style is invalid: %v: %q", ErrUnknownLinebreakStyle, value)
		}
	}

	return rule, nil
}

func (r *LinebreakRule) New(ops map[string]interface{}) (Rule, error) {
	return NewLinebreakRule(ops)
}

func (r *LinebreakRule) MetaData() *MetaData {
	return metadataLinebreakRule
}

func (r *LinebreakRule) Lint(s []byte) (*Result, error) {
	res := NewResult()

	formatTarget := detectLinebreakStyle(s)
	if bytes.Equal(formatTarget, r.Style) {
		res.Set(s)
		return res, nil
	}

	errmsg := fmt.Sprintf(
		`Expected linebreaks to be %s but found %s`,
		r.Style.text(),
		formatTarget.text(),
	)
	res.AddReport(-1, -1, errmsg)
	res.Set(bytes.Replace(s, formatTarget, r.Style, -1))

	return res, nil
}

func detectLinebreakStyle(bs []byte) LinebreakStyle {
	for _, b := range bs {
		if b == '\r' {
			return WindowsStyleLinebreak
		}
	}

	return UnixStyleLinebreak
}

func init() {
	definedRules.Set(&LinebreakRule{})
}
