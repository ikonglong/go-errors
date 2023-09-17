package errors

import (
	"fmt"
	"io"
)

func FormatError(err error, fmtS fmt.State, verb rune) {
	if IsNil(err) {
		return
	}

	s := &state{
		State:        fmtS,
		endingFnName: "",
	}

	switch verb {
	case 'v':
		if fmtS.Flag('+') {
			s.formatRecursive(err, true)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, err.Error())
	case 'q':
		fmt.Fprintf(s, "%q", err.Error())
	}
}

// state tracks error printing state.
type state struct {
	// state inherits fmt.State.
	fmt.State

	// endingFnName identifies the end frame of the error stack on an error level.
	endingFnName string
}

func (s *state) formatRecursive(err error, isOutermost bool) {
	// print the head of error info
	if isOutermost {
		fmt.Fprintf(s, "\nError occurred: %s", err.Error())
	} else {
		fmt.Fprintf(s, "\nCaused by: %s", err.Error())
	}

	switch err.(type) {
	case fmt.Formatter:
		if stp, ok := err.(StackTraceProvider); ok {
			st := stp.StackTrace()
			for _, frame := range st {
				fmt.Fprintf(s, "\n%+v", frame)
				if s.endingFnName != "" && s.endingFnName == frame.name() {
					break
				}
			}

			if len(st) > 0 {
				// set the ending func name which identifies the ending stack frame of the error
				// at the next error layer
				s.endingFnName = st[0].name()
			}
		} else {
			// print nothing, because the return value of .Error() was already
			// printed when the head of error info was printed at the beginning of
			// this function.
			//fmt.Fprintf(s, "\n%s", err.Error())
		}

	default:
		// print nothing, because the return value of .Error() was already
		// printed when the head of error info was printed at the beginning of
		// this function.
		//fmt.Fprintf(s, "\n%s", err.Error())
	}

	cause := UnwrapOnce(err)
	if cause != nil {
		s.formatRecursive(cause, false)
	}
}

type StackTraceProvider interface {
	StackTrace() StackTrace
}
