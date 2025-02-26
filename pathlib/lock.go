package pathlib

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/robocorp/rcc/common"
)

var (
	slashPattern = regexp.MustCompile("[/\\\\]+")
)

type Releaser interface {
	Release() error
}

type Locked struct {
	*os.File
	Marker string
}

type fake bool

func (it fake) Release() error {
	return common.Trace("LOCKER: lockless mode release.")
}

func Fake() Releaser {
	common.Trace("LOCKER: lockless mode.")
	return fake(true)
}

func waitingLockNotification(message string, latch chan bool) {
	delay := 5 * time.Second
	counter := 0
	for {
		select {
		case <-latch:
			return
		case <-time.After(delay):
			counter += 1
			delay *= 3
			common.Log("#%d: %s (rcc lock wait warning)", counter, message)
			common.Timeline("waiting for lock")
		}
	}
}

func LockWaitMessage(message string) func() {
	latch := make(chan bool)
	go waitingLockNotification(message, latch)
	return func() {
		latch <- true
	}
}

func unslash(text string) string {
	parts := slashPattern.Split(text, -1)
	return strings.Join(parts, "_")
}

func lockPidFilename(lockfile string) string {
	now := time.Now().Format("20060102150405")
	base := filepath.Base(lockfile)
	username := "unspecified"
	who, err := user.Current()
	if err == nil {
		username = unslash(who.Username)
	}
	marker := fmt.Sprintf("%s_%s_%s_%s_%d_%s", now, username, common.ControllerType, common.HolotreeSpace, os.Getpid(), base)
	return filepath.Join(common.HololibPids(), marker)
}
