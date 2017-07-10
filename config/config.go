package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/fatih/set"
	zglob "github.com/mattn/go-zglob"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/mohae/deepcopy"
	"github.com/synchro-food/filelint/lib"
)

type Config struct {
	File    File     `yaml:"files"`
	Targets []Target `yaml:"targets"`
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

	if err := yaml.Unmarshal(src, &userConfig); err != nil {
		return nil, err
	}

	conf.Merge(userConfig)

	return conf, nil
}

func NewDefaultConfig() (*Config, error) {
	conf := &Config{}
	src, err := configDefaultYmlBytes()
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(src, &conf); err != nil {
		return nil, err
	}

	return conf, nil
}

func (src *Config) Merge(dst *Config) {
	if len(dst.File.Include) > 0 {
		src.File.Include = dst.File.Include
	}
	src.File.Exclude = append(src.File.Exclude, dst.File.Exclude...)
	src.Targets = append(src.Targets, dst.Targets...)
}

func (cfg *Config) MatchedRule(file string) RuleMap {
	rm := make(RuleMap)

	for _, t := range cfg.Targets {
		if match(file, t.Patterns) {
			rm = rm.Merge(t.Rule)
		}
	}

	return rm
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
	f.Include = addGlobSignIfDir(f.Include...)
	f.Exclude = addGlobSignIfDir(f.Exclude...)

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

type Target struct {
	Patterns []string `yaml:"patterns"`
	Rule     RuleMap  `yaml:"rules"`
}

type RuleMap map[string]map[string]interface{}

func (rm RuleMap) Merge(dst RuleMap) (ret RuleMap) {
	ret = deepcopy.Copy(rm).(RuleMap)

	for ruleName, options := range dst {
		if _, ok := ret[ruleName]; ok {
			for optionKey, value := range options {
				ret[ruleName][optionKey] = value
			}
		} else {
			ret[ruleName] = options
		}
	}

	return ret
}

var (
	fileName   = ".filelint.yml"
	searchPath = "."
)

func SearchConfigFile() (f string, ok bool, err error) {
	if f := filepath.Join(searchPath, fileName); lib.IsExist(f) {
		return f, true, nil
	}

	if gitRoot, err := lib.FindGitRootPath(searchPath); err != nil {
		if err != lib.ErrNotGitRepository {
			return "", false, err
		}
	} else {
		if f := filepath.Join(gitRoot, fileName); lib.IsExist(f) {
			return f, true, nil
		}
	}

	if home, err := homedir.Dir(); err != nil {
		return "", false, err
	} else {
		if f := filepath.Join(home, fileName); lib.IsExist(f) {
			return f, true, nil
		}
	}

	return "", false, nil
}
