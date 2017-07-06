package lint

import "bytes"

var metadataNoEOLSpaceRule = &MetaData{
	Name: "no-eol-space",
	rank: 4, // this rule should be called before first-newline and final-newline
}

type NoEOLSpaceRule struct{}

func NewNoEOLSpaceRule(ops map[string]interface{}) (Rule, error) {
	return &NoEOLSpaceRule{}, nil
}

func (r *NoEOLSpaceRule) New(ops map[string]interface{}) (Rule, error) {
	return NewNoEOLSpaceRule(ops)
}

func (r *NoEOLSpaceRule) MetaData() *MetaData {
	return metadataNoEOLSpaceRule
}

func (r *NoEOLSpaceRule) Lint(s []byte) (*Result, error) {
	res := NewResult()
	errmsg := "Trailing spaces/tabs at the end of lines are disallowed"

	linebreak := detectLinebreakStyle(s)

	ls := bytes.Split(s, linebreak)
	for i, l := range ls {
		ls[i] = bytes.TrimRight(l, " \t")
		if bytes.HasSuffix(l, []byte(" ")) || bytes.HasSuffix(l, []byte("\t")) {
			res.AddReport(i+1, -1, errmsg)
		}
	}
	res.Set(bytes.Join(ls, linebreak))

	return res, nil
}

func init() {
	DefinedRules.Set(&NoEOLSpaceRule{})
}
