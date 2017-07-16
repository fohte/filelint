package dispatcher

import (
	"fmt"

	gitignore "github.com/sabhiram/go-gitignore"
	"github.com/synchro-food/filelint/config"
	"github.com/synchro-food/filelint/lint"
)

type Dispatcher struct {
	config *config.Config
}

func NewDispatcher(cfg *config.Config) *Dispatcher {
	return &Dispatcher{
		config: cfg,
	}
}

func (dp *Dispatcher) Dispatch(
	gitignorePath string,
	onDipatched func(file string, rules []lint.Rule) error,
) error {
	files, err := dp.config.File.FindTargets()
	if err != nil {
		return err
	}

	if gitignorePath != "" {
		files, err = excludeFilesWithGitIgnore(files, gitignorePath)
		if err != nil {
			return err
		}
	}

	for _, file := range files {
		definedRules := lint.GetDefinedRules()
		rules := make([]lint.Rule, 0, definedRules.Size())
		userRules := dp.config.MatchedRule(file)

		for ruleName, options := range userRules {
			if !definedRules.Has(ruleName) {
				return fmt.Errorf("%s is undefined", ruleName)
			}
			if options["enforce"] != true {
				continue
			}
			rule, err := definedRules.Get(ruleName).New(options)
			if err != nil {
				return err
			}
			rules = append(rules, rule)
		}

		if err := onDipatched(file, rules); err != nil {
			return err
		}
	}

	return nil
}

func excludeFilesWithGitIgnore(files []string, gitignorePath string) ([]string, error) {
	gi, err := gitignore.CompileIgnoreFile(gitignorePath)
	if err != nil {
		return nil, err
	}

	newFiles := make([]string, 0, len(files))

	for _, f := range files {
		if !gi.MatchesPath(f) {
			newFiles = append(newFiles, f)
		}
	}

	return newFiles, nil
}
