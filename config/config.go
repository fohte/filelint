package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/fatih/set"
	"github.com/fohte/filelint/lib"
	zglob "github.com/mattn/go-zglob"
	"github.com/mohae/deepcopy"
)

type Config struct {
	File    File      `yaml:"files"`
	Targets TargetMap `yaml:"targets"`
}

func NewConfig(configFile string) (*Config, error) {
	conf, err := NewDefaultConfig()
	if err != nil {
		return nil, err
	}

	userConfig := &Config{}
	src, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	yaml.Unmarshal(src, &userConfig)

	conf.Merge(userConfig)
	conf.Targets.ExtendDefaultTarget()

	return conf, nil
}

func NewDefaultConfig() (*Config, error) {
	conf := &Config{}
	src, err := configDefaultYmlBytes()
	if err != nil {
		return nil, err
	}
	yaml.Unmarshal(src, &conf)
	conf.Targets.ExtendDefaultTarget()
	return conf, nil
}

func (src *Config) Merge(dst *Config) {
	if len(dst.File.Include) > 0 {
		src.File.Include = dst.File.Include
	}
	src.File.Exclude = append(src.File.Exclude, dst.File.Exclude...)

	for key, target := range dst.Targets {
		if _, ok := src.Targets[key]; ok {
			src.Targets[key] = src.Targets[key].Merge(target)
		} else {
			src.Targets[key] = target
		}
	}
}

func (cfg *Config) MatchedRule(file string) RuleMap {
	for k, t := range cfg.Targets {
		if k == "default" {
			continue
		}
		if match(file, t.Pattern) {
			return t.Rule
		}
	}

	return cfg.Targets["default"].Rule
}

func match(file string, patterns []string) bool {
	if len(patterns) == 0 {
		return true
	}

	for _, pattern := range patterns {
		ok, err := zglob.Match(pattern, file)
		if err != nil {
			return false
		}
		if ok {
			return true
		}
	}

	return false
}

type File struct {
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude"`
}

func (f File) FindTargets() ([]string, error) {
	f.addGlobSignIfDir()

	inSet, err := expandGlob(f.Include)
	if err != nil {
		return nil, err
	}

	exSet, _ := expandGlob(f.Exclude)
	if err != nil {
		return nil, err
	}

	fileSet := set.Difference(inSet, exSet)
	files := set.StringSlice(fileSet)
	files = lib.FindTextFiles(files)

	return files, nil
}

func (f File) addGlobSignIfDir() {
	f.Include = addGlobSignIfDir(f.Include...)
	f.Exclude = addGlobSignIfDir(f.Exclude...)
}

func expandGlob(files []string) (*set.Set, error) {
	fileSet := set.New()

	for _, f := range files {
		expanded, err := zglob.Glob(f)
		if err != nil {
			return nil, err
		}

		for _, ex := range expanded {
			fileSet.Add(ex)
		}
	}

	return fileSet, nil
}

func addGlobSignIfDir(files ...string) (dst []string) {
	dst = make([]string, 0, len(files))

	for _, f := range files {
		if !strings.Contains(f, "*") {
			info, err := os.Stat(f)
			if err != nil {
				continue
			}

			if info.IsDir() {
				f = filepath.Join(f, "**", "*")
			}
		}

		dst = append(dst, f)
	}

	return dst
}

type TargetMap map[string]Target

func (tm TargetMap) ExtendDefaultTarget() {
	if !tm.hasDefaultTarget() {
		return
	}

	for k, t := range tm {
		if k == "default" {
			continue
		}
		tm[k] = tm["default"].Merge(t)
	}
}

func (tm TargetMap) hasDefaultTarget() bool {
	for k := range tm {
		if k == "default" {
			return true
		}
	}
	return false
}

type Target struct {
	Pattern []string `yaml:"patterns"`
	Rule    RuleMap  `yaml:"rules"`
}

func (src Target) Merge(dst Target) (ret Target) {
	ret = deepcopy.Copy(src).(Target)

	if len(dst.Pattern) > 0 {
		ret.Pattern = dst.Pattern
	}

	for ruleName, options := range dst.Rule {
		if _, ok := ret.Rule[ruleName]; ok {
			for optionKey, value := range options {
				ret.Rule[ruleName][optionKey] = value
			}
		} else {
			ret.Rule[ruleName] = options
		}
	}

	return ret
}

type RuleMap map[string]map[string]interface{}

var (
	fileName   = ".filelint.yml"
	searchPath = "."
)

func SearchConfigFile() (f string, ok bool) {
	if f := filepath.Join(searchPath, fileName); lib.IsExist(f) {
		return f, true
	}

	if gitRoot, err := lib.FindGitRootPath(searchPath); err != nil {
		if err != nil && err != lib.ErrNotGitRepository {
			return "", false
		}
		if f := filepath.Join(gitRoot, fileName); lib.IsExist(f) {
			return f, true
		}
	}

	if f := filepath.Join(lib.GetHomeDir(), fileName); lib.IsExist(f) {
		return f, true
	}

	return "", false
}
