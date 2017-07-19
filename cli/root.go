package cli

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"

	yaml "gopkg.in/yaml.v2"

	"github.com/synchro-food/filelint/config"
	"github.com/synchro-food/filelint/dispatcher"
	"github.com/synchro-food/filelint/lib"
	"github.com/synchro-food/filelint/lint"

	"github.com/spf13/cobra"
)

const Version = "0.2.0"

var rootCmd = &cobra.Command{
	Use:           "filelint [files...]",
	Short:         "lint any text file following some coding style",
	Long:          `Filelint is a CLI tool for linting any text file following some coding style.`,
	RunE:          execute,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var (
	configFile       string
	userRules        []string
	isShowVersion    bool
	isPrintConfig    bool
	isPrintTarget    bool
	isAutofix        bool
	isQuiet          bool
	useDefaultConfig bool
	useGitIgnore     bool
)

func init() {
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "specify configuration file")
	rootCmd.Flags().StringArrayVar(&userRules, "rule", []string{}, "specify rules")
	rootCmd.Flags().BoolVarP(&isShowVersion, "version", "v", false, "print the version and quit")
	rootCmd.Flags().BoolVar(&isPrintConfig, "print-config", false, "print the configuration")
	rootCmd.Flags().BoolVar(&isPrintTarget, "print-targets", false, "print all lint target files and quit")
	rootCmd.Flags().BoolVar(&isAutofix, "fix", false, "automatically fix problems")
	rootCmd.Flags().BoolVarP(&isQuiet, "quiet", "q", false, "don't print lint errors or fixed files")
	rootCmd.Flags().BoolVar(&useDefaultConfig, "no-config", false, "don't use config file (use the application default config)")
	rootCmd.Flags().BoolVar(&useGitIgnore, "use-gitignore", true, "(experimental) read and use .gitignore file for excluding target files")
}

var (
	ErrNoSuchConfigFile = errors.New("no such config file")
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		exitStatus := DefaultExitStatus

		if ee, ok := err.(ExitError); ok {
			exitStatus = ee.ExitStatus()
		}

		switch exitStatus {
		case LintFailedExitStatus:
			break
		case DefaultExitStatus:
			fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
			rootCmd.Usage()
		default:
			panic(err.Error())
		}

		os.Exit(exitStatus)
	}
}

func execute(cmd *cobra.Command, args []string) error {
	var out io.Writer
	if isQuiet {
		out = ioutil.Discard
	} else {
		out = os.Stdout
	}

	if isShowVersion {
		showVersion()
		return nil
	}

	cfg, err := loadConfig(configFile, useDefaultConfig)
	if err != nil {
		return Raise(err)
	}

	if len(userRules) > 0 {
		userRuleMap := make(config.RuleMap)
		for _, r := range userRules {
			if err := yaml.Unmarshal([]byte(r), &userRuleMap); err != nil {
				return err
			}
		}
		cfg.Targets = append(cfg.Targets, config.Target{
			Patterns: []string{"**/*"},
			Rule:     userRuleMap,
		})
	}

	if len(args) > 0 {
		cfg.File.Include = args
	}

	if isPrintConfig {
		if err := printConfig(out, cfg); err != nil {
			return Raise(err)
		}
		return nil
	}

	if isPrintTarget {
		if err := printTarget(out, cfg.File); err != nil {
			return Raise(err)
		}
		return nil
	}

	var gitignorePath string
	if useGitIgnore {
		var err error
		gitignorePath, err = lib.FindGitIgnore()
		if err != nil {
			return Raise(err)
		}
	}

	if err := runLint(out, isAutofix, cfg, gitignorePath); err != nil {
		return Raise(err)
	}

	return nil
}

func showVersion() {
	fmt.Printf("filelint v%s [%s %s-%s]\n", Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

func printConfig(out io.Writer, cfg *config.Config) error {
	yml, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "%s", yml)
	return nil
}

func printTarget(out io.Writer, f config.File) error {
	fs, err := f.FindTargets()
	if err != nil {
		return err
	}
	for _, f := range fs {
		fmt.Fprintln(out, f)
	}
	return nil
}

func runLint(out io.Writer, isAutofix bool, cfg *config.Config, gitignorePath string) error {
	dp := dispatcher.NewDispatcher(cfg)

	var (
		numErrors      int
		numFixedErrors int
		numErrorFiles  int
		numFixedFiles  int
	)

	if err := dp.Dispatch(gitignorePath, func(file string, rules []lint.Rule) error {
		linter, err := lint.NewLinter(file, rules)
		if err != nil {
			return err
		}

		result, err := linter.Lint()
		if err != nil {
			return err
		}

		if num := len(result.Reports); num > 0 {
			numErrors += num
			numErrorFiles++

			for _, report := range result.Reports {
				if isAutofix {
					fmt.Fprintf(out, "[autofixed]")
					numFixedErrors++
				}
				fmt.Fprintf(out, "%s:%s\n", file, report.String())
			}

			if isAutofix {
				if err := writeFile(file, result.Fixed); err != nil {
					return err
				}
				numFixedFiles++
			}
		}

		return nil
	}); err != nil {
		return err
	}

	if !isAutofix && numErrors > 0 {
		fmt.Fprintf(out, "%d lint error(s) detected in %d file(s)\n", numErrors, numErrorFiles)
		return errLintFailed
	}

	if numFixedFiles > 0 {
		fmt.Fprintf(out, "%d lint error(s) autofixed in %d file(s)\n", numFixedErrors, numFixedFiles)
		return nil
	}

	return nil
}

func loadConfig(configFile string, useDefault bool) (*config.Config, error) {
	if useDefault {
		cfg, err := config.NewDefaultConfig()
		if err != nil {
			return nil, err
		}
		return cfg, err
	}

	if configFile != "" && !lib.IsExist(configFile) {
		return nil, ErrNoSuchConfigFile
	}

	if configFile == "" {
		var exist bool
		var err error
		configFile, exist, err = config.SearchConfigFile()
		if err != nil {
			return nil, err
		}
		if !exist {
			return loadConfig("", true)
		}
	}

	cfg, err := config.NewConfig(configFile)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func writeFile(filename string, src []byte) error {
	var fp *os.File
	var err error

	if lib.IsExist(filename) {
		fp, err = os.Open(filename)
	} else {
		fp, err = os.Create(filename)
	}
	if err != nil {
		return err
	}
	defer fp.Close()

	fi, err := fp.Stat()
	if err != nil {
		return err
	}
	perm := fi.Mode().Perm()

	err = ioutil.WriteFile(filename, src, perm)
	if err != nil {
		return err
	}

	return nil
}
