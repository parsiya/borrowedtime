package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/mholt/archiver"
	"github.com/mitchellh/go-homedir"
	"github.com/parsiya/borrowedtime/shared"
	"github.com/parsiya/golnk"
)

// initiateConfig creates the homedir/borrowedtime directory and copies the config
// files. Next it opens the config file with notepad.exe on Windows.
// TODO: Add default editors for other OS.
// TODO: Editor detection, detect some popular editors and create commented
// entries for them in the config file. ~~Needs lnk parser.~~ Lnk parser is done,
// need some popular editors.
// TODO: Add logging. E.g. creating file X.
func initiateConfig(overwrite bool) error {
	// Delete templates directory and config.json.
	if overwrite {
		// Delete config file.
		cfgFile, _ := ConfigFilePath()
		err := shared.DeletePath(cfgFile)
		if err != nil {
			return fmt.Errorf("config.initiateConfig: %s", err.Error())
		}

		tmplDir, _ := TemplateDir()
		err = shared.DeletePath(tmplDir)
		if err != nil {
			return fmt.Errorf("config.initiateConfig: %s", err.Error())
		}
	}

	exists, err := configDirExists()
	if err != nil {
		return fmt.Errorf("config.initiateConfig: %s", err.Error())
	}
	// If exists, return an error because we do not want extra inits to overwrite everything.
	if exists && !overwrite {
		return fmt.Errorf("config.initiateConfig: config already exists, use Reset()")
	}

	// Create "borrowedtime/templates" in home directory. MkdirAll creates the
	// parents if needed.
	// We can ignore the error because we already called ConfigDir()
	tmplDir, _ := TemplateDir()
	err = os.MkdirAll(tmplDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("config.initiateConfig: create templates directory - %s", err.Error())
	}

	bckDir, _ := backupDir()
	err = os.MkdirAll(bckDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("config.initiateConfig: create backups directory - %s", err.Error())
	}

	// Copy default templates and configurations to populate templates.
	// Add everything from defaultTemplates map.
	for name, content := range defaultTemplates {
		err = AddTemplate(name, content, true)
		if err != nil {
			return fmt.Errorf("config.initiateConfig: add template %s - %s", name, err.Error())
		}
	}

	// Create data directory and copy data files.
	dataDir, _ := dataDir()
	err = os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("config.initiateConfig: create data directory - %s", err.Error())
	}

	// for name, content := range defaultData {
	// 	err = AddData(name, content, true)
	// 	if err != nil {
	// 		return fmt.Errorf("config.initiateConfig: add data %s - %s", name, err.Error())
	// 	}
	// }

	// This needs to be done after template and data file creation so we can add
	// them to the config file.
	err = createDefaultConfig()
	if err != nil {
		return fmt.Errorf("config.initiateConfig: %s", err.Error())
	}
	return nil
}

// Deploy calls initiateConfig() to create the configuration structure and file.
// Returns an error if config already exists.
func Deploy() error {
	return initiateConfig(false)
}

// Reset creates a backup and then calls initiateConfig(overwrite=true)
// and overwrites deployment. If no backup filename is provided, it will
// create one based on timestamp.
func Reset(createBackup bool, backupFile string) error {
	if createBackup {
		// Create a backup.
		if backupFile == "" {
			// Create a backup file based on timestamp.
			backupFile = time.Now().Format("2006-01-02-15-04-05") + "-reset"
		}
		err := Backup(backupFile)
		if err != nil {
			return fmt.Errorf("config.Reset: create backup - %s", err.Error())
		}
	}
	return initiateConfig(true)
}

