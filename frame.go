package runtimex

import (
	"fmt"
	"path"
	"runtime"
)

// FrameFormatterFn defines the function signature for manipulating parts of the frame before printing it.
type FrameFormatterFn func(string) string

// Frame represents a single frame on the stack.
type Frame struct {
	runtime.Frame

	// formatFileFn is called to format the filename before it is printed via Format().
	formatFileFn FrameFormatterFn

	// formatFuncFn is called to format the function name before it is printed via Format().
	formatFuncFn FrameFormatterFn
}

// FrameFromPC returns a frame from the given PC.
func FrameFromPC(pc uintptr) Frame {
	frames := runtime.CallersFrames([]uintptr{pc})
	frame, _ := frames.Next()
	return Frame{
		Frame: frame,
	}
}

// Caller returns the frame from the location at which this function was called.
func Caller(skip int) Frame {
	pc := make([]uintptr, 1)
	runtime.Callers(skip+2, pc)
	frames := runtime.CallersFrames(pc)
	frame, _ := frames.Next()
	return Frame{
		Frame: frame,
	}
}

// WithFormatFileFn returns a new frame using the given function to format the filename.
func (f Frame) WithFormatFileFn(fn FrameFormatterFn) Frame {
	return Frame{
		Frame:        f.Frame,
		formatFileFn: fn,
		formatFuncFn: f.formatFuncFn,
	}
}

// WithFormatFuncFn returns a new frame using the given function to format the function name.
func (f Frame) WithFormatFuncFn(fn FrameFormatterFn) Frame {
	return Frame{
		Frame:        f.Frame,
		formatFileFn: f.formatFileFn,
		formatFuncFn: fn,
	}
}

// Format is used to format a frame using Printf()-like functions.
//
// The following formats are supported:
//
//	%s - Prints filename:line (where filename is only the file name)
//	%+s - Prints filepath:line (where filepath is the full path to the file)
//	%v - Prints function followed by a newline and tab then filepath:line
//	%+v - Prints function (address) followed by a newline and tab then filepath:line
func (f Frame) Format(state fmt.State, verb rune) {
	if f.formatFileFn != nil {
		f.Frame.File = f.formatFileFn(f.Frame.File)
	}
	if f.formatFuncFn != nil {
		f.Frame.Function = f.formatFuncFn(f.Frame.Function)
	}

	switch verb {
	case 's':
		switch {
		case state.Flag('+'):
			fmt.Fprintf(state, "%s:%d", f.Frame.File, f.Frame.Line)
		default:
			fmt.Fprintf(state, "%s:%d", path.Base(f.Frame.File), f.Frame.Line)
		}
	case 'v':
		switch {
		case state.Flag('+'):
			fmt.Fprintf(state, "%s (%p)\n\t%s:%d", f.Function, f.Func, f.File, f.Line)
		default:
			fmt.Fprintf(state, "%s\n\t%s:%d", f.Function, f.File, f.Line)
		}
	default:
		fmt.Fprintf(state, "%p", &f.Frame)
	}
}
