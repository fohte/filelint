package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_Merge(t *testing.T) {
	tests := []struct {
		src  *Config
		dst  *Config
		want *Config
	}{
		{
			src:  &Config{File: File{Include: []string{"a"}}},
			dst:  &Config{File: File{Include: []string{"b"}}},
			want: &Config{File: File{Include: []string{"b"}}},
		},
		{
			src:  &Config{File: File{Include: []string{}}},
			dst:  &Config{File: File{Include: []string{"b"}}},
			want: &Config{File: File{Include: []string{"b"}}},
		},
		{
			src:  &Config{File: File{Include: []string{"a"}}},
			dst:  &Config{File: File{Include: []string{}}},
			want: &Config{File: File{Include: []string{"a"}}},
		},
		{
			src:  &Config{File: File{Exclude: []string{"a"}}},
			dst:  &Config{File: File{Exclude: []string{"b"}}},
			want: &Config{File: File{Exclude: []string{"a", "b"}}},
		},
		{
			src:  &Config{File: File{Exclude: []string{}}},
			dst:  &Config{File: File{Exclude: []string{"b"}}},
			want: &Config{File: File{Exclude: []string{"b"}}},
		},
		{
			src:  &Config{File: File{Exclude: []string{"a"}}},
			dst:  &Config{File: File{Exclude: []string{}}},
			want: &Config{File: File{Exclude: []string{"a"}}},
		},
		{
			src: &Config{Targets: []Target{Target{
				Patterns: []string{"*"},
				Rule:     RuleMap{},
			}}},
			dst: &Config{Targets: []Target{Target{
				Patterns: []string{"*"},
				Rule:     RuleMap{},
			}}},
			want: &Config{Targets: []Target{
				{
					Patterns: []string{"*"},
					Rule:     RuleMap{},
				},
				{
					Patterns: []string{"*"},
					Rule:     RuleMap{},
				},
			}},
		},
	}

	for _, tt := range tests {
		tt.src.Merge(tt.dst)
		assert.Equal(t, tt.want, tt.src)
	}
}

func TestConfig_MatchedRule(t *testing.T) {
	tests := []struct {
		src  []Target
		file string
		want RuleMap
	}{
		{
			src: []Target{
				{
					Patterns: []string{"**/*"},
					Rule:     RuleMap{"a": {"A": 1}},
				},
			},
			file: "path/to/a",
			want: RuleMap{"a": {"A": 1}},
		},
		{
			src: []Target{
				{
					Patterns: []string{"**/*"},
					Rule:     RuleMap{"a": {"A": 1}},
				},
				{
					Patterns: []string{"**/*.go"},
					Rule:     RuleMap{"a": {"A": 2}},
				},
			},
			file: "path/to/a",
			want: RuleMap{"a": {"A": 1}},
		},
		{
			src: []Target{
				{
					Patterns: []string{"**/*"},
					Rule:     RuleMap{"a": {"A": 1}},
				},
				{
					Patterns: []string{"**/*.go"},
					Rule:     RuleMap{"a": {"A": 2}},
				},
			},
			file: "path/to/a.go",
			want: RuleMap{"a": {"A": 2}},
		},
		{
			src: []Target{
				{
					Patterns: []string{"**/*"},
					Rule:     RuleMap{"a": {"A": 1}},
				},
				{
					Patterns: []string{"**/*"},
					Rule:     RuleMap{"a": {"A": 2}},
				},
			},
			file: "path/to/a",
			want: RuleMap{"a": {"A": 2}},
		},
	}

	for _, tt := range tests {
		c := &Config{Targets: tt.src}
		got := c.MatchedRule(tt.file)
		assert.Equal(t, tt.want, got)
	}
}

func TestRuleMap_Merge(t *testing.T) {
	tests := []struct {
		src  RuleMap
		dst  RuleMap
		want RuleMap
	}{
		{
			src:  RuleMap{"a": {"op": 1}},
			dst:  RuleMap{"a": {"op": 2}},
			want: RuleMap{"a": {"op": 2}},
		},
		{
			src:  RuleMap{"a": {"op": 1}},
			dst:  RuleMap{"a": {}},
			want: RuleMap{"a": {"op": 1}},
		},
		{
			src:  RuleMap{"a": {}},
			dst:  RuleMap{"a": {"op": 1}},
			want: RuleMap{"a": {"op": 1}},
		},
		{
			src:  RuleMap{},
			dst:  RuleMap{"a": {"op": 1}},
			want: RuleMap{"a": {"op": 1}},
		},
		{
			src:  RuleMap{"a": {"op": 1}},
			dst:  RuleMap{},
			want: RuleMap{"a": {"op": 1}},
		},
		{
			src:  RuleMap{"a": {"op": false}},
			dst:  RuleMap{"a": {"op": true}},
			want: RuleMap{"a": {"op": true}},
		},
		{
			src:  RuleMap{"a": {"op": true}},
			dst:  RuleMap{"a": {"op": false}},
			want: RuleMap{"a": {"op": false}},
		},
		{
			src:  RuleMap{"a": {"op": true}},
			dst:  RuleMap{"b": {"op": false}},
			want: RuleMap{"a": {"op": true}, "b": {"op": false}},
		},
	}

	for _, tt := range tests {
		got := tt.src.Merge(tt.dst)
		assert.Equal(t, tt.want, got)
	}
}