// Backup creates a zip file from the "homedir/borrowedtime/templates"
// directory and stores it in the "homedir/borrowedtime/backups" directory.
// Go's zip directory is very basic so we use: https://github.com/mholt/archiver.
// If input is empty, file name will be timestamp.
func Backup(filename string) error {
	var backupFilename string
	backupDir, err := backupDir()
	if err != nil {
		return fmt.Errorf("config.Backup: %s", err.Error())
	}
	dataDir, _ := dataDir()
	// If input is empty, use the timestamp.
	if filename == "" {
		backupTimestamp := time.Now().Format("2006-01-02-15-04-05")

		backupFilename = path.Join(backupDir, backupTimestamp+".zip")
	} else {

		if path.Ext(filename) == "" {
			// If filename does not have an extension, pass zip.
			backupFilename = path.Join(backupDir, filename+".zip")
		} else {
			// Otherwise, use the extension in the filename.
			backupFilename = path.Join(backupDir, filename)
		}
	}

	templatesDir, err := TemplateDir()
	if err != nil {
		return fmt.Errorf("config.Backup: %s", err.Error())
	}
	cfgPath, _ := ConfigFilePath()
	return archiver.Zip.Make(backupFilename, []string{templatesDir, cfgPath, dataDir})
}

// Restore restores config.json and the templates directory from a backup file.
// Don't provide fullpath, just the filename including the extension.
func Restore(filename string) error {
	cfgDir, _ := configDir()
	backupDir, _ := backupDir()
	backupFile := ""
	// Check if backup path is absolute or relative to backupDir.
	if filepath.IsAbs(filename) {
		backupFile = filename
	} else {
		// It's relative to backupDir, create absolute path.
		backupFile = filepath.Join(backupDir, filename)
	}
	// Check if backup file exists.
	exists, err := shared.PathExists(backupFile)
	if err != nil {
		return fmt.Errorf("config.Restore: %s", err.Error())
	}
	if !exists {
		return fmt.Errorf("config.Restore: back up file %s not found", backupFile)
	}
	return archiver.Zip.Open(backupFile, cfgDir)
}

// BackupFiles returns all files in the backup directory.
func BackupFiles() (fi []string, err error) {
	backupDir, err := backupDir()
	if err != nil {
		return fi, fmt.Errorf("config.BackupFiles: %s", err.Error())
	}
	// List all files in the backup directory.
	return shared.ListFiles(backupDir, "*.*")
}

// CreateDefault creates the default config file and overwrites whatever
// is at "borrowedtime/config.json", finally opens it with OS default editor.
func CreateDefault() error {
	return createDefaultConfig()
}

// createDefaultConfig creates the default file and then opens it with notepad
// on Windows.
// TODO: Add a default editor function somewhere based on OS.
// TODO: Convert config creation to a template and pass a config struct instead.
func createDefaultConfig() error {

	cfgPath, err := ConfigFilePath()
	if err != nil {
		return fmt.Errorf("config.createDefaultConfig: get config file path - %s", err.Error())
	}

	// Read defaultConfig and create the config slice.
	defaultCfg := make(ConfigMap)
	err = json.Unmarshal([]byte(defaultWorkspaceConfig), &defaultCfg)
	if err != nil {
		return fmt.Errorf("config.createDefaultConfig: unmarshal default config - %s", err.Error())
	}

	// Set project structure template name.
	defaultCfg.Set("projectstructure", "project-structure")

	// Set workspace to Desktop\\projects on Windows.
	if runtime.GOOS == "windows" {
		desktop, err := shared.DesktopPath()
		if err != nil {
			return fmt.Errorf("config.createDefaultConfig: get Windows desktop path - %s", err.Error())
		}
		prjPath := path.Join(desktop, "projects")
		defaultCfg.Set("workspace", filepath.FromSlash(prjPath))

		// Set default editor to VS Code if it exists.
		defaultCfg.Set("editor", detectApp("code.exe"))
	}

	if err := Write(defaultCfg); err != nil {
		return fmt.Errorf("config.createDefaultConfig: %s", err.Error())
	}

	return shared.OpenWithDefaultEditor(cfgPath)
}

