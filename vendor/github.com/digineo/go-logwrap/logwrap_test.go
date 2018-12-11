package logwrap

import (
	"bytes"
	"fmt"
	"log"
	"path"
	"runtime"
	"testing"
)

// getShortfile returns (file name, line number) of the caller.
func getShortfile() (string, int) {
	_, f, l, ok := runtime.Caller(1)
	if !ok {
		panic("runtime.Caller should not fail")
	}
	return path.Base(f), l
}

func check(t *testing.T, buf *bytes.Buffer, expected string) {
	line := buf.String()
	buf.Reset()

	if actual := string(line); actual != expected {
		t.Errorf("expected log output to be %q, got %q", expected, actual)
	}
}

func TestInstance(t *testing.T) {
	// log.LstdFlags is a timestamp, and we don't want to mock time here
	// (which would make this exercise unnecessary complex)
	var buf bytes.Buffer
	l := log.New(&buf, "", log.Lshortfile)
	logger := &Instance{}
	logger.o = l.Output

	{
		logger.Infof("test %d", 1)
		fn, ln := getShortfile()
		check(t, &buf, fmt.Sprintf("%s:%d: INFO - test 1\n", fn, ln-1))
	}
	{
		logger.Errorf("test %d", 2)
		fn, ln := getShortfile()
		check(t, &buf, fmt.Sprintf("%s:%d: ERROR - test 2\n", fn, ln-1))
	}
}

type altLogger struct {
	buf bytes.Buffer
}

func (alt *altLogger) Infof(format string, a ...interface{}) {
	fmt.Fprintf(&alt.buf, format, a...)
	alt.buf.WriteRune('\n')
}

func (alt *altLogger) Errorf(format string, a ...interface{}) {
	fmt.Fprintf(&alt.buf, format, a...)
	alt.buf.WriteRune('\n')
}

func TestAlternative(t *testing.T) {
	alt := &altLogger{}
	logger := &Instance{}
	logger.SetLogger(alt)

	{
		logger.Infof("test alt %d", 1)
		check(t, &alt.buf, fmt.Sprintf("test alt 1\n"))
	}
	{
		logger.Errorf("test alt %d", 2)
		check(t, &alt.buf, fmt.Sprintf("test alt 2\n"))
	}
}
