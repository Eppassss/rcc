package conda

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows/registry"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/pretty"
	"github.com/robocorp/rcc/settings"
	"github.com/robocorp/rcc/shell"
)

const (
	mingwSuffix    = "\\mingw-w64"
	Newline        = "\r\n"
	librarySuffix  = "\\Library"
	scriptSuffix   = "\\Scripts"
	usrSuffix      = "\\usr"
	binSuffix      = "\\bin"
	activateScript = "@echo off\n" +
		"set \"MAMBA_ROOT_PREFIX={{.MambaRootPrefix}}\"\n" +
		"for /f \"tokens=* usebackq\" %%a in ( `call \"{{.Robocorphome}}\\bin\\micromamba.exe\" shell -s cmd.exe activate -p \"{{.Live}}\"` ) do ( call \"%%a\" )\n" +
		"call \"{{.Rcc}}\" internal env -l after\n"
	commandSuffix = ".cmd"
)

func MicromambaLink() string {
	return settings.Global.DownloadsLink(micromambaLink("windows64", "micromamba.exe"))
}

var (
	Shell          = []string{"cmd.exe", "/K"}
	FileExtensions = []string{".exe", ".com", ".bat", ".cmd", ""}
)

func CondaEnvironment() []string {
	env := os.Environ()
	env = append(env, fmt.Sprintf("MAMBA_ROOT_PREFIX=%s", common.MambaRootPrefix()))
	tempFolder := common.RobocorpTemp()
	env = append(env, fmt.Sprintf("TEMP=%s", tempFolder))
	env = append(env, fmt.Sprintf("TMP=%s", tempFolder))
	return env
}

func BinMicromamba() string {
	return common.ExpandPath(filepath.Join(common.BinLocation(), "micromamba.exe"))
}

func CondaPaths(prefix string) []string {
	return []string{
		prefix,
		prefix + librarySuffix + mingwSuffix + binSuffix,
		prefix + librarySuffix + usrSuffix + binSuffix,
		prefix + librarySuffix + binSuffix,
		prefix + scriptSuffix,
		prefix + binSuffix,
	}
}

func IsWindows() bool {
	return true
}

func HasLongPathSupport() bool {
	baseline := []string{common.RobocorpHome(), "stump"}
	stumpath := filepath.Join(baseline...)
	defer os.RemoveAll(stumpath)

	for count := 0; count < 24; count++ {
		baseline = append(baseline, fmt.Sprintf("verylongpath%d", count+1))
	}
	fullpath := filepath.Join(baseline...)

	code, err := shell.New(nil, ".", "cmd.exe", "/c", "mkdir", fullpath).Transparent()
	common.Trace("Checking long path support with MKDIR '%v' (%d characters) -> %v [%v] {%d}", fullpath, len(fullpath), err == nil, err, code)
	if err != nil {
		longPathSupportArticle := settings.Global.DocsLink("product-manuals/robocorp-lab/troubleshooting#windows-has-to-have-long-filenames-support-on")
		common.Log("%sWARNING!  Long path support failed. Reason: %v.%s", pretty.Red, err, pretty.Reset)
		common.Log("%sWARNING!  See %v for more details.%s", pretty.Red, longPathSupportArticle, pretty.Reset)
		return false
	}
	return true
}

func EnforceLongpathSupport() error {
	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, `SYSTEM\CurrentControlSet\Control\FileSystem`, registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer key.Close()
	return key.SetDWordValue("LongPathsEnabled", 1)
}
