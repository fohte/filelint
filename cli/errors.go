package cli

import "errors"

var (
	errLintFailed = errors.New("lint errors detected")
)

const (
	DefaultExitStatus    = 2
	LintFailedExitStatus = 1
)

type ExitError interface {
	error
	ExitStatus() int
}

type exitError struct {
	exitStatus int
	message    string
}

func (ee *exitError) Error() string {
	return ee.message
}

func (ee *exitError) ExitStatus() int {
	return ee.exitStatus
}

func Raise(err error) ExitError {
	var exitStatus int

	switch err {
	case errLintFailed:
		exitStatus = LintFailedExitStatus
	default:
		exitStatus = DefaultExitStatus
	}

	return &exitError{
		exitStatus: exitStatus,
		message:    err.Error(),
	}
}
