package lint

import (
	"io/ioutil"
	"sort"
)

type Linter struct {
	filename string
	source   []byte
	rules    RankedRules
}

type RankedRules []Rule

func (r RankedRules) Len() int {
	return len(r)
}

func (r RankedRules) Less(i int, j int) bool {
	return r[i].MetaData().rank < r[j].MetaData().rank
}

func (r RankedRules) Swap(i int, j int) {
	r[i], r[j] = r[j], r[i]
}

func NewLinter(filename string, rules []Rule) (*Linter, error) {
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	rs := RankedRules(rules)
	sort.Sort(rs)

	linter := &Linter{
		filename: filename,
		source:   src,
		rules:    rs,
	}

	return linter, nil
}

func (linter *Linter) Lint() (*Result, error) {
	result := NewResult()
	src := make([]byte, len(linter.source))
	copy(src, linter.source)

	if len(linter.source) != 0 {
		for _, rule := range linter.rules {
			r, err := rule.Lint(src)
			if err != nil {
				return nil, err
			}
			result.Reports = append(result.Reports, r.Reports...)
			src = r.Fixed
		}
	}

	result.Set(src)

	return result, nil
}
