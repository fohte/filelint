package cli

import (
	"errors"
	"fmt"
	"strings"
)

const (
	optionSep = ":"
	ruleSep   = "="
	multiSep  = ","
)

type ruleValue struct {
	name    string
	options map[string]interface{}
}

func newRuleValue() *ruleValue {
	m := make(map[string]interface{})
	return &ruleValue{
		options: m,
	}
}

func (r *ruleValue) String() string {
	if len(r.options) == 0 {
		return r.name
	}

	ops := make([]string, 0, len(r.options))
	for k, v := range r.options {
		ops = append(ops, fmt.Sprintf("%s%s%s", k, optionSep, v))
	}
	return fmt.Sprintf("%s%s%s", r.name, ruleSep, strings.Join(ops, multiSep))
}

func (r *ruleValue) Set(str string) error {
	rule, err := parseRuleDSL(str)
	if err != nil {
		return err
	}
	*r = *rule
	return nil
}

func (r *ruleValue) Type() string {
	return "rule"
}

func parseRuleDSL(str string) (*ruleValue, error) {
	s := strings.SplitN(str, ruleSep, 2)

	r := newRuleValue()
	r.name = s[0]

	if len(s) == 1 {
		return r, nil
	}

	ops := strings.Split(s[1], multiSep)
	for _, op := range ops {
		o := strings.SplitN(op, optionSep, 2)
		if len(o) < 2 {
			return nil, errors.New(fmt.Sprintf("options should be `key%svalue`", optionSep))
		}
		r.options[o[0]] = o[1]
	}

	return r, nil
}
