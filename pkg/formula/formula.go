package formula

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/ZupIT/ritchie-cli/pkg/api"
)

const (
	PathPattern               = "%s/formulas/%s"
	TmpDirPattern             = "%s/tmp/%s"
	TmpBinDirPattern          = "%s/tmp/%s/%s"
	DefaultConfig             = "config.json"
	ConfigPattern             = "%s/%s"
	CommandEnv                = "COMMAND"
	PwdEnv                    = "PWD"
	BinPattern                = "%s%s"
	BinPathPattern            = "%s/bin"
	windows                   = "windows"
	darwin                    = "darwin"
	linux                     = "linux"
	EnvPattern                = "%s=%s"
	CachePattern              = "%s/.%s.cache"
	DefaultCacheNewLabel      = "Type new value?"
	DefaultCacheQty           = 5
	FormCreatePathPattern     = "%s/ritchie-formulas-local"
	TreeCreatePathPattern     = "%s/tree/tree.json"
	MakefileCreatePathPattern = "%s/%s"
	Makefile                  = "Makefile"
	TreeJson                  = "tree.json"
)

// Config type that represents formula config
type Config struct {
	Name        string  `json:"name"`
	Command     string  `json:"command"`
	Description string  `json:"description"`
	Language    string  `json:"language"`
	Inputs      []Input `json:"inputs"`
}

// Input type that represents input config
type Input struct {
	Name    string   `json:"name"`
	Type    string   `json:"type"`
	Default string   `json:"default"`
	Label   string   `json:"label"`
	Items   []string `json:"items"`
	Cache   Cache    `json:"cache"`
}

type Cache struct {
	Active   bool   `json:"active"`
	Qty      int    `json:"qty"`
	NewLabel string `json:"newLabel"`
}
type Create struct {
	FormulaCmd   string `json:"formulaCmd"`
	Lang         string `json:"lang"`
	LocalRepoDir string `json:"localRepoDir"`
}

// Definition type that represents a Formula
type Definition struct {
	Path     string
	Bin      string
	LBin     string
	MBin     string
	WBin     string
	Bundle   string
	Config   string
	RepoURL  string
	RepoName string
}

type Setup struct {
	pwd            string
	formulaPath    string
	binPath        string
	tmpDir         string
	tmpBinDir      string
	tmpBinFilePath string
	config         Config
	containerId    string
}

// FormulaPath builds the formula path from ritchie home
func (d *Definition) FormulaPath(home string) string {
	return fmt.Sprintf(PathPattern, home, d.Path)
}

// TmpWorkDirPath builds the tmp paths to run formula, first parameter is tmpDir created
// second parameter is tmpBinDir
func (d *Definition) TmpWorkDirPath(home, uuidHash string) (string, string) {
	tmpDir := fmt.Sprintf(TmpDirPattern, home, uuidHash)
	tmpBinDir := fmt.Sprintf(TmpBinDirPattern, home, uuidHash, d.Path)
	return tmpDir, tmpBinDir
}

// BinName builds the bin name from definition params
func (d *Definition) BinName() string {
	bName := d.Bin
	so := runtime.GOOS
	switch so {
	case windows:
		if d.WBin != "" {
			bName = d.WBin
		}
	case darwin:
		if d.MBin != "" {
			bName = d.MBin
		}
	case linux:
		if d.LBin != "" {
			bName = d.LBin
		}
	default:
		bName = d.Bin
	}

	if strings.Contains(bName, "${so}") {
		suffix := ""
		if so == windows {
			suffix = ".exe"
		}
		binSO := strings.ReplaceAll(bName, "${so}", so)

		return fmt.Sprintf(BinPattern, binSO, suffix)
	}
	return bName
}

// BinName builds the bin name from definition params
func (d *Definition) BundleName() string {
	if strings.Contains(d.Bundle, "${so}") {
		so := runtime.GOOS
		bundleSO := strings.ReplaceAll(d.Bundle, "${so}", so)

		return bundleSO
	}
	return d.Bundle
}

// BinPath builds the bin path from formula path
func (d *Definition) BinPath(formula string) string {
	return fmt.Sprintf(BinPathPattern, formula)
}

// BinFilePath builds the bin file path from binPath and binName
func (d *Definition) BinFilePath(binPath, binName string) string {
	return fmt.Sprintf("%s/%s", binPath, binName)
}

// BundleURL builds the bundle url
func (d *Definition) BundleURL() string {
	return fmt.Sprintf("%s/%s/%s", d.RepoURL, d.Path, d.BundleName())
}

// ConfigName resolver de config name
func (d *Definition) ConfigName() string {
	if d.Config != "" {
		return d.Config
	}
	return DefaultConfig
}

// ConfigPath builds the config path from formula path and config name
func (d *Definition) ConfigPath(formulaPath, configName string) string {
	return fmt.Sprintf(ConfigPattern, formulaPath, configName)
}

// ConfigURL builds the config url
func (d *Definition) ConfigURL(configName string) string {
	return fmt.Sprintf("%s/%s/%s", d.RepoURL, d.Path, configName)
}

type PreRunner interface {
	PreRun(def Definition) (Setup, error)
}

type Runner interface {
	Run(def Definition, inputType api.TermInputType) error
}

type PostRunner interface {
	PostRun(p Setup, docker bool) error
}

type InputRunner interface {
	Inputs(cmd *exec.Cmd, setup Setup, inputType api.TermInputType) error
}

type Setuper interface {
	Setup(def Definition) (Setup, error)
}

type Creator interface {
	Create(cf Create) (CreateManager, error)
}
