package lint

import (
	"fmt"
	"strconv"
)

var definedRules = NewRuleMap()

func GetDefinedRules() *RuleMap {
	return definedRules
}

type RuleMap struct {
	m map[string]Rule
}

func NewRuleMap(rules ...Rule) *RuleMap {
	m := make(map[string]Rule)
	rmap := &RuleMap{m: m}
	for _, rule := range rules {
		rmap.Set(rule)
	}
	return rmap
}

func (rmap *RuleMap) Set(rule Rule) {
	rmap.m[rule.MetaData().Name] = rule
}

func (rmap *RuleMap) Get(ruleName string) (r Rule) {
	r, _ = rmap.m[ruleName]
	return r
}

func (rmap *RuleMap) Has(ruleName string) bool {
	_, ok := rmap.m[ruleName]
	return ok
}

func (rmap *RuleMap) Size() int {
	return len(rmap.m)
}

func (rmap *RuleMap) GetAllRuleNames() []string {
	names := make([]string, len(rmap.m))
	for name, _ := range rmap.m {
		names = append(names, name)
	}
	return names
}

type Rule interface {
	New(ops map[string]interface{}) (Rule, error)
	MetaData() *MetaData
	Lint(s []byte) (*Result, error)
}

type MetaData struct {
	Name        string
	Description string

	// rank is the order called by linter.
	// this is evaluated in descending order.
	rank int
}

type Result struct {
	Fixed   []byte
	Reports []*Report
}

func NewResult() *Result {
	var reps []*Report
	reps = []*Report{}
	return &Result{
		Reports: reps,
	}
}

func (res *Result) Set(fixed []byte) {
	if fixed == nil {
		res.Fixed = []byte{}
		return
	}

	res.Fixed = fixed
}

func (res *Result) AddReport(col, row int, message string) {
	res.Reports = append(res.Reports, NewReport(col, row, message))
}

type Report struct {
	position *Position
	message  string
}

func NewReport(col, row int, message string) *Report {
	return &Report{
		position: &Position{col, row},
		message:  message,
	}
}

func (rep *Report) String() string {
	return fmt.Sprintf("%s: %s", rep.position.String(), rep.message)
}

type Position struct {
	column int
	row    int
}

func (pos *Position) String() string {
	var x, y string

	x = strconv.Itoa(pos.row)
	y = strconv.Itoa(pos.column)

	return fmt.Sprintf("%s:%s", y, x)
}

type columnRange struct {
	begin int
	end   int
}

func (cr *columnRange) in(col int) bool {
	return cr.begin <= col && col <= cr.end
}
