package runtimex_test

import (
	"fmt"
	"strings"
	"testing"

	"go.joshhogle.dev/runtimex"
)

// TODO: implement testing and benchmarks

func TestStack(t *testing.T) {
	test1(t)
}

func test1(t *testing.T) {
	test2(t)
}

func test2(t *testing.T) {
	test3(t)
}

func test3(t *testing.T) {
	s := runtimex.Stack(0, formatFrameFile)

	fmt.Printf("%s\n", s)
	fmt.Printf("%+s\n", s)
	fmt.Printf("%v\n", s)
	fmt.Printf("%+v\n", s)
	fmt.Printf("%d\n", s)
}

func formatFrameFile(file string) string {
	return strings.TrimPrefix(file, "/Users/joshhogle/workspace/src/")
}
