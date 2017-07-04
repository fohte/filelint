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
			src:  &Config{Targets: map[string]Target{"default": {}}},
			dst:  &Config{Targets: map[string]Target{"go": {}}},
			want: &Config{Targets: map[string]Target{"default": {}, "go": {}}},
		},
		{
			src:  &Config{Targets: map[string]Target{}},
			dst:  &Config{Targets: map[string]Target{"go": {}}},
			want: &Config{Targets: map[string]Target{"go": {}}},
		},
		{
			src:  &Config{Targets: map[string]Target{"default": {}}},
			dst:  &Config{Targets: map[string]Target{}},
			want: &Config{Targets: map[string]Target{"default": {}}},
		},
	}

	for _, tt := range tests {
		tt.src.Merge(tt.dst)
		assert.Equal(t, tt.want, tt.src)
	}
}

func TestTargetMap_ExtendDefaultTarget(t *testing.T) {
	tests := []struct {
		msg  string
		src  TargetMap
		want TargetMap
	}{
		{
			msg: "non `default` existing values should not overwrite with `default` values",
			src: TargetMap{
				"default": Target{Rule: RuleMap{"a": {"op": "default"}}},
				"go":      Target{Rule: RuleMap{"a": {"op": "go"}}},
			},
			want: TargetMap{
				"default": Target{Rule: RuleMap{"a": {"op": "default"}}},
				"go":      Target{Rule: RuleMap{"a": {"op": "go"}}},
			},
		},
		{
			msg: "if non `default` dosen't have values then should be set with `default` values",
			src: TargetMap{
				"default": Target{Rule: RuleMap{"a": {"A": 1, "B": 1}}},
				"go":      Target{Rule: RuleMap{"a": {"B": 2}}},
			},
			want: TargetMap{
				"default": Target{Rule: RuleMap{"a": {"A": 1, "B": 1}}},
				"go":      Target{Rule: RuleMap{"a": {"A": 1, "B": 2}}},
			},
		},
		{
			msg: "if non `default` dosen't have values then should be set with `default` values",
			src: TargetMap{
				"default": Target{Rule: RuleMap{"a": {"op": "default"}}},
				"go":      Target{Rule: RuleMap{}},
			},
			want: TargetMap{
				"default": Target{Rule: RuleMap{"a": {"op": "default"}}},
				"go":      Target{Rule: RuleMap{"a": {"op": "default"}}},
			},
		},
		{
			msg: "if non `default` dosen't have values then should be set with `default` values",
			src: TargetMap{
				"default": Target{Rule: RuleMap{"a": {"op": "default"}}},
				"go":      Target{Rule: RuleMap{"a": {}}},
			},
			want: TargetMap{
				"default": Target{Rule: RuleMap{"a": {"op": "default"}}},
				"go":      Target{Rule: RuleMap{"a": {"op": "default"}}},
			},
		},
		{
			src: TargetMap{
				"go": Target{Rule: RuleMap{"a": {"op": "go"}}},
			},
			want: TargetMap{
				"go": Target{Rule: RuleMap{"a": {"op": "go"}}},
			},
		},
	}

	for _, tt := range tests {
		tt.src.ExtendDefaultTarget()
		assert.Equal(t, tt.want, tt.src, tt.msg)
	}
}

func TestTarget_Merge(t *testing.T) {
	tests := []struct {
		src  Target
		dst  Target
		want Target
	}{
		{
			src:  Target{Pattern: []string{"*"}},
			dst:  Target{Pattern: []string{"**"}},
			want: Target{Pattern: []string{"**"}},
		},
		{
			src:  Target{Pattern: []string{"*", "*"}},
			dst:  Target{Pattern: []string{"**"}},
			want: Target{Pattern: []string{"**"}},
		},
		{
			src:  Target{Pattern: []string{"*"}},
			dst:  Target{Pattern: []string{"**", "**"}},
			want: Target{Pattern: []string{"**", "**"}},
		},
		{
			src:  Target{Pattern: []string{}},
			dst:  Target{Pattern: []string{"*"}},
			want: Target{Pattern: []string{"*"}},
		},
		{
			src:  Target{Pattern: []string{"*"}},
			dst:  Target{Pattern: []string{}},
			want: Target{Pattern: []string{"*"}},
		},
		{
			src:  Target{Rule: RuleMap{"a": {"op": 1}}},
			dst:  Target{Rule: RuleMap{"a": {"op": 2}}},
			want: Target{Rule: RuleMap{"a": {"op": 2}}},
		},
		{
			src:  Target{Rule: RuleMap{"a": {"op": 1}}},
			dst:  Target{Rule: RuleMap{"a": {}}},
			want: Target{Rule: RuleMap{"a": {"op": 1}}},
		},
		{
			src:  Target{Rule: RuleMap{"a": {}}},
			dst:  Target{Rule: RuleMap{"a": {"op": 1}}},
			want: Target{Rule: RuleMap{"a": {"op": 1}}},
		},
		{
			src:  Target{Rule: RuleMap{}},
			dst:  Target{Rule: RuleMap{"a": {"op": 1}}},
			want: Target{Rule: RuleMap{"a": {"op": 1}}},
		},
		{
			src:  Target{Rule: RuleMap{"a": {"op": 1}}},
			dst:  Target{Rule: RuleMap{}},
			want: Target{Rule: RuleMap{"a": {"op": 1}}},
		},
		{
			src:  Target{Rule: RuleMap{"a": {"op": false}}},
			dst:  Target{Rule: RuleMap{"a": {"op": true}}},
			want: Target{Rule: RuleMap{"a": {"op": true}}},
		},
		{
			src:  Target{Rule: RuleMap{"a": {"op": true}}},
			dst:  Target{Rule: RuleMap{"a": {"op": false}}},
			want: Target{Rule: RuleMap{"a": {"op": false}}},
		},
		{
			src:  Target{Rule: RuleMap{"a": {"op": true}}},
			dst:  Target{Rule: RuleMap{"b": {"op": false}}},
			want: Target{Rule: RuleMap{"a": {"op": true}, "b": {"op": false}}},
		},
	}

	for _, tt := range tests {
		got := tt.src.Merge(tt.dst)
		assert.Equal(t, tt.want, got)
	}
}
