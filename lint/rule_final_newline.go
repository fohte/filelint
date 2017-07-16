package lint

import (
	"bytes"
	"fmt"
)

var metadataFinalNewline = &MetaData{
	Name: "final-newline",
	rank: 5,
}

type FinalNewlineRule struct {
	Num int
}

func NewFinalNewlineRule(ops map[string]interface{}) (Rule, error) {
	rule := &FinalNewlineRule{}

	if v, ok := ops["num"]; ok {
		if value, ok := v.(int); ok {
			rule.Num = value
		} else {
			return nil, fmt.Errorf("final-newline.num is only allow numbers: %v", v)
		}
	}

	return rule, nil
}

func (r *FinalNewlineRule) New(ops map[string]interface{}) (Rule, error) {
	return NewFinalNewlineRule(ops)
}

func (r *FinalNewlineRule) MetaData() *MetaData {
	return metadataFinalNewline
}

func (r *FinalNewlineRule) Lint(s []byte) (*Result, error) {
	res := NewResult()

	n := countFinalNewlines(s)
	linebreak := detectLinebreakStyle(s)

	if n == r.Num {
		res.Set(s)
		return res, nil
	}

	errmsg := fmt.Sprintf("Files should end with %d newline(s) but %d newline(s)", r.Num, n)
	res.AddReport(0, 0, errmsg)

	trimmed := bytes.TrimRight(s, string(linebreak))

	if r.Num == 0 {
		res.Set(trimmed)
		return res, nil
	}

	appended := make([]byte, 0, len(trimmed)+r.Num*len(linebreak))
	copy(appended, trimmed)

	lines := bytes.Repeat(linebreak, r.Num)
	appended = append([]byte{}, trimmed...)
	appended = append(appended, lines...)

	res.Set(appended)

	return res, nil
}

func countFinalNewlines(s []byte) int {
	n := 0

	for i := len(s) - 1; i >= 0; i-- {
		c := s[i]

		switch c {
		case '\n':
			n++
		case '\r':
			continue
		default:
			return n
		}
	}

	return n
}

func init() {
	definedRules.Set(&FinalNewlineRule{})
}
