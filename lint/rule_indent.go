package lint

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

var metadataIndentRule = &MetaData{
	Name: "indent",
	rank: 5,
}

var (
	ErrUnknownIndentStyle = errors.New("unknown indent style")
)

type IndentStyle byte

const (
	IndentStyleSoft IndentStyle = ' '
	IndentStyleHard IndentStyle = '\t'
)

func NewIndentStyle(str string) (IndentStyle, error) {
	switch strings.ToLower(str) {
	case "soft", "space":
		return IndentStyleSoft, nil
	case "hard", "tab":
		return IndentStyleHard, nil
	default:
		return 0, ErrUnknownIndentStyle
	}
}

type IndentRule struct {
	Style IndentStyle
	Size  int
}

func NewIndentRule(ops map[string]interface{}) (Rule, error) {
	rule := &IndentRule{}

	if v, ok := ops["style"]; ok {
		if value, ok := v.(string); ok {
			style, err := NewIndentStyle(value)
			if err != nil {
				return nil, fmt.Errorf("indent.style is invalid: %v: %q", err, value)
			}
			rule.Style = style
		} else {
			return nil, fmt.Errorf("indent.style is invalid: %v: %v", ErrUnknownIndentStyle, v)
		}
	}

	if v, ok := ops["size"]; ok {
		if value, ok := v.(int); ok {
			rule.Size = value
		} else {
			return nil, fmt.Errorf("indent.size is only allow numbers: %v", v)
		}
	}

	return rule, nil
}

func (r *IndentRule) New(ops map[string]interface{}) (Rule, error) {
	return NewIndentRule(ops)
}

func (r *IndentRule) MetaData() *MetaData {
	return metadataIndentRule
}

func (r *IndentRule) Lint(src []byte) (*Result, error) {
	res := NewResult()

	var indent []byte
	var expectMsg string
	switch r.Style {
	case IndentStyleSoft:
		expectMsg = fmt.Sprintf(`Expected indent with %d space(s)`, r.Size)
		indent = bytes.Repeat([]byte{byte(r.Style)}, r.Size)
	case IndentStyleHard:
		expectMsg = `Expected indent with hardtabs (\t)`
		indent = []byte{byte(r.Style)}
	}

	linebreak := detectLinebreakStyle(src)
	lines := bytes.Split(src, linebreak)

	javadocPos := searchJavadocComments(lines)
	softIndentWidth := detectSoftIndentWidth(lines, javadocPos)

	for i, line := range lines {
		col := i + 1
		hasJavadocSpace := hasJavadocSpace(col, javadocPos)

		var depth int
		if hasJavadocSpace {
			idx := bytes.Index(line, []byte(" *"))
			depth = detectIndentDepth(line[:idx], softIndentWidth)
		} else {
			depth = detectIndentDepth(line, softIndentWidth)
		}

		if depth == 0 {
			continue
		}

		line = bytes.TrimLeft(line, " \t")
		indentLength := r.Size * depth

		if hasJavadocSpace {
			indentLength += 1
		}

		newLine := make([]byte, 0, len(line)+indentLength)
		indentBytes := bytes.Repeat(indent, depth)

		newLine = append([]byte{}, indentBytes...)

		if hasJavadocSpace {
			newLine = append(newLine, ' ')
		}

		newLine = append(newLine, line...)

		if len(newLine) != len(lines[i]) {
			var errmsg string
			if bytes.HasPrefix(lines[i], []byte("\t")) {
				errmsg = fmt.Sprintf(`%s but used hardtabs (\t)`, expectMsg)
			} else if bytes.HasPrefix(lines[i], []byte(" ")) {
				errmsg = fmt.Sprintf(`%s but used %d space(s)`, expectMsg, softIndentWidth)
			}
			res.AddReport(col, -1, errmsg)
		}

		lines[i] = newLine
	}
	res.Set(bytes.Join(lines, linebreak))

	return res, nil
}

func searchJavadocComments(lines [][]byte) []*columnRange {
	var crs []*columnRange
	var cr *columnRange
	var inJavadoc bool

	for i := 1; i < len(lines); i++ {
		bef := i - 1
		beforeLine := bytes.TrimLeft(lines[bef], " \t")
		currentLine := bytes.TrimLeft(lines[i], " \t")

		if bytes.HasPrefix(beforeLine, []byte("/**")) && bytes.HasPrefix(currentLine, []byte("*")) {
			inJavadoc = true
			cr = &columnRange{
				begin: bef + 1,
			}
			continue
		}

		if inJavadoc {
			if bytes.HasPrefix(beforeLine, []byte("*")) {
				if bytes.HasPrefix(currentLine, []byte("*/")) {
					inJavadoc = false
					cr.end = i + 1
				} else if !bytes.HasPrefix(beforeLine, []byte("*")) {
					inJavadoc = false
					cr = nil
				}
			}
		}

		if !inJavadoc && cr != nil {
			crs = append(crs, cr)
			cr = nil
		}
	}

	return crs
}

func hasJavadocSpace(col int, javadocPos []*columnRange) bool {
	for _, pos := range javadocPos {
		p := &columnRange{
			begin: pos.begin + 1, // +1 means skip the `/**` line (starts javadoc comments)
			end:   pos.end,
		}

		// before line is ` *` or ` */`
		if p.in(col) {
			return true
		}
	}
	return false
}

func detectIndentDepth(line []byte, softIndentWidth int) int {
	n := 0
	softCount := 0

	for _, c := range line {
		switch c {
		case '\t':
			n++
		case ' ':
			if softIndentWidth > 0 {
				softCount++
				if softCount >= softIndentWidth {
					n++
					softCount = 0
				}
			}
		default:
			return n
		}
	}

	return n
}

func detectSoftIndentWidth(lines [][]byte, javadocPos []*columnRange) int {
	predict := map[int]int{} // { indentWidth: frequency }

	for i := 1; i < len(lines); i++ {
		beforeLine := lines[i-1]
		currentLine := lines[i]

		if len(beforeLine) == 0 || len(currentLine) == 0 { // if line has no characters
			continue
		}

		col := i + 1

		// [1:] means truncate first one space (` *`, ` */`) for not recognizing as an indentation space
		if hasJavadocSpace(col-1, javadocPos) {
			beforeLine = beforeLine[1:]
		}
		if hasJavadocSpace(col, javadocPos) {
			currentLine = currentLine[1:]
		}

		before := countFirstSpaces(beforeLine)
		current := countFirstSpaces(currentLine)

		if before == current {
			continue
		}

		var diff int
		if current > before {
			diff = current - before
		} else {
			diff = before - current
		}
		predict[diff]++
	}

	var max, width int
	for n, freq := range predict {
		if freq > max {
			max = freq
			width = n
		} else if freq == max {
			width = 0
		}
	}

	return width
}

func countFirstSpaces(cs []byte) int {
	n := 0

	for _, c := range cs {
		if c != ' ' {
			break
		}
		n++
	}
	return n
}

func init() {
	DefinedRules.Set(&IndentRule{})
}
