package lint

import (
	"bytes"
	"fmt"
)

var metadataFirstNewlineRule = &MetaData{
	Name: "first-newline",

	// this rule should called before final-newline rule
	rank: 2,
}

type FirstNewlineRule struct {
	Num int
}

func NewFirstNewlineRule(ops map[string]interface{}) (Rule, error) {
	rule := &FirstNewlineRule{}

	if v, ok := ops["num"]; ok {
		if value, ok := v.(int); ok {
			rule.Num = value
		} else {
			return nil, fmt.Errorf("first-newline.num is only allow numbers: %v", v)
		}
	}

	return rule, nil
}

func (r *FirstNewlineRule) New(ops map[string]interface{}) (Rule, error) {
	return NewFirstNewlineRule(ops)
}

func (r *FirstNewlineRule) MetaData() *MetaData {
	return metadataFirstNewlineRule
}

func (r *FirstNewlineRule) Lint(s []byte) (*Result, error) {
	res := NewResult()

	n := countFirstNewlines(s)
	linebreak := detectLinebreakStyle(s)

	if n == r.Num {
		res.Set(s)
		return res, nil
	}

	errmsg := fmt.Sprintf("Files should begin with %d newline(s) but %d newline(s)", r.Num, n)
	res.AddReport(-1, -1, errmsg)

	trimmed := bytes.TrimLeft(s, string(linebreak))

	if r.Num == 0 {
		res.Set(trimmed)
		return res, nil
	}

	appended := make([]byte, 0, len(trimmed)+r.Num*len(linebreak))
	copy(appended, trimmed)

	lines := bytes.Repeat(linebreak, r.Num)
	appended = append([]byte{}, lines...)
	appended = append(appended, trimmed...)

	res.Set(appended)

	return res, nil
}

func countFirstNewlines(s []byte) int {
	n := 0

	for _, c := range s {
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
	definedRules.Set(&FirstNewlineRule{})
}
