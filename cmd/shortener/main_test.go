package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"testing"
)

func SplitTestArgs(args []string) (testArgs, appArgs []string) {
	for i, arg := range args {
		switch {
		case strings.HasPrefix(arg, "-test."):
			testArgs = append(testArgs, arg)
		case i == 0:
			appArgs = append(appArgs, arg)
			testArgs = append(testArgs, arg)
		default:
			appArgs = append(appArgs, arg)
		}
	}
	return
}

func TestEmpty(t *testing.T) {}

func TestMain(m *testing.M) {
	if strings.HasSuffix(os.Args[0], ".test") {
		log.Printf("skip launching server when invoked via go test")
		return
	}

	testArgs, _ := SplitTestArgs(os.Args)
	notifier := make(chan os.Signal, 1)
	signal.Notify(notifier, syscall.SIGINT, syscall.SIGTERM)

	// This will generate coverage files:
	os.Args = testArgs
	m.Run()
}