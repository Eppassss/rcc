package common

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	ROBOCORP_HOME_VARIABLE                = `ROBOCORP_HOME`
	VERBOSE_ENVIRONMENT_BUILDING          = `RCC_VERBOSE_ENVIRONMENT_BUILDING`
	ROBOCORP_OVERRIDE_SYSTEM_REQUIREMENTS = `ROBOCORP_OVERRIDE_SYSTEM_REQUIREMENTS`
)

var (
	NoBuild            bool
	Silent             bool
	DebugFlag          bool
	TraceFlag          bool
	DeveloperFlag      bool
	StrictFlag         bool
	SharedHolotree     bool
	LogLinenumbers     bool
	NoCache            bool
	NoOutputCapture    bool
	Liveonly           bool
	UnmanagedSpace     bool
	StageFolder        string
	ControllerType     string
	HolotreeSpace      string
	EnvironmentHash    string
	SemanticTag        string
	ForcedRobocorpHome string
	When               int64
	ProgressMark       time.Time
	Clock              *stopwatch
	randomIdentifier   string
)

func init() {
	Clock = &stopwatch{"Clock", time.Now()}
	When = Clock.When()
	ProgressMark = time.Now()

	randomIdentifier = fmt.Sprintf("%016x", rand.Uint64()^uint64(os.Getpid()))

	// Note: HololibCatalogLocation, HololibLibraryLocation and HololibUsageLocation
	//       are force created from "htfs" direcotry.go init function
	// Also: HolotreeLocation creation is left for actual holotree commands
	//       to prevent accidental access right problem during usage

	SharedHolotree = isFile(HoloInitUserFile())

	ensureDirectory(JournalLocation())
	ensureDirectory(TemplateLocation())
	ensureDirectory(BinLocation())
	ensureDirectory(PipCache())
	ensureDirectory(WheelCache())
	ensureDirectory(RobotCache())
	ensureDirectory(MambaPackages())
}

func RobocorpHome() string {
	if len(ForcedRobocorpHome) > 0 {
		return ExpandPath(ForcedRobocorpHome)
	}
	home := os.Getenv(ROBOCORP_HOME_VARIABLE)
	if len(home) > 0 {
		return ExpandPath(home)
	}
	return ExpandPath(defaultRobocorpLocation)
}

func RobocorpLock() string {
	return filepath.Join(RobocorpHome(), "robocorp.lck")
}

func VerboseEnvironmentBuilding() bool {
	return DebugFlag || TraceFlag || len(os.Getenv(VERBOSE_ENVIRONMENT_BUILDING)) > 0
}

func OverrideSystemRequirements() bool {
	return len(os.Getenv(ROBOCORP_OVERRIDE_SYSTEM_REQUIREMENTS)) > 0
}

func BinRcc() string {
	self, err := os.Executable()
	if err != nil {
		return os.Args[0]
	}
	return self
}

func OldEventJournal() string {
	return filepath.Join(RobocorpHome(), "event.log")
}

func EventJournal() string {
	return filepath.Join(JournalLocation(), "event.log")
}

func JournalLocation() string {
	return filepath.Join(RobocorpHome(), "journals")
}

func TemplateLocation() string {
	return filepath.Join(RobocorpHome(), "templates")
}

func RobocorpTempRoot() string {
	return filepath.Join(RobocorpHome(), "temp")
}

func RobocorpTempName() string {
	return filepath.Join(RobocorpTempRoot(), randomIdentifier)
}

func RobocorpTemp() string {
	tempLocation := RobocorpTempName()
	fullpath, err := filepath.Abs(tempLocation)
	if err != nil {
		fullpath = tempLocation
	}
	ensureDirectory(fullpath)
	if err != nil {
		Log("WARNING (%v) -> %v", tempLocation, err)
	}
	return fullpath
}

func BinLocation() string {
	return filepath.Join(RobocorpHome(), "bin")
}

