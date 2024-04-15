package runtimex

import (
	"fmt"
	"io"
	"runtime"
)

// stack is stack of Frames from innermost (newest) to outermost (oldest).
type stack struct {
	Frames []Frame
}

// Stack returns the call stack at the position at which the function was called.
//
// You can optionally specify up to 2 format functions - the first being the formatter to call for the frame's
// file and the second being the formatter to call for the frame's function. Making either function nil will
// cause it to be ignored.
func Stack(skipFrames int, formatFns ...FrameFormatterFn) stack {
	stack := stack{
		Frames: []Frame{},
	}

	// get callers
	bufSize := 32
	pc := make([]uintptr, bufSize)
	count := runtime.Callers(skipFrames+2, pc)
	for count == bufSize {
		bufSize *= 2
		pc = make([]uintptr, bufSize)
		count = runtime.Callers(skipFrames+2, pc)
	}

	// turn callers into frames
	frames := runtime.CallersFrames(pc)
	more := true
	var frame runtime.Frame
	for more {
		frame, more = frames.Next()
		f := Frame{Frame: frame}
		numFn := len(formatFns)
		if numFn >= 2 {
			f = f.WithFormatFileFn(formatFns[0]).WithFormatFuncFn(formatFns[1])
		} else if numFn == 1 {
			f = f.WithFormatFileFn(formatFns[0])
		}
		stack.Frames = append(stack.Frames, f)
	}
	return stack
}

// Format formats the stack of Frames according to the fmt.Formatter interface.
//
//	%s  lists source files for each Frame in the stack
//	%v  lists the source file and line number for each Frame in the stack
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//	%+v  Prints filename, function, and line number for each Frame in the stack.
func (s stack) Format(state fmt.State, verb rune) {
	for _, f := range s.Frames {
		io.WriteString(state, "\n")
		f.Format(state, verb)
	}
}
