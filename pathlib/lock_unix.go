//go:build darwin || linux || !windows
// +build darwin linux !windows

package pathlib

import (
	"os"
	"syscall"
	"time"

	"github.com/robocorp/rcc/common"
)

func Locker(filename string, trycount int) (Releaser, error) {
	if Lockless {
		return Fake(), nil
	}
	if common.TraceFlag {
		defer common.Stopwatch("LOCKER: Got lock on %v in", filename).Report()
	}
	common.Trace("LOCKER: Want lock on: %v", filename)
	_, err := EnsureSharedParentDirectory(filename)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	if err != nil {
		return nil, err
	}
	_, err = shared.MakeSharedFile(filename)
	if err != nil {
		return nil, err
	}
	err = syscall.Flock(int(file.Fd()), int(syscall.LOCK_EX))
	if err != nil {
		return nil, err
	}
	marker := lockPidFilename(filename)
	_, err = file.Write([]byte(marker))
	if err != nil {
		return nil, err
	}
	common.Debug("LOCKER: make marker %v", marker)
	ForceTouchWhen(marker, time.Now())
	return &Locked{file, marker}, nil
}

func (it Locked) Release() error {
	defer os.Remove(it.Marker)
	defer common.Debug("LOCKER: remove marker %v", it.Marker)
	defer it.Close()
	err := syscall.Flock(int(it.Fd()), int(syscall.LOCK_UN))
	common.Trace("LOCKER: release %v with err: %v", it.Name(), err)
	return err
}
