package pathlib

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/fail"
)

func Exists(pathname string) bool {
	_, err := os.Stat(pathname)
	return !os.IsNotExist(err)
}

func Abs(path string) (string, error) {
	fullpath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return filepath.Clean(fullpath), nil
}

func Symlink(pathname string) (string, bool) {
	stat, err := os.Lstat(pathname)
	if err != nil {
		return "", false
	}
	mode := stat.Mode()
	if mode&fs.ModeSymlink == 0 {
		return "", false
	}
	name, err := os.Readlink(pathname)
	if err != nil {
		return "", false
	}
	return name, true
}

func IsDir(pathname string) bool {
	stat, err := os.Stat(pathname)
	return err == nil && stat.IsDir()
}

func IsEmptyDir(pathname string) bool {
	if !IsDir(pathname) {
		return false
	}
	content, err := os.ReadDir(pathname)
	if err != nil {
		return false
	}
	return len(content) == 0
}

func IsFile(pathname string) bool {
	stat, err := os.Stat(pathname)
	return err == nil && !stat.IsDir()
}

func DaysSinceModified(filename string) (int, error) {
	stat, err := os.Stat(filename)
	if err != nil {
		return -1, err
	}
	return common.DayCountSince(stat.ModTime()), nil
}

func Size(pathname string) (int64, bool) {
	stat, err := os.Stat(pathname)
	if err != nil {
		return 0, false
	}
	return stat.Size(), true
}

func Modtime(pathname string) (time.Time, error) {
	stat, err := os.Stat(pathname)
	if err != nil {
		return time.Now(), err
	}
	return stat.ModTime(), nil
}

func hasCorrectMode(stat fs.FileInfo, expected fs.FileMode) bool {
	return expected == (stat.Mode() & expected)
}

func ensureCorrectMode(fullpath string, stat fs.FileInfo, correct fs.FileMode) (string, error) {
	if hasCorrectMode(stat, correct) {
		return fullpath, nil
	}
	err := os.Chmod(fullpath, correct)
	if err != nil {
		return "", err
	}
	return fullpath, nil
}

func makeModedDir(fullpath string, correct fs.FileMode) (path string, err error) {
	defer fail.Around(&err)

	stat, err := os.Stat(fullpath)
	if err == nil && stat.IsDir() {
		return ensureCorrectMode(fullpath, stat, correct)
	}
	fail.On(err == nil, "Path %q exists, but is not a directory!", fullpath)
	_, err = shared.MakeSharedDir(filepath.Dir(fullpath))
	fail.On(err != nil, "%v", err)
	err = os.Mkdir(fullpath, correct)
	fail.On(err != nil, "Failed to create directory %q, reason: %v", fullpath, err)
	stat, err = os.Stat(fullpath)
	fail.On(err != nil, "Failed to stat created directory %q, reason: %v", fullpath, err)
	_, err = ensureCorrectMode(fullpath, stat, correct)
	fail.On(err != nil, "Failed to make created directory shared %q, reason: %v", fullpath, err)
	return fullpath, nil
}

func MakeSharedFile(fullpath string) (string, error) {
	return shared.MakeSharedFile(fullpath)
}

func MakeSharedDir(fullpath string) (string, error) {
	return shared.MakeSharedDir(fullpath)
}

func ForceSharedDir(fullpath string) (string, error) {
	return makeModedDir(fullpath, 0777)
}

func IsSharedDir(fullpath string) bool {
	stat, err := os.Stat(fullpath)
	if err != nil {
		return false
	}
	return stat.IsDir() && hasCorrectMode(stat, 0777)
}

func doEnsureDirectory(directory string, mode fs.FileMode) (string, error) {
	fullpath, err := filepath.Abs(directory)
	if err != nil {
		return "", err
	}
	if IsDir(fullpath) {
		return fullpath, nil
	}
	err = os.MkdirAll(fullpath, mode)
	if err != nil {
		return "", err
	}
	stats, err := os.Stat(fullpath)
	if !stats.IsDir() {
		return "", fmt.Errorf("Path %s is not a directory!", fullpath)
	}
	return fullpath, nil
}

func EnsureSharedDirectory(directory string) (string, error) {
	return shared.MakeSharedDir(directory)
}

func EnsureSharedParentDirectory(resource string) (string, error) {
	return EnsureSharedDirectory(filepath.Dir(resource))
}

func EnsureDirectory(directory string) (string, error) {
	return doEnsureDirectory(directory, 0o750)
}

func RemoveEmptyDirectores(starting string) (err error) {
	defer fail.Around(&err)

	return DirWalk(starting, func(fullpath, relative string, entry os.FileInfo) {
		if IsEmptyDir(fullpath) {
			err = os.Remove(fullpath)
			fail.On(err != nil, "%s", err)
		}
	})
}
