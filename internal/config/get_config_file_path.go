package config 

import (
	"os"
	"path/filepath"
    "fmt"
)

func getConfigFilePath() (string, error){
    dirPath, err := os.UserHomeDir()
    if err != nil{
        return fmt.Sprintln("err"), err
    }
    configPath := filepath.Join(dirPath, ".gatorconfig.json")
    return configPath, nil
}