// detectApp parses the start menu, looks for an executable name
// (e.g. code.exe) in base paths, and returns the complete path.
func detectApp(executable string) string {
	if executable == "" {
		return ""
	}

	paths, err := parseStartMenu()
	if err != nil {
		return ""
	}

	for _, p := range paths {
		if strings.Contains(strings.ToLower(p), strings.ToLower(executable)) {
			return p
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

	// Create the location of user start menu.
	// "homedir/AppData/Roaming/Microsoft/Windows/Start Menu/Programs"
	home, _ := homedir.Dir()
	startMenuUser := path.Join(home, "AppData/Roaming/Microsoft/Windows/Start Menu/Programs")
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
	exists, err := shared.PathExists(root)
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

// Read reads the configuration file at "homedir/borrowedtime/config.json"
// and returns a populated map[string]string of keys.
func Read() (ConfigMap, error) {
	cfg := make(ConfigMap)
	cfgFilePath, err := ConfigFilePath()
	if err != nil {
		return cfg, fmt.Errorf("config.Read: get config file path %s", err.Error())
	}
	cfgContent, err := shared.ReadFileByte(cfgFilePath)
	if err != nil {
		return cfg, fmt.Errorf("config.Read: read config file %s", err.Error())
	}
	err = json.Unmarshal(cfgContent, &cfg)
	return cfg, err
}

// Write writes the config in map[string]string to the config file.
func Write(cfg ConfigMap) error {

	if len(cfg) == 0 {
		return fmt.Errorf("config.Write: empty config map")
	}

	cfgPath, _ := ConfigFilePath()
	configFile, err := os.Create(cfgPath)
	if err != nil {
		return fmt.Errorf("config.Write: %s", err.Error())
	}
	defer configFile.Close()

	switch runtime.GOOS {
	case "windows":
		// If on Windows, we need to replace \n with \r\n so notepad will show the
		// files properly.
		// var sb strings.Builder
		// enc := json.NewEncoder(&sb)
		// enc.SetIndent("", "\t")
		// if err := enc.Encode(cfg); err != nil {
		// 	return fmt.Errorf("config.Write: write config file - %s", err.Error())
		// }

		cfgString, err := shared.StructToJSONString(cfg, true)
		if err != nil {
			return fmt.Errorf("config.Write: write config file - %s", err.Error())
		}
		configFile.WriteString(shared.WindowsifyString(cfgString))
	default:
		// If not Windows, indent cfg and write it to the file.
		enc := json.NewEncoder(configFile)
		enc.SetIndent("", "\t")
		if err := enc.Encode(cfg); err != nil {
			return fmt.Errorf("config.Write: write config file - %s", err.Error())
		}
	}
	return nil
}

// configDir returns the config directory. Default: "homedir/borrowedtime".
func configDir() (string, error) {
	homedir, err := shared.HomeDir()
	if err != nil {
		return "", fmt.Errorf("config.ConfigDir: %s", err.Error())
	}
	return path.Join(homedir, "borrowedtime"), nil
}

// backupDir returns the backup directory.
// "homedir/borrowedtime/backups" or "ConfigDir/backups"
func backupDir() (string, error) {
	configDir, err := configDir()
	if err != nil {
		return "", fmt.Errorf("config.backupDir: %s", err.Error())
	}
	return path.Join(configDir, "backups"), nil
}

// dataDir returns the data directory.
// "homedir/borrowedtime/data" or "ConfigDir/data"
func dataDir() (string, error) {
	configDir, err := configDir()
	if err != nil {
		return "", fmt.Errorf("config.dataDir: %s", err.Error())
	}
	return path.Join(configDir, "data"), nil
}

// configDirExists returns true if it exists and any errors.
func configDirExists() (bool, error) {
	// Get config directory.
	configDir, err := configDir()
	if err != nil {
		return false, fmt.Errorf("config.configDirExists: %s", err.Error())
	}
	// Check if borrowedtime directory already exists.
	exists, err := shared.PathExists(configDir)
	// Return an error if we cannot access it.
	if err != nil {
		return false, fmt.Errorf("config.configDirExists: %s", err.Error())
	}
	// If it exists, return true.
	return exists, nil
}

// ConfigFilePath returns the path of the config file.
func ConfigFilePath() (string, error) {
	configDir, err := configDir()
	if err != nil {
		return "", fmt.Errorf("config.configFile: %s", err.Error())
	}
	return path.Join(configDir, defaultWorkspaceConfigFilename), nil
}

// Edit attempts to open the config file with editor.
func Edit(editor string) error {
	cfg, err := ConfigFilePath()
	if err != nil {
		return fmt.Errorf("config.Edit: %s", err.Error())
	}
	return shared.OpenWithEditor(editor, cfg)
}