func SharedMarkerLocation() string {
	return filepath.Join(HoloLocation(), "shared.yes")
}

func HoloLocation() string {
	return ExpandPath(defaultHoloLocation)
}

func HoloInitLocation() string {
	return filepath.Join(HoloLocation(), "lib", "catalog", "init")
}

func HoloInitUserFile() string {
	return filepath.Join(HoloInitLocation(), UserHomeIdentity())
}

func HoloInitCommonFile() string {
	return filepath.Join(HoloInitLocation(), "commons.tof")
}

func HolotreeLocation() string {
	if SharedHolotree {
		return HoloLocation()
	}
	return filepath.Join(RobocorpHome(), "holotree")
}

func HololibLocation() string {
	if SharedHolotree {
		return filepath.Join(HoloLocation(), "lib")
	}
	return filepath.Join(RobocorpHome(), "hololib")
}

func HololibPids() string {
	return filepath.Join(HololibLocation(), "pids")
}

func HololibCatalogLocation() string {
	return filepath.Join(HololibLocation(), "catalog")
}

func HololibLibraryLocation() string {
	return filepath.Join(HololibLocation(), "library")
}

func HololibUsageLocation() string {
	return filepath.Join(HololibLocation(), "used")
}

func HolotreeLock() string {
	return filepath.Join(HolotreeLocation(), "global.lck")
}

func UsesHolotree() bool {
	return len(HolotreeSpace) > 0
}

func PipCache() string {
	return filepath.Join(RobocorpHome(), "pipcache")
}

func WheelCache() string {
	return filepath.Join(RobocorpHome(), "wheels")
}

func RobotCache() string {
	return filepath.Join(RobocorpHome(), "robots")
}

func MambaRootPrefix() string {
	return RobocorpHome()
}

func MambaPackages() string {
	return ExpandPath(filepath.Join(MambaRootPrefix(), "pkgs"))
}

func PipRcFile() string {
	return ExpandPath(filepath.Join(RobocorpHome(), "piprc"))
}

func MicroMambaRcFile() string {
	return ExpandPath(filepath.Join(RobocorpHome(), "micromambarc"))
}

func SettingsFile() string {
	return ExpandPath(filepath.Join(RobocorpHome(), "settings.yaml"))
}

func CaBundleFile() string {
	return ExpandPath(filepath.Join(RobocorpHome(), "ca-bundle.pem"))
}

func UnifyVerbosityFlags() {
	if Silent {
		DebugFlag = false
		TraceFlag = false
	}
	if TraceFlag {
		DebugFlag = true
	}
}

func UnifyStageHandling() {
	if len(StageFolder) > 0 {
		Liveonly = true
	}
}

func ForceDebug() {
	Silent = false
	DebugFlag = true
	UnifyVerbosityFlags()
}

func Platform() string {
	return strings.ToLower(fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH))
}

func UserAgent() string {
	return fmt.Sprintf("rcc/%s (%s %s) %s", Version, runtime.GOOS, runtime.GOARCH, ControllerIdentity())
}

func ControllerIdentity() string {
	return strings.ToLower(fmt.Sprintf("rcc.%s", ControllerType))
}

func isFile(pathname string) bool {
	stat, err := os.Stat(pathname)
	return err == nil && stat.Mode().IsRegular()
}

func isDir(pathname string) bool {
	stat, err := os.Stat(pathname)
	return err == nil && stat.IsDir()
}

func ensureDirectory(name string) {
	if !isDir(name) {
		Error("mkdir", os.MkdirAll(name, 0o750))
	}
}

func UserHomeIdentity() string {
	if UnmanagedSpace {
		return "UNMNGED"
	}
	location, err := os.UserHomeDir()
	if err != nil {
		return "badcafe"
	}
	digest := fmt.Sprintf("%02x", Siphash(9007799254740993, 2147487647, []byte(location)))
	return digest[:7]
}
