package manager

import (
	"fmt"
	"os"
	"path/filepath"
)

// Check an error and panic if it exists.
func (ms *ManagerService) PanicCheck(err error) {
    if err != nil {
        panic(err)
    }
}

// GetFiles returns a list of files, given a root path to search, excluding
// non-files (., .., directories)
func (ms *ManagerService) GetFilesInPath(root string) ([]string, error) {
        var files []string

        err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
                if !info.IsDir() {
                        files = append(files, path)
                }
                return nil
        })
        return files, err
}

// get the property for config section as a string, if property does not exist
// return err
func (ms *ManagerService) GetSectionProperty(section string, property string) (string, error) {
	if ms.Config.Get(section + "." + property) != nil {
		return ms.Config.Get(section + "." + property).(string), nil
	} else {
		return "", fmt.Errorf("Configuration missing '%s' property under [%s] section.\n", property, section)
	}
}

// get the property for config section as a string, if property does not exist
// return the default property
func (ms *ManagerService) GetSectionPropertyOrDefault(section string, property string, def string) string {
	if ms.Config.Get(section + "." + property) != nil {
		return ms.Config.Get(section + "." + property).(string)
	} else {
		return def
	}
}
