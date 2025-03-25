package config

import (
	"encoding/json"
	"os"
    "fmt"
)

func (c *Config) SetUser(userName string) error {
    c.CurrentUserName = userName
    fileName, err := getConfigFilePath()
    if err != nil {
        return err
    }
    
    // Marshal to JSON bytes
    data, err := json.MarshalIndent(c, "", "  ")
    if err != nil {
        return err
    }
    
    // Write file completely from scratch
    err = os.WriteFile(fileName, data, 0644)
    if err != nil {
        return err
    }
    return nil
}
