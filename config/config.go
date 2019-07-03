package config

import (
	"compress/flate"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/mholt/archiver"
	"github.com/parsiya/borrowedtime/shared"
)

// initiateConfig creates the homedir/borrowedtime directory and copies the
// configfiles. Next it opens the config file with notepad.exe on Windows.
// 1. Check if config directory exists.
// 2. Return with an error if it exists and overwrite is not set.
// 3. Delete the config directory.
// 4. Create the directory structure.
// TODO: Add default editors for other OS.
// TODO: Editor detection, detect some popular editors and create commented
// entries for them in the config file. ~~Needs lnk parser.~~ Lnk parser is done,
// need some popular editors.
func initiateConfig(overwrite bool) error {

	// 1. Check if config directory exists.
	exists, err := configDirExists()
	if err != nil {
		return fmt.Errorf("config.initiateConfig: %s", err.Error())
	}
	// 2. Return with an error if it exists and overwrite is not set.
	if exists && !overwrite {
		return fmt.Errorf("config.initiateConfig: config already exists, use \"config reset\"")
	}

	// 3. Delete the config directory.
	// We can safely ignore any errors here because configDirExists was executed
	// successfully.
	configDir, _ := configDir()
	if err = shared.DeletePath(configDir); err != nil {
		return fmt.Errorf("config.initiateConfig: delete the config directory - %s", err.Error())
	}

	// 4. Create the directory structure.
	// TODO: There should be a better way of doing this. See issue-20.
	// Create "borrowedtime/templates" in home directory. MkdirAll creates
	// parents if needed.
	// We can ignore the error here because we already called ConfigDir().
	tmplDir, _ := templateDir()
	err = os.MkdirAll(tmplDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("config.initiateConfig: create the templates directory - %s", err.Error())
	}
	// Create the file sub-directory inside templates.
	fileTmplDir, _ := fileTemplateDir()
	err = os.MkdirAll(fileTmplDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("config.initiateConfig: create the %s directory - %s",
			fileTmplDir, err.Error())
	}
	// Create the project sub-directory inside templates.
	prjTmplDir, _ := projectTemplateDir()
	err = os.MkdirAll(prjTmplDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("config.initiateConfig: create the %s directory - %s",
			prjTmplDir, err.Error())
	}

	bckDir, _ := backupDir()
	err = os.MkdirAll(bckDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("config.initiateConfig: create backups directory - %s", err.Error())
	}

	// Copy default templates and configurations to populate templates.
	// Add everything from defaultFileTemplates.
	for name, content := range defaultFileTemplates {
		err = addFileTemplate(name, content, true)
		if err != nil {
			return fmt.Errorf("config.initiateConfig: add template %s - %s", name, err.Error())
		}
	}
	// Add everything from defaultProjectTemplates.
	for name, content := range defaultProjectTemplates {
		err = addFileTemplate(name, content, true)
		if err != nil {
			return fmt.Errorf("config.initiateConfig: add template %s - %s", name, err.Error())
		}
	}

	// Create data directory and copy data files if any.
	dataDir, _ := dataDir()
	err = os.MkdirAll(dataDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("config.initiateConfig: create data directory - %s", err.Error())
	}

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

		backupFilename = filepath.Join(backupDir, backupTimestamp+".zip")
	} else {

		if filepath.Ext(filename) == "" {
			// If filename does not have an extension, pass zip.
			filename = shared.AddExtension(filename, "zip")
			backupFilename = filepath.Join(backupDir, filename)
		} else {
			// Otherwise, use the extension in the filename.
			backupFilename = filepath.Join(backupDir, filename)
		}
	}

	templatesDir, err := templateDir()
	if err != nil {
		return fmt.Errorf("config.Backup: %s", err.Error())
	}
	cfgPath, _ := ConfigFilePath()
	zip := archiver.Zip{CompressionLevel: flate.DefaultCompression}
	return zip.Archive([]string{templatesDir, cfgPath, dataDir}, backupFilename)
}

// Restore restores config.json and the templates directory from a backup file.
// Don't provide fullpath, just filename and extension.
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
	zip := archiver.Zip{
		CompressionLevel:  flate.DefaultCompression,
		MkdirAll:          true,
		OverwriteExisting: true,
	}
	return zip.Unarchive(backupFile, cfgDir)
}

// BackupFiles returns all files in the backup directory.
func BackupFiles() (fi []string, err error) {
	backupDir, err := backupDir()
	if err != nil {
		return fi, fmt.Errorf("config.BackupFiles: %s", err.Error())
	}
	// List all files in the backup directory.
	return shared.ListFiles(backupDir, "*")
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
	err = json.Unmarshal([]byte(defaultConfig), &defaultCfg)
	if err != nil {
		return fmt.Errorf("config.createDefaultConfig: unmarshal default config - %s", err.Error())
	}

	// Set project structure template name.
	defaultCfg.Set("projectstructure", "project-structure")

	// Set workspace to Desktop/projects on Windows.
	if runtime.GOOS == "windows" {
		desktop, err := shared.DesktopPath()
		if err != nil {
			return fmt.Errorf("config.createDefaultConfig: get Windows desktop path - %s", err.Error())
		}
		prjPath := filepath.Join(desktop, "projects")
		defaultCfg.Set("workspace", filepath.ToSlash(prjPath))

		// Set default editor to VS Code if it exists.
		defaultCfg.Set("editor", shared.DetectApp("code.exe"))
	}

	if err := Write(defaultCfg); err != nil {
		return fmt.Errorf("config.createDefaultConfig: %s", err.Error())
	}

	// Open the config file with the default editor.
	if err := shared.OpenWithDefaultEditor(cfgPath); err != nil {
		return fmt.Errorf("config.createDefaultConfig: open config with default editor - %s", err.Error())
	}
	return nil
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
		// If on Windows, we need to replace \n with \r\n so notepad will show
		// the files properly.
		cfgString, err := shared.StructToJSONString(cfg, true)
		if err != nil {
			return fmt.Errorf("config.Write: write config file - %s", err.Error())
		}
		_, err = configFile.WriteString(shared.WindowsifyString(cfgString))
		if err != nil {
			return fmt.Errorf("config.Write: write config file - %s", err.Error())
		}
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
	return filepath.Join(homedir, "borrowedtime"), nil
}

// backupDir returns the backup directory.
// "homedir/borrowedtime/backups" or "ConfigDir/backups"
func backupDir() (string, error) {
	configDir, err := configDir()
	if err != nil {
		return "", fmt.Errorf("config.backupDir: %s", err.Error())
	}
	return filepath.Join(configDir, "backups"), nil
}

// DataDir returns the data directory.
// "homedir/borrowedtime/data" or "ConfigDir/data"
func dataDir() (string, error) {
	configDir, err := configDir()
	if err != nil {
		return "", fmt.Errorf("config.DataDir: %s", err.Error())
	}
	return filepath.Join(configDir, "data"), nil
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
	return filepath.Join(configDir, defaultConfigFilename), nil
}

// Edit attempts to open the config file and the borrowed time directory with
// editor.
func Edit(editor string) error {
	cfg, err := ConfigFilePath()
	if err != nil {
		return fmt.Errorf("config.Edit: %s", err.Error())
	}
	// We can ignore the error here because we have already called this and
	// handled the error (if any) in "ConfigFilePath()."
	cfgDir, _ := configDir()
	return shared.OpenWithEditor(editor, cfg, cfgDir)
}

// ConfigDir is the exported version of configDir.
func ConfigDir() (string, error) {
	return configDir()
}
