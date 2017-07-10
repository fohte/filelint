# Filelint

[![Build Status](https://travis-ci.org/synchro-food/filelint.svg?branch=master)](https://travis-ci.org/synchro-food/filelint)

Filelint is a CLI tool for linting any text file following some coding style.

![filelint-example](https://user-images.githubusercontent.com/11088009/27952943-16962632-6345-11e7-896f-f6d43aff084b.gif)

## Installation

You can download the binary from [the release page](https://github.com/synchro-food/filelint/releases) and place it in `$PATH` directories.

Or you can use `go get`:

```
$ go get -u github.com/synchro-food/filelint
```

## Usage

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

### Options

Filelint is available some flags:

```
$ filelint --help
Filelint is a CLI tool for linting any text file following some coding style.

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

## Configulation

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
  - patterns: ['**/*'] # `patterns` is an array of string, this pattern specify target files for this group
    rules:
      <rule-name>:
        enforce: true # or false
        <option-key>: <option-value>
        # ...
      # ...
  - patterns: # these rules apply to only .md and .mkd files
      - '**/*.md'
      - '**/*.mkd'
    rules:
      # ...
  # and other patterns and rules ...
```

The default configulation is [here](https://github.com/synchro-food/filelint/blob/master/config/default.yml).

## Rules

### `linebreak`

This rule enforces consistent linebreak style to Unix style (LF) or Windows style (CRLF).

- default: enforce

#### Options

##### `style`

This option specify the line endings (LF or CRLF).
LF is the Unix line endings (`\n`), and CRLF is the Windows line endings (`\r\n`).

- default: `lf`
- available values: `lf` or `crlf` (case insensitive)

### `first-newline`

This rule enforces some newlines at first of files.

- default: enforce

#### Options

##### `num`

This option specify the number of newlines at first of files.

- default: `0`
- available values: positive integers

### `final-newline`

This rule enforces some newlines at final of files.

- default: enforce

#### Options

##### `num`

This option specify the number of newlines at first of files.

- default: `1`
- available values: positive integers

### `no-eol-space`

This rule enforces no trailing whitespaces and tabs at the end of lines.

- default: enforce

#### Options

This rule has no options.

### `indent` (experimental)

**This rule is experimental. We recommend to use other linter tools that can analyze syntax of codes.**

This rule enforces a consistent indentation style to the any code and text.

- default: no enforce

##### `style`

This option specify the indentation style.

- default: `soft`
- available values: `soft`/`space` or `hard`/`tab`

##### `size`

This option specify number of whitespaces for the `soft` indentation style.

- default: `2`
- available values: positive numbers

### `no-bom`

This rule enforces no byte order marks (BOM) of UTF-8 to any text files.

- default: enforce

#### Options

This rule has no options.
