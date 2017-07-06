# Filelint

[![Build Status](https://travis-ci.org/synchro-food/filelint.svg?branch=master)](https://travis-ci.org/synchro-food/filelint)

Filelint is a CLI tool for linting any text file following some file format.

## Installation

You can download the binary from [the release page](https://github.com/synchro-food/filelint/releases) and place it in `$PATH` directories.

Or you can use `go get`:

```
$ go get -u github.com/synchro-food/filelint
```

## Usage

### CLI

Filelint is available some flags:

```
$ filelint --help
Filelint is a CLI tool for linting any text file following some file format.

Usage:
  filelint [files...] [flags]

Flags:
  -c, --config string   specify configuration file
      --fix             automatically fix problems
  -h, --help            help for filelint
      --no-config       don't use config file (use the application default config)
      --print-config    print the configuration
      --print-targets   print all lint target files and quit
  -q, --quiet           don't print lint errors or fixed files
      --use-gitignore   (experimental) read and use .gitignore file for excluding target files (default true)
  -v, --version         print the version and quit
```

The `files` optional argument is linting target files.
If not pass `files` then all text files in current directory recursively.

#### Example

You can run Filelint on all text files in current directory recursively:
```
$ filelint
```

Or you can specify any directory:
```
$ filelint scripts/
```

Or you can specify any files and fix them:
```
$ filelint README.md --fix
```

Or you can specify any path and fix them:
```
$ filelint README.md scripts/ --fix
```

### Configulation

Filelint can configure lint rule settings and format target files via `.filelint.yml`.  
`.filelint.yml` is searched in current directory, repo root directory if you use git, or `$HOME`.

The `.filelint.yml` can use following style:

```yaml
files:
  include:
    - 'inculde/path/to/**/files'
  exclude:
    - 'exclude/path/to'
targets:
  # `default` group is applies all files and other groups extend from this group
  default:
    rules:
      <rule-name>:
        enforce: true # or false
        <option-key>: <option-value>
        # ...
      # ...
  # `group-name` can use anything except `default` available
  <group-name>:
    # `patterns` is the specific target files for this group
    patterns:
      - '**/*.md'
    rules:
      # ...
  # ...
```

The default configulation is [here](https://github.com/synchro-food/filelint/blob/master/config/default.yml).
