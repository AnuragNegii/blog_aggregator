package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func Read() (Config, error){
    config := Config{}

    configPath, err := getConfigFilePath()

    file, err := os.Open(configPath)
    if err != nil {
        return config, fmt.Errorf("Error getting file: %s\n",err) 
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil{
        return config, fmt.Errorf("Error reading data: %s\n",err) 
    }
    
    err = json.Unmarshal(data, &config)
    if err != nil{
        return config, fmt.Errorf("Error unmarshalling file: %s\n",err) 
    }
    return config, nil
}
