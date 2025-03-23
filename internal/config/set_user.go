package config

import (
	"encoding/json"
	"os"
)

func (c *Config) SetUser(userName string) error{
    c.CurrentUserName = userName
    fileName, err := getConfigFilePath()
    if err != nil{
        return err
    }
    file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
    if err != nil{
        return err
    }
    defer file.Close()
    encoder := json.NewEncoder(file)
    encoder.SetIndent("", " ")
    err = encoder.Encode(c)
    if err != nil{
        return err
    }
    return nil
}
