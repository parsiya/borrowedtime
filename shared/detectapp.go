package shared

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	lnk "github.com/parsiya/golnk"
)

// Start menu path
var startMenuAllUsers = "C:/ProgramData/Microsoft/Windows/Start Menu/Programs"

// DetectApp parses the start menu, looks for an executable name
// (e.g. code.exe) in base paths, and returns the complete path.
func DetectApp(executable string) string {
	if executable == "" {
		return ""
	}

	paths, err := parseStartMenu()
	if err != nil {
		return ""
	}

	for _, p := range paths {
		if strings.Contains(strings.ToLower(p), strings.ToLower(executable)) {
			return filepath.ToSlash(p)
		}
	}
	return ""
}

// parseStartMenu returns where the lnk files in Windows start menu are pointing
// to. To do this, I implemented a lnk parser at https://github.com/parsiya/golnk.
func parseStartMenu() (basePaths []string, err error) {
	// Check OS.
	if runtime.GOOS != "windows" {
		return basePaths,
			fmt.Errorf("config.parseStartMenu: not running on windows, running %s", runtime.GOOS)
	}

	// Now parse two locations.
	// var startMenuAllUsers = "C:/ProgramData/Microsoft/Windows/Start Menu/Programs"
	b1, err := parseLnk(startMenuAllUsers)
	if err != nil {
		return basePaths, fmt.Errorf("config.parseStartMenu: parse all users start menu - %s", err.Error())
	}
	basePaths = append(basePaths, b1...)

	// Get the location of user start menu.
	// "homedir/AppData/Roaming/Microsoft/Windows/Start Menu/Programs"
	home, _ := homedir.Dir()
	startMenuUser := filepath.Join(home, "AppData/Roaming/Microsoft/Windows/Start Menu/Programs")
	b2, err := parseLnk(startMenuUser)
	if err != nil {
		// Do not return an error because we have already parsed the other one.
		return basePaths, nil
	}
	basePaths = append(basePaths, b2...)

	return basePaths, nil
}

// parseLnk parses all lnk files in the target location and subdirectories, and
// return the base paths.
func parseLnk(root string) (basePaths []string, err error) {
	exists, err := PathExists(root)
	if err != nil {
		return basePaths, fmt.Errorf("config.parseLnk: check root - %s", err.Error())
	}
	if !exists {
		return basePaths, fmt.Errorf("config.parseLnk: check root - %s", err.Error())
	}

	err = filepath.Walk(root, func(fpath string, info os.FileInfo, walkErr error) error {
		if filepath.Ext(fpath) == ".lnk" {
			fi, lnkErr := lnk.File(fpath)
			// If file was not parsed, move on.
			if lnkErr != nil {
				return nil
			}
			var targetPath = ""
			if fi.LinkInfo.LocalBasePath != "" {
				targetPath = fi.LinkInfo.LocalBasePath
			}
			if fi.LinkInfo.LocalBasePathUnicode != "" {
				targetPath = fi.LinkInfo.LocalBasePathUnicode
			}

			if targetPath != "" && filepath.Ext(targetPath) == ".exe" {
				basePaths = append(basePaths, targetPath)
			}
		}
		return nil
	})
	if err != nil {
		return basePaths, fmt.Errorf("borrowedtime.parseLnk: %s", err.Error())
	}
	return basePaths, nil
}
