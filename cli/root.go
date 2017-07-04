package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"

	yaml "gopkg.in/yaml.v2"

	"github.com/fohte/filelint/config"
	"github.com/fohte/filelint/dispatcher"
	"github.com/fohte/filelint/lib"
	"github.com/fohte/filelint/lint"

	"github.com/spf13/cobra"
)

const Version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:           "filelint [files...]",
	Short:         "lint any text file following some file format",
	Long:          `Filelint is a CLI tool for linting any text file following some file format.`,
	RunE:          execute,
	SilenceUsage:  true,
	SilenceErrors: true,
}

var (
	showVersion      bool
	configFile       string
	useDefaultConfig bool
	printConfig      bool
	autofix          bool
	quiet            bool
	showTargets      bool
	useGitIgnore     bool
)

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
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
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
	if showVersion {
		fmt.Printf("filelint v%s [%s %s-%s]\n", Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		return nil
	}

	if configFile != "" && !lib.IsExist(configFile) {
		return Raise(ErrNoSuchConfigFile)
	}

	var cfg *config.Config
	if useDefaultConfig {
		var err error
		cfg, err = config.NewDefaultConfig()
		if err != nil {
			return Raise(err)
		}
	} else {
		if configFile == "" {
			configFile, _ = config.SearchConfigFile()
		}
		var err error
		cfg, err = config.NewConfig(configFile)
		if err != nil {
			return Raise(err)
		}
	}

	if len(args) > 0 {
		cfg.File.Include = args
	}

	if showTargets {
		fs, err := cfg.File.FindTargets()
		if err != nil {
			return Raise(err)
		}
		buf := bufio.NewWriter(os.Stdout)
		for _, f := range fs {
			fmt.Fprintln(buf, f)
		}
		buf.Flush()
		return nil
	}

	if printConfig {
		yml, err := yaml.Marshal(cfg)
		if err != nil {
			return Raise(err)
		}
		fmt.Printf("%s", yml)
		return nil
	}

	buf := bufio.NewWriter(os.Stdout)

	linterResult := struct {
		numErrors      int
		numFixedErrors int
		numErrorFiles  int
		numFixedFiles  int
	}{}

	dp := dispatcher.NewDispatcher(cfg)
	if err := dp.Dispatch(useGitIgnore, func(file string, rules []lint.Rule) error {
		linter, err := lint.NewLinter(file, rules)
		if err != nil {
			return err
		}

		result, err := linter.Lint()
		if err != nil {
			return err
		}

		if num := len(result.Reports); num > 0 {
			linterResult.numErrors += num
			linterResult.numErrorFiles++

			for _, report := range result.Reports {
				if autofix {
					fmt.Fprintf(buf, "[autofixed]")
					linterResult.numFixedErrors++
				}
				fmt.Fprintf(buf, "%s:%s\n", file, report.String())
			}

			if autofix {
				if err := writeFile(file, result.Fixed); err != nil {
					return err
				}
				linterResult.numFixedFiles++
			}
		}

		return nil
	}); err != nil {
		return Raise(err)
	}

	if !quiet {
		buf.Flush()
	}

	if !autofix && linterResult.numErrors > 0 {
		fmt.Printf("%d lint error(s) detected in %d file(s)\n", linterResult.numErrors, linterResult.numErrorFiles)
		return Raise(errLintFailed)
	}

	if linterResult.numFixedFiles > 0 && !quiet {
		fmt.Printf("%d lint error(s) autofixed in %d file(s)\n", linterResult.numFixedErrors, linterResult.numFixedFiles)
	}

	return nil
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

func init() {
	rootCmd.Flags().BoolVarP(&showVersion, "version", "v", false, "print the version and quit")
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "", "specify configuration file")
	rootCmd.Flags().BoolVarP(&printConfig, "print-config", "", false, "print the configuration")
	rootCmd.Flags().BoolVarP(&useDefaultConfig, "no-config", "", false, "don't use config file (use the application default config)")
	rootCmd.Flags().BoolVarP(&autofix, "fix", "", false, "automatically fix problems")
	rootCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "don't print lint errors or fixed files")
	rootCmd.Flags().BoolVarP(&showTargets, "print-targets", "", false, "print all lint target files and quit")
	rootCmd.Flags().BoolVarP(&useGitIgnore, "use-gitignore", "", true, "(experimental) read and use .gitignore file for excluding target files")
}
