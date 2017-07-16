package cli

import "strings"

type ruleSliceValue struct {
	value []*ruleValue
}

func (rs *ruleSliceValue) String() string {
	strs := make([]string, 0, len(rs.value))
	for _, v := range rs.value {
		strs = append(strs, v.String())
	}
	return strings.Join(strs, ", ")
}

func (rs *ruleSliceValue) Set(str string) error {
	r := newRuleValue()
	if err := r.Set(str); err != nil {
		return err
	}
	rs.value = append(rs.value, r)
	return nil
}

func (rs *ruleSliceValue) Type() string {
	return "ruleSlice"
}
