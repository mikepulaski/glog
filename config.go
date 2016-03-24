package glog

import (
	"os"
	"strconv"
	"strings"
)

// Config defines output options. See Apply.
type Config struct {
	// Output defines the output medium.
	Output OutputType `toml:"output",json:"output"`

	// OutputDir defines the directory which logs will be written to. Must exist.
	OutputDir string `toml:"output_dir",json:"output_dir"`

	// Verbosity
	Verbosity        int `toml:"verbosity",json:"verbosity"`
	VerbosityModules []struct {
		Module    string `toml:"module",json:"module"`
		Verbosity int    `toml:"verbosity",json:"verbosity"`
	} `toml:"verbosity_modules",json:"verbosity_modules"`

	// StdErrThreshold defines the verbosity threshold which triggers message to log to stderr (in addition to file, if applicable)
	StdErrThreshold int32 `toml:"-",json:"-"`
}

// DefaultConfig defines the default config.
var DefaultConfig Config

func init() {
	DefaultConfig.Verbosity = 9000 // silent by default
	DefaultConfig.StdErrThreshold = 3
	DefaultConfig.Output = OutputStdErr
	DefaultConfig.OutputDir = os.TempDir()

	Apply(DefaultConfig)
}

// OutputType defines the output medium.
type OutputType uint16

const (
	// OutputFile logs output to a file.
	OutputFile OutputType = iota

	// OutputStdErr logs output to stderr.
	OutputStdErr

	// OutputFile logs output to a file and stderr.
	OutputBoth
)

// String implements the fmt.Stringer interface.
func (ot OutputType) String() string {
	switch ot {
	case OutputFile:
		return "File"
	case OutputStdErr:
		return "StdErr"
	case OutputBoth:
		return "Both"
	default:
		return "Unknown"
	}
}

// MarshalText implements the encoding.TextMarshaler interface.
func (ot OutputType) MarshalText() ([]byte, error) {
	return []byte(ot.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (ot *OutputType) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case "file":
		*ot = OutputFile
	case "stderr":
		*ot = OutputStdErr
	case "both":
		*ot = OutputBoth
	}

	return nil
}

// Apply sets the given config. Apply should be called during initialization; it is not
// safe for concurrent use.
func Apply(c Config) {
	// Normalize values
	if c.StdErrThreshold == 0 {
		c.StdErrThreshold = DefaultConfig.StdErrThreshold
	}

	if c.Verbosity == 0 {
		c.Verbosity = DefaultConfig.Verbosity
	}

	if c.OutputDir == "" {
		c.OutputDir = DefaultConfig.OutputDir
	}

	// Mimic flag.Set
	logging.toStderr = (c.Output == OutputStdErr)
	logging.alsoToStderr = (c.Output == OutputBoth)
	logging.stderrThreshold.Set(strconv.FormatInt(int64(c.StdErrThreshold), 10))
	logging.verbosity.Set(strconv.FormatInt(int64(c.Verbosity), 10))
	logging.vmodule.Set(strings.Join(c.VerbosityModulePatterns, ","))

	// See glog_file.go
	logDirs = []string{c.OutputDir}
}
