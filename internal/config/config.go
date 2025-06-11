package Config

import (
	"encoding/json"
	"os"
)
type Config struct{
	DBURL string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error){
	filepath, err := getConfigFilePath()
	if err != nil{
		return Config{}, err
	}

	file, err := os.Open(filepath)
	if err != nil{ 
		return Config{}, err
	}
	defer file.Close()
	
	newConfig := Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&newConfig); err != nil{
		return Config{}, err
	}
 
	return newConfig, nil
}

func (c *Config) SetUser(name string) error{
	c.CurrentUserName = name
	return write(*c) 
}

func getConfigFilePath() (string, error){
	filepath, err := os.UserHomeDir()
	if err != nil {
		return "",err 
	}
	filepath = filepath + "/.gatorconfig.json"

	return filepath, nil
}

func write(cfg Config) error{
	filepath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(cfg); err != nil{
		return err
	}
	return nil
}
