package lint

import "bytes"

var metadataNoBOMRule = &MetaData{
	Name: "no-bom",
	rank: 5,
}

type NoBOMRule struct{}

func NewNoBOMRule(ops map[string]interface{}) (Rule, error) {
	return &NoBOMRule{}, nil
}

func (r *NoBOMRule) New(ops map[string]interface{}) (Rule, error) {
	return NewNoBOMRule(ops)
}

func (r *NoBOMRule) MetaData() *MetaData {
	return metadataNoBOMRule
}

func (r *NoBOMRule) Lint(s []byte) (*Result, error) {
	res := NewResult()
	errmsg := "Byte order mark is disallowed"

	if !bytes.HasPrefix(s, UTF8BOMs) {
		res.Set(s)
		return res, nil
	}

	res.AddReport(0, 0, errmsg)
	res.Set(bytes.TrimPrefix(s, UTF8BOMs))

	return res, nil
}

var UTF8BOMs = []byte{0xEF, 0xbb, 0xbf}

func init() {
	definedRules.Set(&NoBOMRule{})
}
