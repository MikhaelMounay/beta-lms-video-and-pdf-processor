package logger

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Logger struct {
	Buffer     io.Writer
	MsgCounter int
}

func NewLogger(buffer io.Writer) *Logger {
	return &Logger{
		Buffer:     buffer,
		MsgCounter: 0,
	}
}

func (l *Logger) Log(msg string) {
	l.MsgCounter++
	fmt.Fprintf(l.Buffer, "[%02d] %s", l.MsgCounter, msg)
}

func (l *Logger) Logf(format string, a ...any) {
	l.MsgCounter++
	s := fmt.Sprintf(format, a...)
	fmt.Fprintf(l.Buffer, "[%02d] %s", l.MsgCounter, s)
}

func (l *Logger) Prompt(msg string) {
	fmt.Fprintf(l.Buffer, "> %s", msg)
}

func (l *Logger) Promptf(format string, a ...any) {
	s := fmt.Sprintf(format, a...)
	fmt.Fprintf(l.Buffer, "> %s", s)
}

func (l *Logger) Fatal(format string, a ...any) {
	l.MsgCounter++
	fmt.Fprintf(l.Buffer, format, a...)
	fmt.Fprintf(l.Buffer, "\n\nFatal error happened! Press [Enter/Return] to exit... ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	os.Exit(1)
}
